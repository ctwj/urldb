package pan

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

// CaptchaData 存储在数据库中的验证码令牌数据
type CaptchaData struct {
	CaptchaToken string `json:"captcha_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// XunleiExtraData 所有额外数据的容器
type XunleiTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
	Sub          string `json:"sub"`
	TokenType    string `json:"token_type"`
	UserId       string `json:"user_id"`
}

type XunleiExtraData struct {
	captcha *CaptchaData
	token   *XunleiTokenData
}

type XunleiPanService struct {
	*BasePanService
	configMutex sync.RWMutex
	clientId    string
	deviceId    string
	entity      entity.Cks
	cksRepo     repo.CksRepository
	extra       XunleiExtraData // 需要保存到数据库的token信息
}

var (
	xunleiInstance *XunleiPanService
	xunleiOnce     sync.Once
)

// 配置化 API Host
func (x *XunleiPanService) apiHost(apiType string) string {
	if apiType == "user" {
		return "https://xluser-ssl.xunlei.com"
	}
	return "https://api-pan.xunlei.com"
}

func (x *XunleiPanService) setCommonHeader(req *http.Request) {
	for k, v := range x.headers {
		req.Header.Set(k, v)
	}
}

// NewXunleiPanService 创建迅雷网盘服务
func NewXunleiPanService(config *PanConfig) *XunleiPanService {
	xunleiOnce.Do(func() {
		xunleiInstance = &XunleiPanService{
			BasePanService: NewBasePanService(config),
			clientId:       "Xqp0kJBXWhwaTpB6",
			deviceId:       "925b7631473a13716b791d7f28289cad",
			extra:          XunleiExtraData{}, // Initialize extra with zero values
		}

		xunleiInstance.SetHeaders(map[string]string{
			"Accept":             "*/;",
			"Accept-Encoding":    "deflate",
			"Accept-Language":    "zh-CN,zh;q=0.9",
			"Cache-Control":      "no-cache",
			"Content-Type":       "application/json",
			"Origin":             "https://pan.xunlei.com",
			"Pragma":             "no-cache",
			"Priority":           "u=1,i",
			"Referer":            "https://pan.xunlei.com/",
			"sec-ch-ua":          `"Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"`,
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": `"Windows"`,
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-site",
			"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
			"Authorization":      "",
			"x-captcha-token":    "",
			"x-client-id":        xunleiInstance.clientId,
			"x-device-id":        xunleiInstance.deviceId,
		})
	})
	xunleiInstance.UpdateConfig(config)
	return xunleiInstance
}

// SetCKSRepository 设置 CksRepository 和 entity
func (x *XunleiPanService) SetCKSRepository(cksRepo repo.CksRepository, entity entity.Cks) {
	x.cksRepo = cksRepo
	x.entity = entity
	var extra XunleiExtraData
	if err := json.Unmarshal([]byte(x.entity.Extra), &extra); err != nil {
		log.Printf("解析 extra 数据失败: %v，使用空数据", err)
	}
	x.extra = extra
}

// GetXunleiInstance 获取迅雷网盘服务单例实例
func GetXunleiInstance() *XunleiPanService {
	return NewXunleiPanService(nil)
}

func (x *XunleiPanService) GetAccessTokenByRefreshToken(refreshToken string) (map[string]interface{}, error) {
	// 构造请求体
	body := map[string]interface{}{
		"client_id":     x.clientId,
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	// 过滤 headers（移除 Authorization 和 x-captcha-token）
	filteredHeaders := make(map[string]string)
	for k, v := range x.headers {
		if k != "Authorization" && k != "x-captcha-token" {
			filteredHeaders[k] = v
		}
	}

	// 调用 API 获取新的 token
	resp, err := x.requestXunleiApi("https://xluser-ssl.xunlei.com/v1/auth/token", "POST", body, nil, filteredHeaders)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("获取 access_token 请求失败: %v", err)
	}

	// 正确做法：用 exists 判断
	if _, exists := resp["access_token"]; exists {
		// 会输出，即使值为 nil
	} else {
		return map[string]interface{}{}, fmt.Errorf("获取 access_token 请求失败: %v 不存在", "access_token")
	}

	// 计算过期时间（当前时间 + expires_in - 60 秒缓冲）
	currentTime := time.Now().Unix()
	expiresAt := currentTime + int64(resp["expires_in"].(float64)) - 60
	resp["expires_at"] = expiresAt
	return resp, nil
}

// getAccessToken 获取 Access Token（内部包含缓存判断、刷新、保存）- 匹配 PHP 版本
func (x *XunleiPanService) getAccessToken() (string, error) {
	// 检查 Access Token 是否有效
	currentTime := time.Now().Unix()
	if x.extra.token != nil && x.extra.token.AccessToken != "" && x.extra.token.ExpiresAt > currentTime {
		return x.extra.token.AccessToken, nil
	}
	newData, err := x.GetAccessTokenByRefreshToken(x.extra.token.RefreshToken)

	if err != nil {
		return "", fmt.Errorf("获取 access_token 失败: %v", err)
	}
	newAccessToken := newData["access_token"].(string)
	x.extra.token.AccessToken = newAccessToken
	x.extra.token.ExpiresAt = int64(newData["expires_in"].(float64)) - 60

	// 保存到数据库
	extraBytes, err := json.Marshal(x.extra.token)
	if err != nil {
		return "", fmt.Errorf("序列化 extra 数据失败: %v", err)
	}
	x.entity.Extra = string(extraBytes)
	if err := x.cksRepo.UpdateWithAllFields(&x.entity); err != nil {
		return "", fmt.Errorf("保存 access_token 到数据库失败: %v", err)
	}
	return newAccessToken, nil
}

// getCaptchaToken 获取 captcha_token - 匹配 PHP 版本
func (x *XunleiPanService) getCaptchaToken() (string, error) {
	// 检查 Captcha Token 是否有效
	currentTime := time.Now().Unix()
	if x.extra.captcha != nil && x.extra.captcha.CaptchaToken != "" && x.extra.captcha.ExpiresAt > currentTime {
		return x.extra.captcha.CaptchaToken, nil
	}

	// 构造请求体
	body := map[string]interface{}{
		"client_id": x.clientId,
		"action":    "get:/drive/v1/share",
		"device_id": x.deviceId,
		"meta": map[string]interface{}{
			"username":       "",
			"phone_number":   "",
			"email":          "",
			"package_name":   "pan.xunlei.com",
			"client_version": "1.45.0",
			"captcha_sign":   "1.fe2108ad808a74c9ac0243309242726c",
			"timestamp":      "1645241033384",
			"user_id":        "0",
		},
	}

	captchaHeaders := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}

	// 调用 API 获取 captcha_token
	resp, err := x.requestXunleiApi("https://xluser-ssl.xunlei.com/v1/shield/captcha/init", "POST", body, nil, captchaHeaders)
	if err != nil {
		return "", fmt.Errorf("获取 captcha_token 请求失败: %v", err)
	}

	if resp["code"] != 0 || resp["data"] == nil {
		return "", fmt.Errorf("获取 captcha_token 失败: %v", resp)
	}

	data := resp["data"].(map[string]interface{})
	captchaToken := data["captcha_token"].(string)

	// 计算过期时间（当前时间 + expires_in - 10 秒缓冲）
	expiresAt := currentTime + int64(data["expires_in"].(float64)) - 10

	// 更新 extra 数据
	if x.extra.captcha == nil {
		x.extra.captcha = &CaptchaData{}
	}
	x.extra.captcha.CaptchaToken = captchaToken
	x.extra.captcha.ExpiresAt = expiresAt

	// 保存到数据库
	extraBytes, err := json.Marshal(x.extra)
	if err != nil {
		return "", fmt.Errorf("序列化 extra 数据失败: %v", err)
	}
	x.entity.Extra = string(extraBytes)
	if err := x.cksRepo.UpdateWithAllFields(&x.entity); err != nil {
		return "", fmt.Errorf("保存 captcha_token 到数据库失败: %v", err)
	}

	return captchaToken, nil
}

// requestXunleiApi 迅雷 API 通用请求方法 - 使用 BasePanService 方法
func (x *XunleiPanService) requestXunleiApi(url string, method string, data map[string]interface{}, queryParams map[string]string, headers map[string]string) (map[string]interface{}, error) {
	var respData []byte
	var err error

	// 先更新当前请求的 headers
	originalHeaders := make(map[string]string)
	for k, v := range x.headers {
		originalHeaders[k] = v
	}

	// 临时设置请求的 headers
	for k, v := range headers {
		x.SetHeader(k, v)
	}
	defer func() {
		// 恢复原始 headers
		for k, v := range originalHeaders {
			x.SetHeader(k, v)
		}
	}()

	// 根据方法调用相应的 BasePanService 方法
	if method == "GET" {
		respData, err = x.HTTPGet(url, queryParams)
	} else if method == "POST" {
		respData, err = x.HTTPPost(url, data, queryParams)
	} else {
		return nil, fmt.Errorf("不支持的HTTP方法: %s", method)
	}

	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v, raw: %s", err, string(respData))
	}

	return result, nil
}

func (x *XunleiPanService) UpdateConfig(config *PanConfig) {
	if config == nil {
		return
	}
	x.configMutex.Lock()
	defer x.configMutex.Unlock()
	x.config = config
	if config.Cookie != "" {
		x.SetHeader("Cookie", config.Cookie)
	}
}

// GetServiceType 获取服务类型
func (x *XunleiPanService) GetServiceType() ServiceType {
	return Xunlei
}

// Transfer 转存分享链接 - 实现 PanService 接口，匹配 XunleiPan.php 的逻辑
func (x *XunleiPanService) Transfer(shareID string) (*TransferResult, error) {
	// 读取配置（线程安全）
	x.configMutex.RLock()
	config := x.config
	x.configMutex.RUnlock()

	log.Printf("开始处理迅雷分享: %s", shareID)

	// 1️⃣ 获取 AccessToken 和 CaptchaToken
	accessToken, err := x.getAccessToken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取accessToken失败: %v", err)), nil
	}

	captchaToken, err := x.getCaptchaToken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取captchaToken失败: %v", err)), nil
	}

	// 检查是否为检验模式
	if config.IsType == 1 {
		// 检验模式：直接获取分享信息
		urls := map[string]interface{}{
			"title":     "",
			"share_url": config.URL,
			"stoken":    "",
		}
		return SuccessResult("检验成功", urls), nil
	}

	// 转存模式：实现完整的转存流程
	shareID = strings.TrimRight(shareID, "#/")
	thisCode := strings.TrimLeft(x.config.Code, "#")

	// 获取分享详情
	shareDetail, err := x.getShare(shareID, thisCode, accessToken, captchaToken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取分享详情失败: %v", err)), nil
	}

	if shareDetail["code"].(int) != 200 {
		message := "获取分享详情失败"
		if shareDetail["message"] != nil {
			message = shareDetail["message"].(string)
		}
		return ErrorResult(message), nil
	}

	shareData := shareDetail["data"].(map[string]interface{})
	files := shareData["files"].([]interface{})

	// 转存到网盘
	parent_id := "0" // 默认存储路径
	if config.ExpiredType == 2 {
		parent_id = "指定路径" // 这里可能需要从配置获取
	}

	fileIDs := make([]string, 0)
	for _, file := range files {
		fileMap := file.(map[string]interface{})
		if fid, ok := fileMap["id"].(string); ok {
			fileIDs = append(fileIDs, fid)
		}
	}

	if len(fileIDs) == 0 {
		return ErrorResult("分享中没有可转存的文件"), nil
	}

	// 转存资源
	restoreResult, err := x.getRestore(shareID, shareData, accessToken, captchaToken, parent_id)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}

	if restoreResult["code"].(int) != 200 {
		return ErrorResult("转存失败"), nil
	}

	// 获取转存任务信息
	restoreData := restoreResult["data"].(map[string]interface{})
	taskID := restoreData["restore_task_id"].(string)

	// 等待转存完成
	_, err = x.waitForTask(taskID, accessToken, captchaToken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("等待转存完成失败: %v", err)), nil
	}

	// 获取任务结果并解析文件ID
	// 获取任务结果以获取文件ID
	existingFileIds := make([]string, 0)
	taskResp, err := x.getTasks(taskID, accessToken, captchaToken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取任务结果失败: %v", err)), nil
	}

	if taskData, ok := taskResp["data"].(map[string]interface{}); ok {
		if params, ok2 := taskData["params"].(map[string]interface{}); ok2 {
			if traceIds, ok3 := params["trace_file_ids"].(string); ok3 {
				traceData := make(map[string]interface{})
				json.Unmarshal([]byte(traceIds), &traceData)
				for _, fid := range traceData {
					existingFileIds = append(existingFileIds, fid.(string))
				}
			}
		}
	}

	// 处理广告过滤（这里简化处理）
	// TODO: 添加广告文件过滤逻辑

	// 创建分享链接
	expirationDays := "-1"
	if config.ExpiredType == 2 {
		expirationDays = "2"
	}

	// 根据share_id获取到分享链接
	passwordResult, err := x.getSharePassword(existingFileIds, accessToken, captchaToken, expirationDays)
	if err != nil {
		return ErrorResult(fmt.Sprintf("创建分享链接失败: %v", err)), nil
	}

	if passwordResult["code"].(int) != 200 {
		return ErrorResult("创建分享链接失败"), nil
	}

	passwordData := passwordResult["data"].(map[string]interface{})
	shareTitle := ""
	if len(files) > 0 {
		file0 := files[0].(map[string]interface{})
		if name, ok := file0["name"].(string); ok {
			shareTitle = name
		}
	}

	result := map[string]interface{}{
		"title":     shareTitle,
		"share_url": passwordData["share_url"].(string) + "?pwd=" + passwordData["pass_code"].(string),
		"code":      passwordData["pass_code"].(string),
		"fid":       existingFileIds,
	}

	return SuccessResult("转存成功", result), nil
}

// waitForTask 等待任务完成 - 使用 HTTPGet 方法
func (x *XunleiPanService) waitForTask(taskID string, accessToken, captchaToken string) (*XLTaskResult, error) {
	maxRetries := 50
	retryDelay := 2 * time.Second

	for retryIndex := 0; retryIndex < maxRetries; retryIndex++ {
		result, err := x.getTaskStatus(taskID, retryIndex, accessToken, captchaToken)
		if err != nil {
			return nil, err
		}

		if result.Status == 2 { // 任务完成
			return result, nil
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("任务超时")
}

// getTaskStatus 获取任务状态 - 使用 HTTPGet 方法
func (x *XunleiPanService) getTaskStatus(taskID string, retryIndex int, accessToken, captchaToken string) (*XLTaskResult, error) {
	apiURL := x.apiHost("") + "/drive/v1/task"
	queryParams := map[string]string{
		"task_id":     taskID,
		"retry_index": fmt.Sprintf("%d", retryIndex),
	}

	// 设置 request 所需的 headers
	headers := map[string]string{
		"Authorization":   "Bearer " + accessToken,
		"x-captcha-token": captchaToken,
	}

	resp, err := x.requestXunleiApi(apiURL, "GET", nil, queryParams, headers)
	if err != nil {
		return nil, err
	}

	if code, ok := resp["code"].(float64); ok && code != 200 {
		return nil, fmt.Errorf("获取任务状态失败")
	}

	var data XLTaskResult
	resultBytes, _ := json.Marshal(resp)
	if err := json.Unmarshal(resultBytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// GetUserInfoByEntity 根据 entity.Cks 获取用户信息（待实现）
func (x *XunleiPanService) GetUserInfoByEntity(cks entity.Cks) (*UserInfo, error) {
	return nil, nil
}

// getShare 获取分享详情 - 匹配 PHP 版本
func (x *XunleiPanService) getShare(shareID, passCode, accessToken, captchaToken string) (map[string]interface{}, error) {
	// 设置 headers
	headers := make(map[string]string)
	for k, v := range x.headers {
		headers[k] = v
	}
	headers["Authorization"] = "Bearer " + accessToken
	headers["x-captcha-token"] = captchaToken

	queryParams := map[string]string{
		"share_id":        shareID,
		"pass_code":       passCode,
		"limit":           "100",
		"pass_code_token": "",
		"page_token":      "",
		"thumbnail_size":  "SIZE_SMALL",
	}

	return x.requestXunleiApi("https://api-pan.xunlei.com/drive/v1/share", "GET", nil, queryParams, headers)
}

// getRestore 转存到网盘 - 匹配 PHP 版本
func (x *XunleiPanService) getRestore(shareID string, infoData map[string]interface{}, accessToken, captchaToken, parentID string) (map[string]interface{}, error) {
	ids := make([]string, 0)
	if files, ok := infoData["files"].([]interface{}); ok {
		for _, file := range files {
			if fileMap, ok2 := file.(map[string]interface{}); ok2 {
				if id, ok3 := fileMap["id"].(string); ok3 {
					ids = append(ids, id)
				}
			}
		}
	}

	passCodeToken := ""
	if token, ok := infoData["pass_code_token"]; ok {
		if tokenStr, ok2 := token.(string); ok2 {
			passCodeToken = tokenStr
		}
	}

	data := map[string]interface{}{
		"parent_id":         parentID,
		"share_id":          shareID,
		"pass_code_token":   passCodeToken,
		"ancestor_ids":      []string{},
		"specify_parent_id": true,
		"file_ids":          ids,
	}

	headers := make(map[string]string)
	for k, v := range x.headers {
		headers[k] = v
	}
	headers["Authorization"] = "Bearer " + accessToken
	headers["x-captcha-token"] = captchaToken

	return x.requestXunleiApi("https://api-pan.xunlei.com/drive/v1/share/restore", "POST", data, nil, headers)
}

// getTasks 获取转存任务状态 - 匹配 PHP 版本
func (x *XunleiPanService) getTasks(taskID, accessToken, captchaToken string) (map[string]interface{}, error) {
	headers := make(map[string]string)
	for k, v := range x.headers {
		headers[k] = v
	}
	headers["Authorization"] = "Bearer " + accessToken
	headers["x-captcha-token"] = captchaToken

	return x.requestXunleiApi("https://api-pan.xunlei.com/drive/v1/tasks/"+taskID, "GET", nil, nil, headers)
}

// getSharePassword 创建分享链接 - 匹配 PHP 版本
func (x *XunleiPanService) getSharePassword(fileIDs []string, accessToken, captchaToken, expirationDays string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"file_ids": fileIDs,
		"share_to": "copy",
		"params": map[string]interface{}{
			"subscribe_push":     "false",
			"WithPassCodeInLink": "true",
		},
		"title":           "云盘资源分享",
		"restore_limit":   "-1",
		"expiration_days": expirationDays,
	}

	headers := make(map[string]string)
	for k, v := range x.headers {
		headers[k] = v
	}
	headers["Authorization"] = "Bearer " + accessToken
	headers["x-captcha-token"] = captchaToken

	return x.requestXunleiApi("https://api-pan.xunlei.com/drive/v1/share", "POST", data, nil, headers)
}

// getShareInfo 获取分享信息（用于检验模式）
func (x *XunleiPanService) getShareInfo(shareID string) (*XLShareInfo, error) {
	// 使用现有的 GetShareFolder 方法获取分享信息
	shareDetail, err := x.GetShareFolder(shareID, "", "")
	if err != nil {
		return nil, err
	}

	// 构造分享信息
	shareInfo := &XLShareInfo{
		ShareID: shareID,
		Title:   fmt.Sprintf("迅雷分享_%s", shareID),
		Files:   make([]XLFileInfo, 0),
	}

	// 处理文件信息
	for _, file := range shareDetail.Data.Files {
		shareInfo.Files = append(shareInfo.Files, XLFileInfo{
			FileID: file.FileID,
			Name:   file.Name,
		})
	}

	return shareInfo, nil
}

// GetFiles 获取文件列表 - 匹配 PHP 版本接口调用
func (x *XunleiPanService) GetFiles(pdirFid string) (*TransferResult, error) {
	log.Printf("开始获取迅雷网盘文件列表，目录ID: %s", pdirFid)

	// 获取 tokens
	accessToken, err := x.getAccessToken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取accessToken失败: %v", err)), nil
	}

	captchaToken, err := x.getCaptchaToken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取captchaToken失败: %v", err)), nil
	}

	// 设置 headers
	headers := make(map[string]string)
	for k, v := range x.headers {
		headers[k] = v
	}
	headers["Authorization"] = "Bearer " + accessToken
	headers["x-captcha-token"] = captchaToken

	filters := map[string]interface{}{
		"phase": map[string]interface{}{
			"eq": "PHASE_TYPE_COMPLETE",
		},
		"trashed": map[string]interface{}{
			"eq": false,
		},
	}

	filtersStr, _ := json.Marshal(filters)
	queryParams := map[string]string{
		"parent_id":      pdirFid,
		"filters":        string(filtersStr),
		"with_audit":     "true",
		"thumbnail_size": "SIZE_SMALL",
		"limit":          "50",
	}

	result, err := x.requestXunleiApi("https://api-pan.xunlei.com/drive/v1/files", "GET", nil, queryParams, headers)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取文件列表失败: %v", err)), nil
	}

	if code, ok := result["code"].(float64); ok && code != 0 {
		return ErrorResult("获取文件列表失败"), nil
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		if files, ok2 := data["files"]; ok2 {
			return SuccessResult("获取成功", files), nil
		}
	}

	return SuccessResult("获取成功", []interface{}{}), nil
}

// DeleteFiles 删除文件 - 实现 PanService 接口
func (x *XunleiPanService) DeleteFiles(fileList []string) (*TransferResult, error) {
	log.Printf("开始删除迅雷网盘文件，文件数量: %d", len(fileList))

	// 使用现有的 ShareBatchDelete 方法删除分享
	result, err := x.ShareBatchDelete(fileList)
	if err != nil {
		return ErrorResult(fmt.Sprintf("删除文件失败: %v", err)), nil
	}

	if result.Code != 0 {
		return ErrorResult(fmt.Sprintf("删除文件失败: %s", result.Msg)), nil
	}

	return SuccessResult("删除成功", nil), nil
}

// GetUserInfo 获取用户信息 - 实现 PanService 接口，cookie 参数为 refresh_token，先获取 access_token 再访问 API
func (x *XunleiPanService) GetUserInfo(cookie string) (*UserInfo, error) {
	log.Printf("开始获取迅雷网盘用户信息（cookie 为 refresh_token）")

	// 使用 refresh_token 获取 access_token
	accessTokenData, err := x.GetAccessTokenByRefreshToken(cookie)
	if err != nil {
		return nil, fmt.Errorf("获取 access_token 失败: %v", err)
	}

	accessToken := accessTokenData["access_token"].(string)
	log.Printf("成功获取 access_token")

	// 设置 Authorization header 为 Bearer token
	x.SetHeader("Authorization", "Bearer "+accessToken)

	// 获取用户信息
	respData, err := x.HTTPGet(x.apiHost("user")+"/v1/user/me", nil)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	var response struct {
		Username string `json:"name"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	var userInfo *UserInfo
	userInfo = &UserInfo{
		Username:    response.Username,
		ServiceType: "xunlei",
	}

	return userInfo, nil
}

// GetShareList 严格对齐 GET + query（使用 BasePanService）
func (x *XunleiPanService) GetShareList(pageToken string) (*XLShareListResp, error) {
	api := x.apiHost("") + "/drive/v1/share/list"
	queryParams := map[string]string{
		"limit":          "100",
		"thumbnail_size": "SIZE_SMALL",
	}
	if pageToken != "" {
		queryParams["page_token"] = pageToken
	}

	respData, err := x.HTTPGet(api, queryParams)
	if err != nil {
		return nil, fmt.Errorf("获取分享列表失败: %v", err)
	}

	var data XLShareListResp
	if err := json.Unmarshal(respData, &data); err != nil {
		return nil, fmt.Errorf("解析分享列表失败: %v", err)
	}
	return &data, nil
}

// FileBatchShare 创建分享（使用 BasePanService）
func (x *XunleiPanService) FileBatchShare(ids []string, needPassword bool, expirationDays int) (*XLBatchShareResp, error) {
	apiURL := x.apiHost("") + "/drive/v1/share/batch"
	body := map[string]interface{}{
		"file_ids":        ids,
		"need_password":   needPassword,
		"expiration_days": expirationDays,
	}

	respData, err := x.HTTPPost(apiURL, body, nil)
	if err != nil {
		return nil, fmt.Errorf("创建分享失败: %v", err)
	}

	var data XLBatchShareResp
	if err := json.Unmarshal(respData, &data); err != nil {
		return nil, fmt.Errorf("解析分享响应失败: %v", err)
	}
	return &data, nil
}

// ShareBatchDelete 取消分享（使用 BasePanService）
func (x *XunleiPanService) ShareBatchDelete(ids []string) (*XLCommonResp, error) {
	apiURL := x.apiHost("") + "/drive/v1/share/batch/delete"
	body := map[string]interface{}{
		"share_ids": ids,
	}

	respData, err := x.HTTPPost(apiURL, body, nil)
	if err != nil {
		return nil, fmt.Errorf("删除分享失败: %v", err)
	}

	var data XLCommonResp
	if err := json.Unmarshal(respData, &data); err != nil {
		return nil, fmt.Errorf("解析删除响应失败: %v", err)
	}
	return &data, nil
}

// GetShareFolder 获取分享内容（使用 BasePanService）
func (x *XunleiPanService) GetShareFolder(shareID, passCodeToken, parentID string) (*XLShareFolderResp, error) {
	apiURL := x.apiHost("") + "/drive/v1/share/detail"
	body := map[string]interface{}{
		"share_id":        shareID,
		"pass_code_token": passCodeToken,
		"parent_id":       parentID,
		"limit":           100,
		"thumbnail_size":  "SIZE_LARGE",
		"order":           "6",
	}

	respData, err := x.HTTPPost(apiURL, body, nil)
	if err != nil {
		return nil, fmt.Errorf("获取分享文件夹失败: %v", err)
	}

	var data XLShareFolderResp
	if err := json.Unmarshal(respData, &data); err != nil {
		return nil, fmt.Errorf("解析分享文件夹失败: %v", err)
	}
	return &data, nil
}

// Restore 转存（使用 BasePanService）
func (x *XunleiPanService) Restore(shareID, passCodeToken string, fileIDs []string) (*XLRestoreResp, error) {
	apiURL := x.apiHost("") + "/drive/v1/share/restore"
	body := map[string]interface{}{
		"share_id":          shareID,
		"pass_code_token":   passCodeToken,
		"file_ids":          fileIDs,
		"folder_type":       "NORMAL",
		"specify_parent_id": true,
		"parent_id":         "",
	}

	respData, err := x.HTTPPost(apiURL, body, nil)
	if err != nil {
		return nil, fmt.Errorf("转存失败: %v", err)
	}

	var data XLRestoreResp
	if err := json.Unmarshal(respData, &data); err != nil {
		return nil, fmt.Errorf("解析转存响应失败: %v", err)
	}
	return &data, nil
}

// 结构体完全对齐 xunleix
type XLShareListResp struct {
	Data struct {
		List []struct {
			ShareID string `json:"share_id"`
			Title   string `json:"title"`
		} `json:"list"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type XLBatchShareResp struct {
	Data struct {
		ShareURL string `json:"share_url"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type XLCommonResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type XLShareFolderResp struct {
	Data struct {
		Files []struct {
			FileID string `json:"file_id"`
			Name   string `json:"name"`
		} `json:"files"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type XLRestoreResp struct {
	Data struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 新增辅助结构体
type XLShareInfo struct {
	ShareID string       `json:"share_id"`
	Title   string       `json:"title"`
	Files   []XLFileInfo `json:"files"`
}

type XLFileInfo struct {
	FileID string `json:"file_id"`
	Name   string `json:"name"`
}

type XLTaskResult struct {
	Status int    `json:"status"`
	TaskID string `json:"task_id"`
	Data   struct {
		ShareID string `json:"share_id"`
	} `json:"data"`
}
