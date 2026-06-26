package pan

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// ============================================================================
// UC 网盘与夸克网盘同构（同一套接口结构、相同的端点路径与请求/响应字段），
// 仅以下 4 项 HTTP 调用配置不同。取值依据：demo/OpenList-main/drivers/quark_uc
// （QuarkOrUC 单驱动同时支撑夸克与 UC，UC 的 Conf 配置）。本项目 common/quark_pan.go
// 为转存链路的镜像模板，本文件按 UC 配置改写，不修改 quark_pan.go（FR-012 零回归）。
// ============================================================================

const (
	// ucAPIBase UC 网盘云盘 API 主机（夸克为 https://drive-pc.quark.cn/1/clouddrive）
	ucAPIBase = "https://pc-api.uc.cn/1/clouddrive"
	// ucPR UC 专属 pr 查询参数（夸克为 ucpro）
	ucPR = "UCBrowser"
	// ucReferer UC Referer 头（夸克为 https://pan.quark.cn/）
	ucReferer = "https://drive.uc.cn"
	// ucUA UC 专属 User-Agent（夸克为 quark-cloud-drive 变体）
	ucUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) uc-cloud-drive/2.5.20 Chrome/100.0.4896.160 Electron/18.3.5.4-b478491100 Safari/537.36 Channel/pckk_other_ch"
)

// UCService UC网盘服务
type UCService struct {
	*BasePanService
	configMutex sync.RWMutex // 保护配置的读写锁
}

// NewUCService 创建UC网盘服务
func NewUCService(config *PanConfig) *UCService {
	service := &UCService{
		BasePanService: NewBasePanService(config),
	}

	// 设置UC网盘的默认请求头（UC 专属 Referer / User-Agent）
	cookie := ""
	if config != nil {
		cookie = config.Cookie
	}
	service.SetHeaders(map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Content-Type":       "application/json;charset=UTF-8",
		"Referer":            ucReferer,
		"User-Agent":         ucUA,
		"Cookie":             cookie,
	})

	service.UpdateConfig(config)

	return service
}

// GetServiceType 获取服务类型
func (u *UCService) GetServiceType() ServiceType {
	return UC
}

// UpdateConfig 更新配置（线程安全）
func (u *UCService) UpdateConfig(config *PanConfig) {
	if config == nil {
		return
	}

	u.configMutex.Lock()
	defer u.configMutex.Unlock()

	u.config = config
	if config.Cookie != "" {
		u.SetHeader("Cookie", config.Cookie)
	}
}

// ucCommonQuery UC 接口公共查询参数（pr / fr）
func (u *UCService) ucCommonQuery() map[string]string {
	return map[string]string{
		"pr":           ucPR,
		"fr":           "pc",
		"uc_param_str": "",
	}
}

// configCode 线程安全地读取提取码
func (u *UCService) configCode() string {
	u.configMutex.RLock()
	defer u.configMutex.RUnlock()
	if u.config == nil {
		return ""
	}
	return u.config.Code
}

// ============================================================================
// Transfer 转存分享链接（与夸克同构：stoken → 详情 → 转存 → 轮询 → 广告处理 → 再分享 → 取码）
// ============================================================================

// Transfer 转存分享链接
func (u *UCService) Transfer(shareID string) (*TransferResult, error) {
	u.configMutex.RLock()
	config := u.config
	u.configMutex.RUnlock()

	log.Printf("开始处理UC分享: %s", shareID)

	// 获取 stoken
	var stoken string
	if config == nil || config.Stoken == "" {
		stokenResult, err := u.getStoken(shareID)
		if err != nil {
			return ErrorResult(fmt.Sprintf("获取stoken失败: %v", err)), nil
		}
		stoken = strings.ReplaceAll(stokenResult.Stoken, " ", "+")
	} else {
		stoken = strings.ReplaceAll(config.Stoken, " ", "+")
	}

	// 获取分享详情
	shareResult, err := u.getShare(shareID, stoken)
	if err != nil || len(shareResult.List) == 0 {
		return ErrorResult(fmt.Sprintf("获取分享详情失败: %v", err)), nil
	}

	// 检验模式：只读取资源信息，不实际转存（FR-007）
	if config != nil && config.IsType == 1 {
		return SuccessResult("检验成功", map[string]interface{}{
			"title":    shareResult.Share.Title,
			"shareUrl": config.URL,
		}), nil
	}

	// 提取文件信息
	fidList := make([]string, 0)
	fidTokenList := make([]string, 0)
	title := shareResult.Share.Title

	for _, item := range shareResult.List {
		fidList = append(fidList, item.Fid)
		fidTokenList = append(fidTokenList, item.ShareFidToken)
	}

	// 转存资源
	saveResult, err := u.getShareSave(shareID, stoken, fidList, fidTokenList)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}

	// 等待转存完成
	myData, err := u.waitForTask(saveResult.TaskID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("等待转存完成失败: %v", err)), nil
	}

	if len(myData.SaveAs.SaveAsTopFids) == 0 {
		return ErrorResult("转存完成但未获取到文件标识"), nil
	}

	// 删除广告文件（如果有配置）
	if err := u.deleteAdFiles(myData.SaveAs.SaveAsTopFids[0]); err != nil {
		log.Printf("删除UC广告文件失败: %v", err)
	}

	// 添加个人自定义广告
	if err := u.addAd(myData.SaveAs.SaveAsTopFids[0]); err != nil {
		log.Printf("添加UC广告文件失败: %v", err)
	}

	// 分享资源
	shareBtnResult, err := u.getShareBtn(myData.SaveAs.SaveAsTopFids, title)
	if err != nil {
		return ErrorResult(fmt.Sprintf("分享失败: %v", err)), nil
	}

	// 等待分享完成
	shareTaskResult, err := u.waitForTask(shareBtnResult.TaskID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("等待分享完成失败: %v", err)), nil
	}

	// 获取分享密码
	passwordResult, err := u.getSharePassword(shareTaskResult.ShareID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取分享密码失败: %v", err)), nil
	}

	// 确定 fid
	var fid string
	if len(myData.SaveAs.SaveAsTopFids) > 1 {
		fid = strings.Join(myData.SaveAs.SaveAsTopFids, ",")
	} else {
		fid = passwordResult.FirstFile.Fid
	}

	return SuccessResult("转存成功", map[string]interface{}{
		"shareUrl": passwordResult.ShareURL,
		"title":    passwordResult.ShareTitle,
		"fid":      fid,
		"code":     passwordResult.Code,
	}), nil
}

// ============================================================================
// GetFiles / DeleteFiles
// ============================================================================

// GetFiles 获取文件列表
func (u *UCService) GetFiles(pdirFid string) (*TransferResult, error) {
	if pdirFid == "" {
		pdirFid = "0"
	}

	queryParams := map[string]string{
		"pr":              ucPR,
		"fr":              "pc",
		"uc_param_str":    "",
		"pdir_fid":        pdirFid,
		"_page":           "1",
		"_size":           "50",
		"_fetch_total":    "1",
		"_fetch_sub_dirs": "0",
		"_sort":           "file_type:asc,updated_at:desc",
	}

	data, err := u.HTTPGet(ucAPIBase+"/file/sort", queryParams)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取UC文件列表失败: %v", err)), nil
	}

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			List []interface{} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return ErrorResult("解析UC响应失败"), nil
	}

	if response.Status != 200 {
		message := response.Message
		if strings.Contains(message, "require login") || strings.Contains(message, "guest") {
			message = "UC未登录，请检查cookie"
		}
		return ErrorResult(message), nil
	}

	return SuccessResult("获取成功", response.Data.List), nil
}

// DeleteFiles 删除文件
func (u *UCService) DeleteFiles(fileList []string) (*TransferResult, error) {
	if len(fileList) == 0 {
		return ErrorResult("文件列表为空"), nil
	}

	for _, fileID := range fileList {
		if err := u.deleteSingleFile(fileID); err != nil {
			log.Printf("删除UC文件 %s 失败: %v", fileID, err)
			return ErrorResult(fmt.Sprintf("删除UC文件 %s 失败: %v", fileID, err)), nil
		}
	}

	return SuccessResult("删除成功", nil), nil
}

// deleteSingleFile 删除单个文件
func (u *UCService) deleteSingleFile(fileID string) error {
	log.Printf("正在删除UC文件: %s", fileID)

	data := map[string]interface{}{
		"action_type":  2,
		"filelist":     []string{fileID},
		"exclude_fids": []string{},
	}

	respData, err := u.HTTPPost(ucAPIBase+"/file/delete", data, u.ucCommonQuery())
	if err != nil {
		return fmt.Errorf("删除UC文件请求失败: %v", err)
	}

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			TaskID string `json:"task_id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return fmt.Errorf("解析UC删除响应失败: %v", err)
	}

	if response.Status != 200 {
		return fmt.Errorf("删除UC文件失败: %s", response.Message)
	}

	if response.Data.TaskID != "" {
		log.Printf("UC删除文件任务ID: %s", response.Data.TaskID)
		if _, err := u.waitForTask(response.Data.TaskID); err != nil {
			return fmt.Errorf("等待UC删除任务完成失败: %v", err)
		}
	}

	return nil
}

// ============================================================================
// UC 转存链路内部方法（端点路径与夸克一致，主机/参数为 UC 取值）
// ============================================================================

// getStoken 提取码换 stoken
func (u *UCService) getStoken(shareID string) (*UCStokenResult, error) {
	data := map[string]interface{}{
		"passcode": u.configCode(), // 加密分享的提取码（公开分享为空）
		"pwd_id":   shareID,
	}

	respData, err := u.HTTPPost(ucAPIBase+"/share/sharepage/token", data, u.ucCommonQuery())
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int           `json:"status"`
		Message string        `json:"message"`
		Data    UCStokenResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("%s", response.Message)
	}

	return &response.Data, nil
}

// getShare 获取分享详情
func (u *UCService) getShare(shareID, stoken string) (*UCShareResult, error) {
	queryParams := map[string]string{
		"pr":            ucPR,
		"fr":             "pc",
		"uc_param_str":  "",
		"pwd_id":        shareID,
		"stoken":        stoken,
		"pdir_fid":      "0",
		"force":         "0",
		"_page":         "1",
		"_size":         "100",
		"_fetch_banner": "1",
		"_fetch_share":  "1",
		"_fetch_total":  "1",
		"_sort":         "file_type:asc,updated_at:desc",
	}

	respData, err := u.HTTPGet(ucAPIBase+"/share/sharepage/detail", queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int           `json:"status"`
		Message string        `json:"message"`
		Data    UCShareResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("%s", response.Message)
	}

	return &response.Data, nil
}

// getShareSave 转存分享到根目录
func (u *UCService) getShareSave(shareID, stoken string, fidList, fidTokenList []string) (*UCSaveResult, error) {
	return u.getShareSaveToDir(shareID, stoken, fidList, fidTokenList, "0")
}

// getShareSaveToDir 转存分享到指定目录
func (u *UCService) getShareSaveToDir(shareID, stoken string, fidList, fidTokenList []string, toPdirFid string) (*UCSaveResult, error) {
	data := map[string]interface{}{
		"pwd_id":         shareID,
		"stoken":         stoken,
		"fid_list":       fidList,
		"fid_token_list": fidTokenList,
		"to_pdir_fid":    toPdirFid,
	}

	respData, err := u.HTTPPost(ucAPIBase+"/share/sharepage/save", data, u.ucCommonQuery())
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int          `json:"status"`
		Message string       `json:"message"`
		Data    UCSaveResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("%s", response.Message)
	}

	return &response.Data, nil
}

// generateTimestamp 生成指定长度的时间戳（用于 task 轮询的 __t 参数）
func (u *UCService) generateTimestamp(length int) int64 {
	timestamp := utils.GetCurrentTime().UnixNano() / int64(time.Millisecond)
	timestampStr := strconv.FormatInt(timestamp, 10)
	if len(timestampStr) > length {
		timestampStr = timestampStr[:length]
	}
	timestamp, _ = strconv.ParseInt(timestampStr, 10, 64)
	return timestamp
}

// getShareBtn 创建分享
func (u *UCService) getShareBtn(fidList []string, title string) (*UCShareBtnResult, error) {
	data := map[string]interface{}{
		"fid_list":     fidList,
		"title":        title,
		"url_type":     1,
		"expired_type": 1, // 永久分享
	}

	respData, err := u.HTTPPost(ucAPIBase+"/share", data, u.ucCommonQuery())
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int              `json:"status"`
		Message string           `json:"message"`
		Data    UCShareBtnResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("%s", response.Message)
	}

	return &response.Data, nil
}

// getShareTask 获取任务状态
func (u *UCService) getShareTask(taskID string, retryIndex int) (*UCTaskResult, error) {
	queryParams := map[string]string{
		"pr":           ucPR,
		"fr":           "pc",
		"uc_param_str": "",
		"task_id":      taskID,
		"retry_index":  fmt.Sprintf("%d", retryIndex),
		"__dt":         "21192",
		"__t":          fmt.Sprintf("%d", u.generateTimestamp(13)),
	}

	respData, err := u.HTTPGet(ucAPIBase+"/task", queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int          `json:"status"`
		Message string       `json:"message"`
		Data    UCTaskResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("%s", response.Message)
	}

	return &response.Data, nil
}

// getSharePassword 获取分享密码与链接
func (u *UCService) getSharePassword(shareID string) (*UCPasswordResult, error) {
	data := map[string]interface{}{
		"share_id": shareID,
	}

	respData, err := u.HTTPPost(ucAPIBase+"/share/password", data, u.ucCommonQuery())
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int              `json:"status"`
		Message string           `json:"message"`
		Data    UCPasswordResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("%s", response.Message)
	}

	return &response.Data, nil
}

// waitForTask 等待任务完成
func (u *UCService) waitForTask(taskID string) (*UCTaskResult, error) {
	maxRetries := 50
	retryDelay := 2 * time.Second

	for retryIndex := 0; retryIndex < maxRetries; retryIndex++ {
		result, err := u.getShareTask(taskID, retryIndex)
		if err != nil {
			if strings.Contains(err.Error(), "capacity limit") {
				return nil, fmt.Errorf("UC账号容量不足")
			}
			return nil, err
		}

		if result.Status == 2 { // 任务完成
			return result, nil
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("UC转存任务超时")
}

// getDirFile 获取指定文件夹的文件列表
func (u *UCService) getDirFile(pdirFid string) ([]map[string]interface{}, error) {
	log.Printf("正在遍历UC父文件夹: %s", pdirFid)

	queryParams := map[string]string{
		"pr":              ucPR,
		"fr":              "pc",
		"uc_param_str":    "",
		"pdir_fid":        pdirFid,
		"_page":           "1",
		"_size":           "50",
		"_fetch_total":    "1",
		"_fetch_sub_dirs": "0",
		"_sort":           "updated_at:desc",
	}

	respData, err := u.HTTPGet(ucAPIBase+"/file/sort", queryParams)
	if err != nil {
		log.Printf("获取UC目录文件失败: %v", err)
		return nil, err
	}

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			List []map[string]interface{} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		log.Printf("解析UC目录文件响应失败: %v", err)
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("%s", response.Message)
	}

	return response.Data.List, nil
}

// ============================================================================
// 广告处理（复用 pan_ad.go 共享工具 + UC 自身的转存方法）
// ============================================================================

// deleteAdFiles 删除转存目录内命中广告关键词的文件
func (u *UCService) deleteAdFiles(pdirFid string) error {
	log.Printf("开始删除UC广告文件，目录ID: %s", pdirFid)

	fileList, err := u.getDirFile(pdirFid)
	if err != nil {
		return err
	}

	if len(fileList) == 0 {
		log.Printf("UC目录为空，无需删除广告文件")
		return nil
	}

	for _, file := range fileList {
		fileName, ok := file["file_name"].(string)
		if !ok {
			continue
		}
		if containsAdKeywords(fileName) { // pan_ad.go 共享函数
			if fid, ok := file["fid"].(string); ok {
				log.Printf("删除UC广告文件: %s (FID: %s)", fileName, fid)
				if _, err := u.DeleteFiles([]string{fid}); err != nil {
					log.Printf("删除UC广告文件失败: %v", err)
				}
			}
		}
	}

	return nil
}

// addAd 添加个人自定义广告到转存目录（随机选一条系统配置的广告）
func (u *UCService) addAd(dirID string) error {
	log.Printf("开始添加UC自定义广告到目录: %s", dirID)

	autoInsertAdStr, err := getAdSystemConfigValue(entity.ConfigKeyAutoInsertAd)
	if err != nil {
		log.Printf("获取自动插入广告配置失败: %v", err)
		return err
	}
	if autoInsertAdStr == "" {
		log.Printf("没有配置自动插入广告，跳过UC广告插入")
		return nil
	}

	adURLs := splitAdURLs(autoInsertAdStr)
	if len(adURLs) == 0 {
		return nil
	}
	adFileIDs := extractAdFileIDs(adURLs)
	if len(adFileIDs) == 0 {
		log.Printf("没有有效的UC广告文件ID，跳过广告插入")
		return nil
	}

	rand.Seed(utils.GetCurrentTimestampNano())
	selectedAdID := adFileIDs[rand.Intn(len(adFileIDs))]
	log.Printf("选择UC广告文件ID: %s", selectedAdID)

	stokenResult, err := u.getStoken(selectedAdID)
	if err != nil {
		return err
	}
	adDetail, err := u.getShare(selectedAdID, stokenResult.Stoken)
	if err != nil {
		return err
	}
	if len(adDetail.List) == 0 {
		return fmt.Errorf("UC广告文件详情为空")
	}

	adFile := adDetail.List[0]
	saveResult, err := u.getShareSaveToDir(selectedAdID, stokenResult.Stoken, []string{adFile.Fid}, []string{adFile.ShareFidToken}, dirID)
	if err != nil {
		return err
	}
	if _, err := u.waitForTask(saveResult.TaskID); err != nil {
		return err
	}

	log.Printf("UC广告文件添加成功")
	return nil
}

// ============================================================================
// GetUserInfo 账号信息查询
// ============================================================================

// GetUserInfo 获取用户信息
func (u *UCService) GetUserInfo(cookie *string) (*UserInfo, error) {
	// 设置Cookie
	u.SetHeader("Cookie", *cookie)

	// 调用UC网盘用户信息API
	userInfoURL := "https://drive.uc.cn/api/user/info"

	resp, err := u.HTTPGet(userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("获取UC用户信息失败: %v", err)
	}

	// 解析响应
	var result struct {
		Code int `json:"code"`
		Data struct {
			Username   string `json:"username"`
			Nickname   string `json:"nickname"`
			VipStatus  int    `json:"vip_status"`
			TotalSpace int64  `json:"total_space"`
			UsedSpace  int64  `json:"used_space"`
		} `json:"data"`
	}

	if err := u.ParseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("解析UC用户信息失败: %v", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("UC API返回错误: %d", result.Code)
	}

	// 转换VIP状态
	vipStatus := result.Data.VipStatus > 0

	// 使用nickname或username
	username := result.Data.Nickname
	if username == "" {
		username = result.Data.Username
	}
	if username == "" {
		username = "UC网盘用户"
	}

	return &UserInfo{
		Username:    username,
		VIPStatus:   vipStatus,
		UsedSpace:   result.Data.UsedSpace,
		TotalSpace:  result.Data.TotalSpace,
		ServiceType: "uc",
	}, nil
}

// GetUserInfoByEntity 根据 entity.Cks 获取用户信息
func (u *UCService) GetUserInfoByEntity(cks entity.Cks) (*UserInfo, error) {
	ck := cks.Ck
	if ck == "" {
		return nil, fmt.Errorf("UC账号Cookie为空")
	}
	return u.GetUserInfo(&ck)
}

// SetCKSRepository 设置CKS仓储（与夸克保持一致的空实现形态）
func (u *UCService) SetCKSRepository(cksRepo repo.CksRepository, entity entity.Cks) {
}

// ============================================================================
// UC 转存结果结构体（与夸克同构，加 UC 前缀避免同包命名冲突）
// ============================================================================

type UCStokenResult struct {
	Stoken string `json:"stoken"`
	Title  string `json:"title"`
}

type UCShareResult struct {
	Share struct {
		Title string `json:"title"`
	} `json:"share"`
	List []struct {
		Fid           string `json:"fid"`
		ShareFidToken string `json:"share_fid_token"`
	} `json:"list"`
}

type UCSaveResult struct {
	TaskID string `json:"task_id"`
}

type UCShareBtnResult struct {
	TaskID string `json:"task_id"`
}

type UCTaskResult struct {
	Status  int    `json:"status"`
	ShareID string `json:"share_id"`
	SaveAs  struct {
		SaveAsTopFids []string `json:"save_as_top_fids"`
	} `json:"save_as"`
}

type UCPasswordResult struct {
	ShareURL   string `json:"share_url"`
	ShareTitle string `json:"share_title"`
	Code       string `json:"code"`
	FirstFile  struct {
		Fid string `json:"fid"`
	} `json:"first_file"`
}
