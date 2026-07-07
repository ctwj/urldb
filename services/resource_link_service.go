package services

import (
	"context"
	"fmt"
	"time"

	panutils "github.com/ctwj/urldb/common"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// TransferResult 转存结果（由 handlers 迁移至 services，供网页端与电报机器人共用）
type TransferResult struct {
	Success  bool   `json:"success"`
	Fid      string `json:"fid"`
	SaveURL  string `json:"save_url"`
	ErrorMsg string `json:"error_msg"`
}

// LinkResolution 取链结果
type LinkResolution struct {
	URL         string
	Type        string // "transferred" | "original"
	Platform    string // pan.Remark（展示用）
	Transferred bool
	Note        string // 可选提示（如易和谐提醒）
}

// ResourceLinkService 资源取链服务：网页端 GetResourceLink 与电报机器人共用，避免逻辑漂移。
// 访问计数 / 来源归因由调用方各自处理（网页=web，机器人=telegram）。
// Feature: 011-telegram-bot-enhance
type ResourceLinkService interface {
	Resolve(ctx context.Context, resource *entity.Resource) (LinkResolution, error)
}

type resourceLinkServiceImpl struct {
	cksRepo      repo.CksRepository
	panRepo      repo.PanRepository
	configRepo   repo.SystemConfigRepository
	resourceRepo repo.ResourceRepository
}

// NewResourceLinkService 创建取链服务
func NewResourceLinkService(cksRepo repo.CksRepository, panRepo repo.PanRepository, configRepo repo.SystemConfigRepository, resourceRepo repo.ResourceRepository) ResourceLinkService {
	return &resourceLinkServiceImpl{cksRepo: cksRepo, panRepo: panRepo, configRepo: configRepo, resourceRepo: resourceRepo}
}

// Resolve 按网页端 GetResourceLink 同一逻辑解析可用链接：
// 非 quark/xunlei/baidu → 原链；已转存 → SaveURL；未开自动转存 → 原链；否则执行自动转存。
func (s *resourceLinkServiceImpl) Resolve(ctx context.Context, resource *entity.Resource) (LinkResolution, error) {
	_ = ctx
	if resource == nil {
		return LinkResolution{}, fmt.Errorf("resource 为空")
	}

	var panName, panRemark string
	if resource.PanID != nil {
		if pan, err := s.panRepo.FindByID(*resource.PanID); err == nil && pan != nil {
			panName = pan.Name
			panRemark = pan.Remark
		}
	}
	platform := panRemark
	if platform == "" {
		platform = panName
	}

	// 仅 quark/xunlei/baidu 支持详情页自动转存；其他平台直接返回原链
	if panName != "quark" && panName != "xunlei" && panName != "baidu" {
		return LinkResolution{URL: resource.URL, Type: "original", Platform: platform}, nil
	}
	// 已存在转存链接
	if resource.SaveURL != "" {
		return LinkResolution{URL: resource.SaveURL, Type: "transferred", Platform: platform, Transferred: true}, nil
	}
	// 自动转存开关
	autoTransfer, err := s.configRepo.GetConfigBool(entity.ConfigKeyAutoTransferEnabled)
	if err != nil || !autoTransfer {
		return LinkResolution{URL: resource.URL, Type: "original", Platform: platform}, nil
	}
	// 执行自动转存
	res := PerformAutoTransfer(s.cksRepo, s.configRepo, s.resourceRepo, resource)
	if res.Success {
		return LinkResolution{URL: res.SaveURL, Type: "transferred", Platform: platform, Transferred: true, Note: "资源易和谐，请及时转存"}, nil
	}
	utils.Error("[RESOURCE_LINK] 自动转存失败 (resource=%d): %s", resource.ID, res.ErrorMsg)
	return LinkResolution{URL: resource.URL, Type: "original", Platform: platform}, nil
}

// PerformAutoTransfer 执行自动转存（由 handlers 迁移，逻辑保持一致）。
// 传入所需仓库，避免依赖包级 repoManager，便于网页端与机器人共用。
func PerformAutoTransfer(cksRepo repo.CksRepository, configRepo repo.SystemConfigRepository, resourceRepo repo.ResourceRepository, resource *entity.Resource) TransferResult {
	utils.Info("开始执行资源转存 - ID: %d, URL: %s", resource.ID, resource.URL)

	panID := resource.PanID
	if panID == nil {
		return TransferResult{Success: false, ErrorMsg: "资源未关联网盘平台"}
	}

	accounts, err := cksRepo.FindByPanID(*panID)
	if err != nil {
		utils.Error("获取网盘账号失败: %v", err)
		return TransferResult{Success: false, ErrorMsg: fmt.Sprintf("获取网盘账号失败: %v", err)}
	}

	autoTransferMinSpace, err := configRepo.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	if err != nil {
		utils.Error("获取最小存储空间配置失败: %v", err)
		autoTransferMinSpace = 5 // 默认5GB
	}

	// 过滤：只保留已激活、同平台、剩余空间足够的账号
	minSpaceBytes := int64(autoTransferMinSpace) * 1024 * 1024 * 1024
	var validAccounts []entity.Cks
	for _, acc := range accounts {
		if !acc.IsValid {
			utils.Warn("跳过账号 ID=%d (%s)：IsValid=false", acc.ID, acc.Username)
			continue
		}
		if acc.PanID != *panID {
			utils.Warn("跳过账号 ID=%d (%s)：PanID 不匹配 (账号=%d, 资源=%d)", acc.ID, acc.Username, acc.PanID, *panID)
			continue
		}
		if acc.LeftSpace < minSpaceBytes {
			utils.Warn("跳过账号 ID=%d (%s)：剩余空间不足 (%d < %d bytes = %dGB)", acc.ID, acc.Username, acc.LeftSpace, minSpaceBytes, autoTransferMinSpace)
			continue
		}
		validAccounts = append(validAccounts, acc)
	}

	if len(validAccounts) == 0 {
		msg := fmt.Sprintf("没有可用的网盘账号 (候选 %d 个, 最小空间要求 %dGB)", len(accounts), autoTransferMinSpace)
		utils.Warn("%s", msg)
		return TransferResult{Success: false, ErrorMsg: msg}
	}

	utils.Info("找到 %d 个可用网盘账号，开始转存处理...", len(validAccounts))
	account := validAccounts[0]
	factory := panutils.NewPanFactory()
	result := transferSingle(cksRepo, resource, account, factory)

	if result.Success {
		// 更新资源的转存信息
		resource.SaveURL = result.SaveURL
		resource.Fid = result.Fid
		resource.CkID = &account.ID
		resource.ErrorMsg = ""
		// 详情页触发自动转存视为中转备份，写入 transferred_at 让自动清理调度器到期回收
		now := time.Now()
		resource.TransferredAt = &now
		// GORM Updates(struct) 会跳过零值字段，所以用 UpdateFields 显式更新
		if err := resourceRepo.UpdateFields(resource.ID, map[string]interface{}{
			"save_url":            result.SaveURL,
			"fid":                 result.Fid,
			"ck_id":               account.ID,
			"error_msg":           "",
			"transferred_at":      now,
			"updated_at":          now,
			"cleaned_at":          nil,
			"clean_error_msg":     "",
			"last_clean_error_at": nil,
		}); err != nil {
			utils.Error("更新资源转存信息失败: %v", err)
		}
	} else {
		resource.ErrorMsg = result.ErrorMsg
		if err := resourceRepo.Update(resource); err != nil {
			utils.Error("更新资源错误信息失败: %v", err)
		}
	}

	return result
}

// transferSingle 转存单个资源（由 handlers 迁移）
func transferSingle(cksRepo repo.CksRepository, resource *entity.Resource, account entity.Cks, factory *panutils.PanFactory) TransferResult {
	utils.Info("开始转存资源 - 资源ID: %d, 账号: %s", resource.ID, account.Username)

	service, err := factory.CreatePanService(resource.URL, &panutils.PanConfig{
		URL:         resource.URL,
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	if err != nil {
		utils.Error("创建网盘服务失败: %v", err)
		return TransferResult{Success: false, ErrorMsg: fmt.Sprintf("创建网盘服务失败: %v", err)}
	}

	// 设置账号信息
	service.SetCKSRepository(cksRepo, account)

	// 提取分享ID
	shareID, _ := panutils.ExtractShareId(resource.URL)
	if shareID == "" {
		return TransferResult{Success: false, ErrorMsg: "无效的分享链接"}
	}

	// 执行转存
	transferResult, err := service.Transfer(shareID)
	if err != nil {
		utils.Error("转存失败: %v", err)
		return TransferResult{Success: false, ErrorMsg: fmt.Sprintf("转存失败: %v", err)}
	}

	if transferResult == nil || !transferResult.Success {
		errMsg := "转存失败"
		if transferResult != nil && transferResult.Message != "" {
			errMsg = transferResult.Message
		}
		utils.Error("转存失败: %s", errMsg)
		return TransferResult{Success: false, ErrorMsg: errMsg}
	}

	// 提取转存链接
	var saveURL string
	var fid string

	if data, ok := transferResult.Data.(map[string]interface{}); ok {
		if v, ok := data["shareUrl"]; ok {
			saveURL, _ = v.(string)
		}
		if v, ok := data["fid"]; ok {
			fid, _ = v.(string)
		}
	}
	if saveURL == "" {
		saveURL = transferResult.ShareURL
	}

	if saveURL == "" {
		return TransferResult{Success: false, ErrorMsg: "转存成功但未获取到分享链接"}
	}

	utils.Info("转存成功 - 资源ID: %d, 转存链接: %s", resource.ID, saveURL)
	return TransferResult{Success: true, SaveURL: saveURL, Fid: fid}
}
