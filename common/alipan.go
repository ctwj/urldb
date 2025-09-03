package pan

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
)

// AlipanService 阿里云盘服务
type AlipanService struct {
	*BasePanService
	accessToken string
	configMutex sync.RWMutex // 保护配置的读写锁
}

// 单例相关变量
var (
	alipanInstance *AlipanService
	alipanOnce     sync.Once
)

// NewAlipanService 创建阿里云盘服务（单例模式）
func NewAlipanService(config *PanConfig) *AlipanService {
	alipanOnce.Do(func() {
		alipanInstance = &AlipanService{
			BasePanService: NewBasePanService(config),
		}

		// 设置阿里云盘的默认请求头
		alipanInstance.SetHeaders(map[string]string{
			"Accept":             "application/json, text/plain, */*",
			"Accept-Language":    "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
			"Content-Type":       "application/json",
			"Origin":             "https://www.alipan.com",
			"Priority":           "u=1, i",
			"Referer":            "https://www.alipan.com/",
			"Sec-Ch-Ua":          `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": `"Windows"`,
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "same-site",
			"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0",
			"X-Canary":           "client=web,app=share,version=v2.3.1",
		})
	})

	// 更新配置
	alipanInstance.UpdateConfig(config)

	return alipanInstance
}

// GetAlipanInstance 获取阿里云盘服务单例实例
func GetAlipanInstance() *AlipanService {
	return NewAlipanService(nil)
}

// UpdateConfig 更新配置（线程安全）
func (a *AlipanService) UpdateConfig(config *PanConfig) {
	if config == nil {
		return
	}

	a.configMutex.Lock()
	defer a.configMutex.Unlock()

	a.config = config
}

// GetServiceType 获取服务类型
func (a *AlipanService) GetServiceType() ServiceType {
	return Alipan
}

// Transfer 转存分享链接
func (a *AlipanService) Transfer(shareID string) (*TransferResult, error) {
	// 读取配置（线程安全）
	a.configMutex.RLock()
	config := a.config
	a.configMutex.RUnlock()

	fmt.Printf("开始处理阿里云盘分享: %s", shareID)

	// 获取access token
	accessToken, err := a.manageAccessToken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取access_token失败: %v", err)), nil
	}

	// 设置Authorization头
	a.SetHeader("Authorization", "Bearer "+accessToken)

	// 获取分享信息
	shareInfo, err := a.getAlipan1(shareID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取分享信息失败: %v", err)), nil
	}

	if config.IsType == 1 {
		// 直接返回资源信息
		return SuccessResult("检验成功", map[string]interface{}{
			"title":    shareInfo.ShareName,
			"shareUrl": config.URL,
		}), nil
	}

	// 获取share token
	shareTokenResult, err := a.getAlipan2(shareID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取share_token失败: %v", err)), nil
	}

	// 确定存储路径
	toPdirFid := "root" // 默认存储路径，可以从配置中读取
	if config.ExpiredType == 2 {
		toPdirFid = "temp" // 临时资源路径，可以从配置中读取
	}

	// 构建批量复制请求
	batchRequests := make([]map[string]interface{}, 0)
	for i, fileInfo := range shareInfo.FileInfos {
		request := map[string]interface{}{
			"body": map[string]interface{}{
				"auto_rename":       true,
				"file_id":           fileInfo.FileID,
				"share_id":          shareID,
				"to_drive_id":       "2008425230",
				"to_parent_file_id": toPdirFid,
			},
			"headers": map[string]string{
				"Content-Type": "application/json",
			},
			"id":     fmt.Sprintf("%d", i),
			"method": "POST",
			"url":    "/file/copy",
		}
		batchRequests = append(batchRequests, request)
	}

	batchData := map[string]interface{}{
		"requests": batchRequests,
		"resource": "file",
	}

	// 执行批量复制
	copyResult, err := a.getAlipan3(batchData, shareTokenResult.ShareToken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("批量复制失败: %v", err)), nil
	}

	// 提取复制后的文件ID
	fileIDList := make([]string, 0)
	for _, response := range copyResult.Responses {
		if response.Body.Code != "" {
			return ErrorResult(fmt.Sprintf("复制失败: %s", response.Body.Message)), nil
		}
		fileIDList = append(fileIDList, response.Body.FileID)
	}

	// 创建分享
	shareData := map[string]interface{}{
		"drive_id":     "2008425230",
		"expiration":   "",
		"share_pwd":    "",
		"file_id_list": fileIDList,
	}

	shareResult, err := a.getAlipan4(shareData)
	if err != nil {
		return ErrorResult(fmt.Sprintf("创建分享失败: %v", err)), nil
	}

	return SuccessResult("转存成功", map[string]interface{}{
		"shareUrl": shareResult.ShareURL,
		"title":    shareResult.ShareTitle,
		"fid":      shareResult.FileIDList,
	}), nil
}

// GetFiles 获取文件列表
func (a *AlipanService) GetFiles(pdirFid string) (*TransferResult, error) {
	// 获取access token
	accessToken, err := a.manageAccessToken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取access_token失败: %v", err)), nil
	}

	// 设置Authorization头
	a.SetHeader("Authorization", "Bearer "+accessToken)

	if pdirFid == "" {
		pdirFid = "root"
	}

	data := map[string]interface{}{
		"all":             false,
		"drive_id":        "2008425230",
		"fields":          "*",
		"limit":           100,
		"order_by":        "updated_at",
		"order_direction": "DESC",
		"parent_file_id":  pdirFid,
		"url_expire_sec":  14400,
	}

	respData, err := a.HTTPPost("https://api.aliyundrive.com/adrive/v3/file/list", data, nil)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取文件列表失败: %v", err)), nil
	}

	var response struct {
		Message string        `json:"message"`
		Items   []interface{} `json:"items"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return ErrorResult("解析响应失败"), nil
	}

	if response.Message != "" {
		return ErrorResult(response.Message), nil
	}

	return SuccessResult("获取成功", response.Items), nil
}

// DeleteFiles 删除文件
func (a *AlipanService) DeleteFiles(fileList []string) (*TransferResult, error) {
	// 获取access token
	accessToken, err := a.manageAccessToken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取access_token失败: %v", err)), nil
	}

	// 设置Authorization头
	a.SetHeader("Authorization", "Bearer "+accessToken)

	data := map[string]interface{}{
		"drive_id":     "2008425230",
		"file_id_list": fileList,
	}

	_, err = a.HTTPPost("https://api.aliyundrive.com/adrive/v3/file/delete", data, nil)
	if err != nil {
		return ErrorResult(fmt.Sprintf("删除文件失败: %v", err)), nil
	}

	return SuccessResult("删除成功", nil), nil
}

// GetUserInfo 获取用户信息
func (a *AlipanService) GetUserInfo(cookie *string) (*UserInfo, error) {
	// 设置Cookie
	a.SetHeader("Cookie", *cookie)

	// 获取access token
	accessToken, err := a.manageAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取access_token失败: %v", err)
	}

	// 设置Authorization头
	a.SetHeader("Authorization", "Bearer "+accessToken)

	// 调用阿里云盘用户信息API
	userInfoURL := "https://api.alipan.com/v2/user/get"
	resp, err := a.HTTPGet(userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 解析响应
	var result struct {
		Code string `json:"code"`
		Data struct {
			NickName  string `json:"nick_name"`
			Avatar    string `json:"avatar"`
			DriveInfo struct {
				TotalSize string `json:"total_size"`
				UsedSize  string `json:"used_size"`
			} `json:"drive_info"`
			VipInfo struct {
				VipStatus string `json:"vip_status"`
			} `json:"vip_info"`
		} `json:"data"`
	}

	if err := a.ParseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	if result.Code != "" {
		return nil, fmt.Errorf("API返回错误: %s", result.Code)
	}

	// 转换VIP状态
	vipStatus := result.Data.VipInfo.VipStatus == "vip"

	// 转换容量字符串为字节数
	totalSizeStr := result.Data.DriveInfo.TotalSize
	usedSizeStr := result.Data.DriveInfo.UsedSize

	// 解析容量字符串
	totalSizeBytes := ParseCapacityString(totalSizeStr)
	usedSizeBytes := ParseCapacityString(usedSizeStr)

	return &UserInfo{
		Username:    result.Data.NickName,
		VIPStatus:   vipStatus,
		UsedSpace:   usedSizeBytes,
		TotalSpace:  totalSizeBytes,
		ServiceType: "alipan",
	}, nil
}

// getAlipan1 通过分享id获取file_id
func (a *AlipanService) getAlipan1(shareID string) (*AlipanShareInfo, error) {
	data := map[string]interface{}{
		"share_id": shareID,
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 临时设置headers
	originalHeaders := a.headers
	a.SetHeaders(headers)
	defer func() { a.headers = originalHeaders }()

	respData, err := a.HTTPPost("https://api.aliyundrive.com/adrive/v3/share_link/get_share_by_anonymous", data, nil)
	if err != nil {
		return nil, err
	}

	var result AlipanShareInfo
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUserInfoByEntity 根据 entity.Cks 获取用户信息（待实现）
func (a *AlipanService) GetUserInfoByEntity(cks entity.Cks) (*UserInfo, error) {
	return nil, nil
}

// getAlipan2 通过分享id获取X-Share-Token
func (a *AlipanService) getAlipan2(shareID string) (*AlipanShareToken, error) {
	data := map[string]interface{}{
		"share_id": shareID,
	}

	respData, err := a.HTTPPost("https://api.aliyundrive.com/v2/share_link/get_share_token", data, nil)
	if err != nil {
		return nil, err
	}

	var result AlipanShareToken
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// getAlipan3 批量复制
func (a *AlipanService) getAlipan3(batchData map[string]interface{}, shareToken string) (*AlipanBatchResult, error) {
	// 设置X-Share-Token头
	a.SetHeader("X-Share-Token", shareToken)

	respData, err := a.HTTPPost("https://api.aliyundrive.com/adrive/v4/batch", batchData, nil)
	if err != nil {
		return nil, err
	}

	var result AlipanBatchResult
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// getAlipan4 创建分享
func (a *AlipanService) getAlipan4(shareData map[string]interface{}) (*AlipanShareResult, error) {
	respData, err := a.HTTPPost("https://api.aliyundrive.com/adrive/v3/share_link/create", shareData, nil)
	if err != nil {
		return nil, err
	}

	var result AlipanShareResult
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// manageAccessToken 管理access token
func (a *AlipanService) manageAccessToken() (string, error) {
	if a.accessToken != "" {
		return a.accessToken, nil
	}

	// 从文件读取token
	tokenFile := filepath.Join("config", "alipan_access_token.json")

	// 检查token文件是否存在
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		// 获取新的access token
		return a.getNewAccessToken()
	}

	// 读取token文件
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return a.getNewAccessToken()
	}

	var tokenInfo struct {
		AccessToken string    `json:"access_token"`
		ExpiresAt   time.Time `json:"expires_at"`
	}

	if err := json.Unmarshal(data, &tokenInfo); err != nil {
		return a.getNewAccessToken()
	}

	// 检查token是否过期
	if utils.GetCurrentTime().After(tokenInfo.ExpiresAt) {
		return a.getNewAccessToken()
	}

	a.accessToken = tokenInfo.AccessToken
	return a.accessToken, nil
}

// getNewAccessToken 获取新的access token
func (a *AlipanService) getNewAccessToken() (string, error) {
	// 这里需要实现获取新token的逻辑
	// 可能需要用户登录或者使用refresh token
	return "", fmt.Errorf("需要实现获取新access token的逻辑")
}

// 定义阿里云盘相关的结构体
type AlipanShareInfo struct {
	ShareName string `json:"share_name"`
	FileInfos []struct {
		FileID string `json:"file_id"`
	} `json:"file_infos"`
}

type AlipanShareToken struct {
	ShareToken string `json:"share_token"`
}

type AlipanBatchResult struct {
	Responses []struct {
		Body struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			FileID  string `json:"file_id"`
		} `json:"body"`
	} `json:"responses"`
}

type AlipanShareResult struct {
	ShareURL   string   `json:"share_url"`
	ShareTitle string   `json:"share_title"`
	FileIDList []string `json:"file_id_list"`
}
