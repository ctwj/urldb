package pan

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	Captcha     *CaptchaData              `json:"captcha,omitempty"`
	Token       *XunleiTokenData          `json:"token,omitempty"`
	Credentials *XunleiAccountCredentials `json:"credentials,omitempty"` // 账号密码信息
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
	xunleiInstance := &XunleiPanService{
		BasePanService: NewBasePanService(config),
		clientId:       XLClientID,
		// 占位默认设备标识；实际按账号派生（见 SetCKSRepository / LoginWithCredentials，R-05）
		deviceId: deriveDeviceID("urldb", "xunlei"),
		extra:    XunleiExtraData{},
	}
	xunleiInstance.SetHeaders(map[string]string{
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Cache-Control":   "no-cache",
		"Content-Type":    "application/json",
		"Origin":          "https://pan.xunlei.com",
		"Pragma":          "no-cache",
		"Referer":         "https://pan.xunlei.com/",
		"User-Agent":      XLUserAgent,
		"Authorization":   "",
		"x-captcha-token": "",
		"x-client-id":     xunleiInstance.clientId,
		"x-device-id":     xunleiInstance.deviceId,
	})

	xunleiInstance.UpdateConfig(config)
	return xunleiInstance
}

// SetCKSRepository 设置 CksRepository 和 entity
func (x *XunleiPanService) SetCKSRepository(cksRepo repo.CksRepository, entity entity.Cks) {
	x.cksRepo = cksRepo
	x.entity = entity
	var extra XunleiExtraData

	// 解析extra字段
	if x.entity.Extra != "" {
		if err := json.Unmarshal([]byte(x.entity.Extra), &extra); err != nil {
			log.Printf("解析 extra 数据失败: %v", err)
		}
	}

	// 从ck字段解析账号密码，并按账号派生稳定的设备标识（R-05）
	if credentials, err := ParseCredentialsFromCk(x.entity.Ck); err == nil {
		extra.Credentials = credentials
		if credentials.Username != "" {
			x.deviceId = deriveDeviceID(credentials.Username, credentials.Password)
			x.SetHeader("x-device-id", x.deviceId)
		}
	}

	x.extra = extra
}

// persistExtra 序列化 extra 并持久化到数据库；同时把 refresh_token 冗余写入 ck 字段以向后兼容。
func (x *XunleiPanService) persistExtra() {
	if x.cksRepo == nil {
		return
	}
	extraBytes, err := json.Marshal(x.extra)
	if err != nil {
		log.Printf("序列化 extra 失败: %v", err)
		return
	}
	x.entity.Extra = string(extraBytes)
	if x.extra.Token != nil && x.extra.Token.RefreshToken != "" {
		x.entity.Ck = x.extra.Token.RefreshToken
	}
	if err := x.cksRepo.UpdateWithAllFields(&x.entity); err != nil {
		log.Printf("保存 extra 到数据库失败: %v", err)
	}
}

// GetXunleiInstance 获取迅雷网盘服务单例实例
func GetXunleiInstance() *XunleiPanService {
	return NewXunleiPanService(nil)
}

func (x *XunleiPanService) GetAccessTokenByRefreshToken(refreshToken string) (XunleiTokenData, error) {
	body := map[string]interface{}{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     XLClientID,
		"client_secret": XLClientSecret,
	}
	resp, err := x.sendXunleiAuthRequest("https://xluser-ssl.xunlei.com/v1/auth/token", body)
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("刷新 access_token 请求失败: %v", err)
	}

	accessToken, _ := resp["access_token"].(string)
	if accessToken == "" {
		return XunleiTokenData{}, fmt.Errorf("刷新 access_token 失败（响应无 access_token）: %v", resp)
	}
	expiresIn := int64(3600)
	if exp, ok := resp["expires_in"].(float64); ok {
		expiresIn = int64(exp)
	}
	refresh := ""
	if rt, ok := resp["refresh_token"].(string); ok {
		refresh = rt
	}
	sub, _ := resp["sub"].(string)
	return XunleiTokenData{
		AccessToken:  accessToken,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		ExpiresAt:    time.Now().Unix() + expiresIn - 60,
		Sub:          sub,
		TokenType:    "Bearer",
		UserId:       sub,
	}, nil
}

// reloginWithCredentials 使用账号密码重新登录
func (x *XunleiPanService) reloginWithCredentials() (XunleiTokenData, error) {
	if x.extra.Credentials == nil {
		return XunleiTokenData{}, fmt.Errorf("无账号密码信息")
	}

	tokenData, err := x.LoginWithCredentials(x.extra.Credentials.Username, x.extra.Credentials.Password)
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("账号密码登录失败: %v", err)
	}

	log.Printf("账号 %s 重新登录成功", x.extra.Credentials.Username)
	return tokenData, nil
}

// getAccessToken 获取 Access Token（内部包含缓存判断、刷新、重新登录、保存）
func (x *XunleiPanService) getAccessToken() (string, error) {
	// 检查 Access Token 是否有效
	currentTime := time.Now().Unix()
	if x.extra.Token != nil && x.extra.Token.AccessToken != "" && x.extra.Token.ExpiresAt > currentTime {
		return x.extra.Token.AccessToken, nil
	}

	// 尝试使用refresh_token刷新
	var newData XunleiTokenData
	var err error

	if x.extra.Token != nil && x.extra.Token.RefreshToken != "" {
		newData, err = x.GetAccessTokenByRefreshToken(x.extra.Token.RefreshToken)
		if err != nil {
			log.Printf("refresh_token刷新失败: %v，尝试使用账号密码重新登录", err)

			// 如果refresh_token失效且有账号密码信息，尝试重新登录
			if x.extra.Credentials != nil && x.extra.Credentials.Username != "" && x.extra.Credentials.Password != "" {
				newData, err = x.reloginWithCredentials()
				if err != nil {
					return "", fmt.Errorf("重新登录失败: %v", err)
				}
			} else {
				return "", fmt.Errorf("refresh_token失效且无账号密码信息，无法重新登录: %v", err)
			}
		}
	} else {
		return "", fmt.Errorf("无有效的refresh_token")
	}

	// 更新token信息
	if x.extra.Token == nil {
		x.extra.Token = &XunleiTokenData{}
	}
	x.extra.Token.AccessToken = newData.AccessToken
	x.extra.Token.RefreshToken = newData.RefreshToken
	x.extra.Token.ExpiresAt = newData.ExpiresAt
	x.extra.Token.ExpiresIn = newData.ExpiresIn
	x.extra.Token.Sub = newData.Sub
	x.extra.Token.TokenType = newData.TokenType
	x.extra.Token.UserId = newData.UserId

	// 持久化（含 ck 字段 refresh_token 同步）
	x.persistExtra()

	return newData.AccessToken, nil
}

// getCaptchaToken 获取 captcha_token（登录后阶段）。
// 登录后阶段需要 user_id（取自令牌）+ 动态 captcha_sign（基于实时时间戳，R-04/R-06）。
func (x *XunleiPanService) getCaptchaToken() (string, error) {
	currentTime := time.Now().Unix()
	if x.extra.Captcha != nil && x.extra.Captcha.CaptchaToken != "" && x.extra.Captcha.ExpiresAt > currentTime {
		return x.extra.Captcha.CaptchaToken, nil
	}

	// user_id 取自令牌（登录后才有）
	userID := ""
	if x.extra.Token != nil {
		userID = x.extra.Token.UserId
		if userID == "" {
			userID = x.extra.Token.Sub
		}
	}

	ts, sign := xlCaptchaSign(x.deviceId)
	meta := map[string]string{
		"client_version": XLClientVersion,
		"package_name":   XLPackageName,
		"user_id":        userID,
		"timestamp":      ts,
		"captcha_sign":   sign,
	}
	return x.captchaInit("GET:/drive/v1/share", meta)
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

	// 风控：review_panel 表示需要短信验证（参考 OpenList thunder），给出明确提示（FR-009）
	if errMsg, _ := result["error"].(string); errMsg == "review_panel" {
		return nil, fmt.Errorf("迅雷账号触发安全验证（review_panel），请在迅雷官方端登录该账号一次后重试")
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

func extractCode(url string) string {
	// 查找 pwd= 的位置
	if pwdIndex := strings.Index(url, "pwd="); pwdIndex != -1 {
		code := url[pwdIndex+4:]
		// 移除 # 及后面的内容（如果存在）
		if hashIndex := strings.Index(code, "#"); hashIndex != -1 {
			code = code[:hashIndex]
		}
		return code
	}
	return ""
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

	// 转存模式：实现完整的转存流程
	thisCode := extractCode(config.URL)

	// 获取分享详情
	shareDetail, err := x.getShare(shareID, thisCode, accessToken, captchaToken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取分享详情失败: %v", err)), nil
	}
	if shareDetail["share_status"].(string) != "OK" {
		return ErrorResult(fmt.Sprintf("获取分享详情失败: %v", "分享状态异常")), nil
	}
	if shareDetail["file_num"].(string) == "0" {
		return ErrorResult(fmt.Sprintf("获取分享详情失败: %v", "文件列表为空")), nil
	}

	parent_id := "" // 默认存储路径

	// 检查是否为检验模式
	if config.IsType == 1 {
		// 检验模式：直接获取分享信息
		urls := map[string]interface{}{
			"title":     shareDetail["title"],
			"share_url": config.URL,
			"stoken":    "",
		}
		return SuccessResult("检验成功", urls), nil
	}

	// files := shareDetail["files"].([]interface{})
	// fileIDs := make([]string, 0)
	// for _, file := range files {
	// 	fileMap := file.(map[string]interface{})
	// 	if fid, ok := fileMap["id"].(string); ok {
	// 		fileIDs = append(fileIDs, fid)
	// 	}
	// }

	// 处理广告过滤（这里简化处理）
	// TODO: 添加广告文件过滤逻辑

	// 转存资源
	restoreResult, err := x.getRestore(shareID, shareDetail, accessToken, captchaToken, parent_id)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}

	// 获取转存任务信息
	taskID := restoreResult["restore_task_id"].(string)

	// 等待转存完成
	taskResp, err := x.waitForTask(taskID, accessToken, captchaToken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("等待转存完成失败: %v", err)), nil
	}

	// 获取任务结果以获取文件ID
	existingFileIds := make([]string, 0)
	if params, ok2 := taskResp["params"].(map[string]interface{}); ok2 {
		if traceIds, ok3 := params["trace_file_ids"].(string); ok3 {
			traceData := make(map[string]interface{})
			json.Unmarshal([]byte(traceIds), &traceData)
			for _, fid := range traceData {
				existingFileIds = append(existingFileIds, fid.(string))
			}
		}
	}

	// 创建分享链接
	expirationDays := "-1"
	if config.ExpiredType == 2 {
		expirationDays = "2"
	}

	// 根据share_id获取到分享链接
	shareResult, err := x.getSharePassword(existingFileIds, accessToken, captchaToken, expirationDays)
	if err != nil {
		return ErrorResult(fmt.Sprintf("创建分享链接失败: %v", err)), nil
	}

	var fid string
	if len(existingFileIds) > 1 {
		fid = strings.Join(existingFileIds, ",")
	} else {
		fid = existingFileIds[0]
	}

	result := map[string]interface{}{
		"title":    "",
		"shareUrl": shareResult["share_url"].(string) + "?pwd=" + shareResult["pass_code"].(string),
		"code":     shareResult["pass_code"].(string),
		"fid":      fid,
	}

	return SuccessResult("转存成功", result), nil
}

// waitForTask 等待任务完成 - 使用 HTTPGet 方法
func (x *XunleiPanService) waitForTask(taskID string, accessToken, captchaToken string) (map[string]interface{}, error) {
	maxRetries := 50
	retryDelay := 2 * time.Second

	for retryIndex := 0; retryIndex < maxRetries; retryIndex++ {
		result, err := x.getTaskStatus(taskID, retryIndex, accessToken, captchaToken)
		if err != nil {
			return nil, err
		}

		if int64(result["progress"].(float64)) == 100 { // 任务完成
			return result, nil
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("任务超时")
}

// getTaskStatus 获取任务状态 - 使用 HTTPGet 方法
func (x *XunleiPanService) getTaskStatus(taskID string, retryIndex int, accessToken, captchaToken string) (map[string]interface{}, error) {
	apiURL := x.apiHost("") + "/drive/v1/tasks/" + taskID
	queryParams := map[string]string{}

	// 设置 request 所需的 headers
	headers := map[string]string{
		"Authorization":   "Bearer " + accessToken,
		"x-captcha-token": captchaToken,
	}

	resp, err := x.requestXunleiApi(apiURL, "GET", nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	return resp, nil
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
func (x *XunleiPanService) GetUserInfo(cookie *string) (*UserInfo, error) {
	userInfo := &UserInfo{}
	accessToken, err := x.getAccessToken()
	if err != nil {
		return nil, err
	}

	captchaToken, err := x.getCaptchaToken()
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for k, v := range x.headers {
		headers[k] = v
	}
	headers["Authorization"] = "Bearer " + accessToken
	headers["x-captcha-token"] = captchaToken

	resp, err := x.requestXunleiApi("https://api-pan.xunlei.com/drive/v1/about", "GET", nil, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	if quota, ok := resp["quota"].(map[string]interface{}); ok {
		if limit, ok := quota["limit"].(string); ok {
			userInfo.TotalSpace, _ = strconv.ParseInt(limit, 10, 64)
		}
		if usage, ok := quota["usage"].(string); ok {
			userInfo.UsedSpace, _ = strconv.ParseInt(usage, 10, 64)
		}
	}

	// 获取用户信息
	respData, err := x.requestXunleiApi(x.apiHost("user")+"/v1/user/me", "GET", nil, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	if name, ok := respData["name"].(string); ok {
		userInfo.Username = name
	}
	isVip := false
	if vipInfo, ok := respData["vip_info"].([]interface{}); ok && len(vipInfo) > 0 {
		if v0, ok := vipInfo[0].(map[string]interface{}); ok {
			if isVipStr, ok := v0["is_vip"].(string); ok {
				isVip = isVipStr != "0"
			}
		}
	}
	userInfo.ServiceType = x.GetServiceType().String()
	userInfo.VIPStatus = isVip
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

// (sendCaptchaRequestForGeneralAPI 已移除：验证码初始化统一由 captchaInit 处理)

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