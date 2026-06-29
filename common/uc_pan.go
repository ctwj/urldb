package pan

import (
	"encoding/json"
	"fmt"
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
//
// 日志约定：诊断日志统一用 utils.Debug（开发环境 DEBUG 级可见，生产 INFO 级自动屏蔽），
// 关键里程碑用 utils.Info，失败用 utils.Error。严禁打印 Cookie 明文。
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

	utils.Debug("[UC] 创建UC服务 cookieLen=%d apiBase=%s", len(cookie), ucAPIBase)
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

	isType := 0
	if config != nil {
		isType = config.IsType
	}
	utils.Info("[UC] 开始处理分享 shareID=%s isType=%d(0=转存,1=校验)", shareID, isType)

	// 获取 stoken
	var stoken string
	if config == nil || config.Stoken == "" {
		stokenResult, err := u.getStoken(shareID)
		if err != nil {
			utils.Error("[UC] 获取stoken失败 shareID=%s err=%v", shareID, err)
			return ErrorResult(fmt.Sprintf("获取stoken失败: %v", err)), nil
		}
		stoken = strings.ReplaceAll(stokenResult.Stoken, " ", "+")
	} else {
		stoken = strings.ReplaceAll(config.Stoken, " ", "+")
	}
	utils.Debug("[UC] 获取stoken成功 shareID=%s stokenLen=%d", shareID, len(stoken))

	// 获取分享详情
	shareResult, err := u.getShare(shareID, stoken)
	if err != nil || len(shareResult.List) == 0 {
		utils.Error("[UC] 获取分享详情失败 shareID=%s listLen=%d err=%v", shareID, len(shareResult.List), err)
		return ErrorResult(fmt.Sprintf("获取分享详情失败: %v", err)), nil
	}
	utils.Debug("[UC] 获取分享详情成功 shareID=%s title=%s fileCount=%d", shareID, shareResult.Share.Title, len(shareResult.List))

	// 检验模式：只读取资源信息，不实际转存（FR-007）
	if isType == 1 {
		utils.Info("[UC] 校验模式完成（不转存）shareID=%s title=%s", shareID, shareResult.Share.Title)
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
		utils.Error("[UC] 转存失败 shareID=%s err=%v", shareID, err)
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}
	utils.Debug("[UC] 转存请求成功 shareID=%s taskID=%s", shareID, saveResult.TaskID)

	// 等待转存完成
	myData, err := u.waitForTask(saveResult.TaskID)
	if err != nil {
		utils.Error("[UC] 等待转存完成失败 shareID=%s taskID=%s err=%v", shareID, saveResult.TaskID, err)
		return ErrorResult(fmt.Sprintf("等待转存完成失败: %v", err)), nil
	}

	if len(myData.SaveAs.SaveAsTopFids) == 0 {
		utils.Error("[UC] 转存完成但未获取到文件标识 shareID=%s", shareID)
		return ErrorResult("转存完成但未获取到文件标识"), nil
	}
	utils.Debug("[UC] 转存完成 shareID=%s topFids=%v", shareID, myData.SaveAs.SaveAsTopFids)

	// 删除广告文件（如果有配置）
	if err := u.deleteAdFiles(myData.SaveAs.SaveAsTopFids[0]); err != nil {
		utils.Debug("[UC] 删除广告文件失败（不阻断）fid=%s err=%v", myData.SaveAs.SaveAsTopFids[0], err)
	}

	// 添加个人自定义广告
	if err := u.addAd(myData.SaveAs.SaveAsTopFids[0]); err != nil {
		utils.Debug("[UC] 添加广告文件失败（不阻断）fid=%s err=%v", myData.SaveAs.SaveAsTopFids[0], err)
	}

	// 分享资源
	shareBtnResult, err := u.getShareBtn(myData.SaveAs.SaveAsTopFids, title)
	if err != nil {
		utils.Error("[UC] 创建再分享失败 shareID=%s err=%v", shareID, err)
		return ErrorResult(fmt.Sprintf("分享失败: %v", err)), nil
	}
	utils.Debug("[UC] 创建再分享请求成功 taskID=%s", shareBtnResult.TaskID)

	// 等待分享完成
	shareTaskResult, err := u.waitForTask(shareBtnResult.TaskID)
	if err != nil {
		utils.Error("[UC] 等待分享完成失败 taskID=%s err=%v", shareBtnResult.TaskID, err)
		return ErrorResult(fmt.Sprintf("等待分享完成失败: %v", err)), nil
	}

	// 获取分享密码
	passwordResult, err := u.getSharePassword(shareTaskResult.ShareID)
	if err != nil {
		utils.Error("[UC] 获取分享密码失败 shareID=%s err=%v", shareTaskResult.ShareID, err)
		return ErrorResult(fmt.Sprintf("获取分享密码失败: %v", err)), nil
	}

	// 确定 fid
	var fid string
	if len(myData.SaveAs.SaveAsTopFids) > 1 {
		fid = strings.Join(myData.SaveAs.SaveAsTopFids, ",")
	} else {
		fid = passwordResult.FirstFile.Fid
	}

	utils.Info("[UC] 转存成功 shareID=%s newShareUrl=%s title=%s fid=%s", shareID, passwordResult.ShareURL, passwordResult.ShareTitle, fid)
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
	utils.Debug("[UC] GetFiles pdirFid=%s", pdirFid)

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
		utils.Error("[UC] GetFiles 请求失败 pdirFid=%s err=%v", pdirFid, err)
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
		utils.Error("[UC] GetFiles 解析响应失败 pdirFid=%s err=%v bodyHead=%s", pdirFid, err, headSnippet(data))
		return ErrorResult("解析UC响应失败"), nil
	}

	if response.Status != 200 {
		message := response.Message
		if strings.Contains(message, "require login") || strings.Contains(message, "guest") {
			message = "UC未登录，请检查cookie"
		}
		utils.Error("[UC] GetFiles 接口返回非200 pdirFid=%s status=%d msg=%s", pdirFid, response.Status, message)
		return ErrorResult(message), nil
	}

	utils.Debug("[UC] GetFiles 成功 pdirFid=%s count=%d", pdirFid, len(response.Data.List))
	return SuccessResult("获取成功", response.Data.List), nil
}

// DeleteFiles 删除文件
func (u *UCService) DeleteFiles(fileList []string) (*TransferResult, error) {
	if len(fileList) == 0 {
		return ErrorResult("文件列表为空"), nil
	}
	utils.Debug("[UC] DeleteFiles count=%d fids=%v", len(fileList), fileList)

	for _, fileID := range fileList {
		if err := u.deleteSingleFile(fileID); err != nil {
			utils.Error("[UC] 删除UC文件失败 fid=%s err=%v", fileID, err)
			return ErrorResult(fmt.Sprintf("删除UC文件 %s 失败: %v", fileID, err)), nil
		}
	}

	utils.Debug("[UC] DeleteFiles 成功 count=%d", len(fileList))
	return SuccessResult("删除成功", nil), nil
}

// deleteSingleFile 删除单个文件
func (u *UCService) deleteSingleFile(fileID string) error {
	utils.Debug("[UC] 删除文件 fid=%s", fileID)

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
		return fmt.Errorf("解析UC删除响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}

	if response.Status != 200 {
		return fmt.Errorf("删除UC文件失败: status=%d msg=%s", response.Status, response.Message)
	}

	if response.Data.TaskID != "" {
		utils.Debug("[UC] 删除文件任务提交 fid=%s taskID=%s", fileID, response.Data.TaskID)
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
	passcode := u.configCode()
	data := map[string]interface{}{
		"passcode": passcode, // 加密分享的提取码（公开分享为空）
		"pwd_id":   shareID,
	}
	utils.Debug("[UC] 请求 stoken shareID=%s passcodeLen=%d", shareID, len(passcode))

	respData, err := u.HTTPPost(ucAPIBase+"/share/sharepage/token", data, u.ucCommonQuery())
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	var response struct {
		Status  int            `json:"status"`
		Message string         `json:"message"`
		Data    UCStokenResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("获取stoken失败: status=%d msg=%s", response.Status, response.Message)
	}

	utils.Debug("[UC] stoken 响应成功 shareID=%s stokenLen=%d title=%s", shareID, len(response.Data.Stoken), response.Data.Title)
	return &response.Data, nil
}

// getShare 获取分享详情
func (u *UCService) getShare(shareID, stoken string) (*UCShareResult, error) {
	queryParams := map[string]string{
		"pr":             ucPR,
		"fr":             "pc",
		"uc_param_str":   "",
		"pwd_id":         shareID,
		"stoken":         stoken,
		"pdir_fid":       "0",
		"force":          "0",
		"_page":          "1",
		"_size":          "100",
		"_fetch_banner":  "1",
		"_fetch_share":   "1",
		"_fetch_total":   "1",
		"_sort":          "file_type:asc,updated_at:desc",
	}

	respData, err := u.HTTPGet(ucAPIBase+"/share/sharepage/detail", queryParams)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	var response struct {
		Status  int           `json:"status"`
		Message string        `json:"message"`
		Data    UCShareResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("获取分享详情失败: status=%d msg=%s", response.Status, response.Message)
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
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	var response struct {
		Status  int          `json:"status"`
		Message string       `json:"message"`
		Data    UCSaveResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("转存失败: status=%d msg=%s", response.Status, response.Message)
	}

	utils.Debug("[UC] save 响应成功 taskID=%s toPdirFid=%s fileCount=%d", response.Data.TaskID, toPdirFid, len(fidList))
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
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	var response struct {
		Status  int              `json:"status"`
		Message string           `json:"message"`
		Data    UCShareBtnResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("创建分享失败: status=%d msg=%s", response.Status, response.Message)
	}

	utils.Debug("[UC] share 创建成功 taskID=%s title=%s fidCount=%d", response.Data.TaskID, title, len(fidList))
	return &response.Data, nil
}

// getShareTask 获取任务状态（高频轮询调用，内部不打日志，避免噪声；由 waitForTask 汇总）
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
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	var response struct {
		Status  int          `json:"status"`
		Message string       `json:"message"`
		Data    UCTaskResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v bodyHead=%s", err, headSnippet(respData))
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
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	var response struct {
		Status  int              `json:"status"`
		Message string           `json:"message"`
		Data    UCPasswordResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("获取分享密码失败: status=%d msg=%s", response.Status, response.Message)
	}

	utils.Debug("[UC] share/password 成功 shareUrl=%s hasCode=%v", response.Data.ShareURL, response.Data.Code != "")
	return &response.Data, nil
}

// waitForTask 等待任务完成
func (u *UCService) waitForTask(taskID string) (*UCTaskResult, error) {
	maxRetries := 50
	retryDelay := 2 * time.Second
	utils.Debug("[UC] 开始轮询任务 taskID=%s maxRetries=%d", taskID, maxRetries)

	for retryIndex := 0; retryIndex < maxRetries; retryIndex++ {
		result, err := u.getShareTask(taskID, retryIndex)
		if err != nil {
			utils.Debug("[UC] 任务查询返回错误 taskID=%s retry=%d rawErr=%v", taskID, retryIndex, err)
			if strings.Contains(err.Error(), "capacity limit") {
				utils.Error("[UC] 任务失败-容量不足 taskID=%s", taskID)
				return nil, fmt.Errorf("UC账号容量不足")
			}
			return nil, err
		}

		if result.Status == 2 { // 任务完成
			utils.Debug("[UC] 任务完成 taskID=%s retry=%d", taskID, retryIndex)
			return result, nil
		}

		time.Sleep(retryDelay)
	}

	utils.Error("[UC] 任务轮询超时 taskID=%s", taskID)
	return nil, fmt.Errorf("UC转存任务超时")
}

// getDirFile 获取指定文件夹的文件列表
func (u *UCService) getDirFile(pdirFid string) ([]map[string]interface{}, error) {
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
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			List []map[string]interface{} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}

	if response.Status != 200 {
		return nil, fmt.Errorf("获取目录文件失败: status=%d msg=%s", response.Status, response.Message)
	}

	utils.Debug("[UC] getDirFile 成功 pdirFid=%s count=%d", pdirFid, len(response.Data.List))
	return response.Data.List, nil
}

// ============================================================================
// 广告处理（复用 pan_ad.go 共享工具 + UC 自身的转存方法）
// ============================================================================

// deleteAdFiles 删除转存目录内命中广告关键词的文件
func (u *UCService) deleteAdFiles(pdirFid string) error {
	fileList, err := u.getDirFile(pdirFid)
	if err != nil {
		return err
	}

	if len(fileList) == 0 {
		utils.Debug("[UC] 目录为空，无广告文件可删 pdirFid=%s", pdirFid)
		return nil
	}

	deleted := 0
	for _, file := range fileList {
		fileName, ok := file["file_name"].(string)
		if !ok {
			continue
		}
		if containsAdKeywords(fileName) { // pan_ad.go 共享函数
			if fid, ok := file["fid"].(string); ok {
				utils.Debug("[UC] 命中广告关键词，删除 fileName=%s fid=%s", fileName, fid)
				if _, err := u.DeleteFiles([]string{fid}); err != nil {
					utils.Debug("[UC] 删除广告文件失败（继续）fileName=%s err=%v", fileName, err)
				} else {
					deleted++
				}
			}
		}
	}
	utils.Debug("[UC] 广告清理完成 pdirFid=%s deleted=%d", pdirFid, deleted)
	return nil
}

// addAd 添加个人自定义广告到转存目录（随机选一条系统配置的广告）
func (u *UCService) addAd(dirID string) error {
	autoInsertAdStr, err := getAdSystemConfigValue(entity.ConfigKeyAutoInsertAd)
	if err != nil {
		return err
	}
	if autoInsertAdStr == "" {
		utils.Debug("[UC] 未配置自动插入广告，跳过")
		return nil
	}

	adURLs := splitAdURLs(autoInsertAdStr)
	if len(adURLs) == 0 {
		return nil
	}
	adFileIDs := extractAdFileIDs(adURLs)
	if len(adFileIDs) == 0 {
		utils.Debug("[UC] 无有效广告文件ID，跳过 adURLCount=%d", len(adURLs))
		return nil
	}

	rand.Seed(utils.GetCurrentTimestampNano())
	selectedAdID := adFileIDs[rand.Intn(len(adFileIDs))]
	utils.Debug("[UC] 选中广告文件 adShareID=%s dirID=%s", selectedAdID, dirID)

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

	utils.Debug("[UC] 广告插入成功 adShareID=%s dirID=%s", selectedAdID, dirID)
	return nil
}

// ============================================================================
// GetUserInfo 账号信息查询
// ============================================================================

// GetUserInfo 获取用户信息（容量 + 会员状态）。
//
// 走 UC 已验证可用的 /member 接口（主机 pc-api.uc.cn/1/clouddrive，pr=UCBrowser），
// 取 total_capacity / use_capacity / member_type。
//
// 注意：早期实现的 https://drive.uc.cn/api/user/info 为无效端点（302 跳转回首页 HTML，
// 导致 "invalid character '<'" 解析错误），已弃用。/member 不返回昵称，而 UC 无公开的
// 昵称接口（/user/info 返回 404），故 Username 使用占位值。
func (u *UCService) GetUserInfo(cookie *string) (*UserInfo, error) {
	// 设置Cookie
	u.SetHeader("Cookie", *cookie)
	utils.Debug("[UC] GetUserInfo 请求 /member cookieLen=%d", len(*cookie))

	queryParams := map[string]string{
		"pr":              ucPR,
		"fr":              "pc",
		"uc_param_str":    "",
		"fetch_subscribe": "false",
		"_ch":             "home",
		"fetch_identity":  "false",
	}

	resp, err := u.HTTPGet(ucAPIBase+"/member", queryParams)
	if err != nil {
		return nil, fmt.Errorf("获取UC用户信息失败: %v", err)
	}

	var result struct {
		Status  int    `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			MemberType    string `json:"member_type"`
			UseCapacity   int64  `json:"use_capacity"`
			TotalCapacity int64  `json:"total_capacity"`
		} `json:"data"`
	}

	if err := u.ParseJSONResponse(resp, &result); err != nil {
		utils.Error("[UC] GetUserInfo 解析失败 err=%v bodyHead=%s", err, headSnippet(resp))
		return nil, fmt.Errorf("解析UC用户信息失败: %v", err)
	}

	if result.Status != 200 || result.Code != 0 {
		utils.Error("[UC] GetUserInfo 接口返回错误 status=%d code=%d msg=%s", result.Status, result.Code, result.Message)
		return nil, fmt.Errorf("UC接口返回错误: status=%d code=%d msg=%s", result.Status, result.Code, result.Message)
	}

	// 会员判定：member_type 非 NORMAL 即视为 VIP（与夸克一致）
	vipStatus := result.Data.MemberType != "NORMAL" && result.Data.MemberType != ""

	utils.Debug("[UC] GetUserInfo 成功 memberType=%s total=%d used=%d vip=%v", result.Data.MemberType, result.Data.TotalCapacity, result.Data.UseCapacity, vipStatus)
	return &UserInfo{
		Username:    "UC网盘用户",
		VIPStatus:   vipStatus,
		UsedSpace:   result.Data.UseCapacity,
		TotalSpace:  result.Data.TotalCapacity,
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

// headSnippet 返回响应体的前 200 字符（用于解析失败时定位是否返回 HTML/错误页），
// 仅截断展示，避免打印过长的响应体；Cookie 不会出现在响应体中。
func headSnippet(data []byte) string {
	n := len(data)
	if n > 200 {
		n = 200
	}
	return string(data[:n])
}
