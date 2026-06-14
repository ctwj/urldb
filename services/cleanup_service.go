package services

import (
	"context"
	"strings"
	"time"

	pan "github.com/ctwj/urldb/common"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// CleanupService 转存文件自动清理服务
// 周期性扫描已转存且超过保留期的资源，调用网盘 API 删除文件并清空转存字段
type CleanupService struct {
	resourceRepo repo.ResourceRepository
	configRepo   repo.SystemConfigRepository
	cksRepo      repo.CksRepository
	panRepo      repo.PanRepository
}

// NewCleanupService 创建清理服务
func NewCleanupService(
	resourceRepo repo.ResourceRepository,
	configRepo repo.SystemConfigRepository,
	cksRepo repo.CksRepository,
	panRepo repo.PanRepository,
) *CleanupService {
	return &CleanupService{
		resourceRepo: resourceRepo,
		configRepo:   configRepo,
		cksRepo:      cksRepo,
		panRepo:      panRepo,
	}
}

// Run 执行一次清理扫描
// 返回处理统计：总数、成功数、失败数
func (s *CleanupService) Run(ctx context.Context) (total, success, failed int, err error) {
	startTime := time.Now()
	utils.Info("[CleanupService] 开始执行清理任务")

	// 读取保留天数配置，缺失或非法时使用默认值 7 天
	retentionDays, cfgErr := s.configRepo.GetConfigInt(entity.ConfigKeyAutoCleanupRetentionDays)
	if cfgErr != nil || retentionDays <= 0 {
		retentionDays = 7
		utils.Warn("[CleanupService] 读取保留天数配置失败或非法，使用默认值 7 天: %v", cfgErr)
	}

	// 每轮处理上限，避免长时间阻塞调度器（可后续接入并发配置）
	batchLimit := 100

	resources, err := s.resourceRepo.FindDueForCleanup(retentionDays, batchLimit)
	if err != nil {
		utils.Error("[CleanupService] 查询待清理资源失败: %v", err)
		return 0, 0, 0, err
	}

	total = len(resources)
	if total == 0 {
		utils.Info("[CleanupService] 没有待清理的资源（保留期=%d天）", retentionDays)
		return 0, 0, 0, nil
	}

	utils.Info("[CleanupService] 找到 %d 个待清理资源，保留期=%d天", total, retentionDays)

	// 按 ck_id 缓存账号信息，避免重复查询
	accountCache := make(map[uint]*entity.Cks)

	for _, res := range resources {
		// 检查上下文是否已取消（调度器停止时立即退出）
		select {
		case <-ctx.Done():
			utils.Warn("[CleanupService] 任务被取消，已处理 %d/%d", success+failed, total)
			return total, success, failed, nil
		default:
		}

		// 无效 fid 直接跳过（数据异常，记录失败但不重试）
		if res.Fid == "" {
			utils.Warn("[CleanupService] 资源 ID=%d 无 fid，跳过删除并标记失败", res.ID)
			_ = s.resourceRepo.MarkCleanError(res.ID, "fid 为空", time.Now())
			failed++
			continue
		}

		// 解析资源对应的账号 cookie（FR-012 防跨账号误删）
		account, accErr := s.resolveAccount(res, accountCache)
		if accErr != nil {
			utils.Error("[CleanupService] 资源 ID=%d 解析账号失败: %v", res.ID, accErr)
			_ = s.resourceRepo.MarkCleanError(res.ID, truncateMsg(accErr.Error()), time.Now())
			failed++
			continue
		}

		// 调用网盘 API 删除文件
		delErr := s.deleteFile(account, res.Fid)
		if delErr == nil {
			// 成功：清空 fid/save_url，写入 cleaned_at
			if err := s.resourceRepo.MarkCleaned(res.ID, time.Now()); err != nil {
				utils.Error("[CleanupService] 资源 ID=%d 文件已删除但更新数据库失败: %v", res.ID, err)
				// 文件已删除，不视为失败但记录错误
			} else {
				utils.Info("[CleanupService] 资源 ID=%d 清理成功", res.ID)
			}
			success++
		} else if isFileNotExist(delErr) {
			// 文件不存在视为成功（FR-009）
			utils.Info("[CleanupService] 资源 ID=%d 网盘文件已不存在，视为清理成功: %v", res.ID, delErr)
			_ = s.resourceRepo.MarkCleaned(res.ID, time.Now())
			success++
		} else {
			// 其他失败：保留 fid/save_url，等待下一轮重试（FR-006 不阻塞）
			utils.Error("[CleanupService] 资源 ID=%d 清理失败: %v", res.ID, delErr)
			_ = s.resourceRepo.MarkCleanError(res.ID, truncateMsg(delErr.Error()), time.Now())
			failed++
		}
	}

	elapsed := time.Since(startTime)
	utils.Info("[CleanupService] 清理任务完成，总计=%d, 成功=%d, 失败=%d, 耗时=%v", total, success, failed, elapsed)
	return total, success, failed, nil
}

// resolveAccount 解析资源对应的账号 cookie
// 通过 ck_id 查询账号，确保使用与转存时同一账号进行删除（防跨账号误删）
func (s *CleanupService) resolveAccount(res *entity.Resource, cache map[uint]*entity.Cks) (*entity.Cks, error) {
	if res.CkID == nil {
		return nil, errNoAccountBound
	}
	accID := *res.CkID

	// 命中缓存
	if acc, ok := cache[accID]; ok && acc != nil {
		return acc, nil
	}

	acc, err := s.cksRepo.FindByIds([]uint{accID})
	if err != nil {
		return nil, err
	}
	if len(acc) == 0 {
		return nil, errNoAccountBound
	}

	cache[accID] = acc[0]
	return acc[0], nil
}

// deleteFile 调用网盘 API 删除指定文件
// 根据账号 ServiceType 选择对应网盘服务实例
func (s *CleanupService) deleteFile(account *entity.Cks, fid string) error {
	serviceType := account.ServiceType
	if serviceType == "" {
		return errUnknownServiceType
	}

	factory := pan.NewPanFactory()
	// 创建带账号 cookie 的配置
	cfg := &pan.PanConfig{
		Cookie: account.Ck,
	}

	// 使用工厂创建对应类型的网盘服务
	service, err := factory.CreatePanServiceByType(toPanServiceType(serviceType), cfg)
	if err != nil {
		return err
	}

	// 设置 cks 仓库（部分操作需要刷新 token）
	if setter, ok := service.(interface {
		SetCKSRepository(repo.CksRepository, entity.Cks)
	}); ok {
		setter.SetCKSRepository(s.cksRepo, *account)
	}

	result, err := service.DeleteFiles([]string{fid})
	if err != nil {
		return err
	}
	if result == nil {
		return errDeleteFailed
	}
	if !result.Success {
		// 将 Message 作为错误返回，便于上层宽松匹配"文件不存在"
		return errMsg(result.Message)
	}
	return nil
}

// toPanServiceType 将 ServiceType 字符串转为枚举
func toPanServiceType(serviceType string) pan.ServiceType {
	switch strings.ToLower(serviceType) {
	case "quark":
		return pan.Quark
	case "alipan", "aliyun":
		return pan.Alipan
	case "baidu":
		return pan.BaiduPan
	case "uc":
		return pan.UC
	case "xunlei":
		return pan.Xunlei
	default:
		return pan.NotFound
	}
}

// isFileNotExist 判断错误是否表示"文件已不存在"
// 宽松匹配中英文关键字，避免被具体错误码绑死（FR-009）
func isFileNotExist(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "不存在") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "no such") ||
		strings.Contains(msg, "already deleted") ||
		strings.Contains(msg, "已删除")
}

// truncateMsg 截断错误信息以匹配字段长度（255）
func truncateMsg(msg string) string {
	if len(msg) > 255 {
		return msg[:255]
	}
	return msg
}

// 错误变量
var (
	errNoAccountBound     = errMsg("资源未绑定账号，无法确定清理使用的 cookie")
	errUnknownServiceType = errMsg("账号 ServiceType 为空，无法确定网盘类型")
	errDeleteFailed       = errMsg("删除文件失败：返回结果为空")
)

// errMsg 简单包装字符串为 error，便于上层匹配关键字
type errMsg string

func (e errMsg) Error() string { return string(e) }
