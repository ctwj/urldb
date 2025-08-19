// 1. 修正接口 Host，增加配置项
// 2. POST/GET 区分（xunleix 的 /drive/v1/share/list 是 GET，不是 POST）
// 3. 参数传递方式严格区分 query/body
// 4. header 应支持 Authorization（Bearer ...）、x-device-id、x-client-id、x-captcha-token 等
// 5. 结构体返回字段需和 xunleix 100%一致（如 data 字段是 map 还是 list），注意 code 字段为 int 还是 string
// 6. 错误处理，返回体未必有 code/msg，需先判断 HTTP 状态码再判断 body
// 7. 建议增加日志和更清晰的错误提示

package pan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type XunleiPanService struct {
	*BasePanService
	configMutex sync.RWMutex
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

// 工具：自动补全必要 header
func (x *XunleiPanService) setCommonHeader(req *http.Request) {
	for k, v := range x.headers {
		req.Header.Set(k, v)
	}
}

func NewXunleiPanService(config *PanConfig) *XunleiPanService {
	xunleiOnce.Do(func() {
		xunleiInstance = &XunleiPanService{
			BasePanService: NewBasePanService(config),
		}
		xunleiInstance.SetHeaders(map[string]string{
			"Content-Type": "application/json",
			// "access-control-allow-credentials": "true",
			// "access-control-allow-methods":     "GET, POST, PUT, DELETE, OPTIONS",
			// "access-control-allow-origin":      "https://pan.xunlei.com",
			// "access-control-max-age":           "86400",
			// "access-control-expose-headers": "Content-Type, Content-Length, Content-Encoding, X-Request-Id, X-Response-Time",
			// "access-control-allow-headers":  "Authorization, Content-Type,Accept, X-Project-Id, X-Device-Id, X-Request-Id, X-Captcha-Token, X-Client-Id, x-sdk-version, x-client-version, x-device-name, x-device-model, x-captcha-token, x-net-work-type, x-os-version, x-protocol-version, x-platform-version, x-provider-name, x-client-channel-id, x-appname, x-appid, x-device-sign, x-auto-login, x-peer-id, x-action",
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			// "X-Client-Id":                      "Xqp0kJBXWhwaTpB6",
			// "X-Device-Id":                      "c24ecadc44c643637d127fb847dbe36d",
			// "X-Device-Sign":                    "wdi10.c24ecadc44c643637d127fb847dbe36d74ee67f56a443148fc801eb27bb3e058",
			// "x-sdk-version":                    "5.2.4",
			// "x-protocol-version":               "301",
			"Authorization": config.Cookie,
		})
	})
	xunleiInstance.UpdateConfig(config)
	return xunleiInstance
}

// GetXunleiInstance 获取迅雷网盘服务单例实例
func GetXunleiInstance() *XunleiPanService {
	return NewXunleiPanService(nil)
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

// Transfer 转存分享链接 - 实现 PanService 接口
func (x *XunleiPanService) Transfer(shareID string) (*TransferResult, error) {
	// 读取配置（线程安全）
	x.configMutex.RLock()
	config := x.config
	x.configMutex.RUnlock()

	log.Printf("开始处理迅雷分享: %s", shareID)

	// 检查是否为检验模式
	if config.IsType == 1 {
		// 检验模式：直接获取分享信息
		shareInfo, err := x.getShareInfo(shareID)
		if err != nil {
			return ErrorResult(fmt.Sprintf("获取分享信息失败: %v", err)), nil
		}

		return SuccessResult("检验成功", map[string]interface{}{
			"title":    shareInfo.Title,
			"shareUrl": config.URL,
		}), nil
	}

	// 转存模式：实现完整的转存流程
	// 1. 获取分享详情
	shareDetail, err := x.GetShareFolder(shareID, "", "")
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取分享详情失败: %v", err)), nil
	}

	// 2. 提取文件ID列表
	fileIDs := make([]string, 0)
	for _, file := range shareDetail.Data.Files {
		fileIDs = append(fileIDs, file.FileID)
	}

	if len(fileIDs) == 0 {
		return ErrorResult("分享中没有可转存的文件"), nil
	}

	// 3. 转存文件
	restoreResult, err := x.Restore(shareID, "", fileIDs)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}

	// 4. 等待转存完成
	taskID := restoreResult.Data.TaskID
	_, err = x.waitForTask(taskID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("等待转存完成失败: %v", err)), nil
	}

	// 5. 创建新的分享
	shareResult, err := x.FileBatchShare(fileIDs, false, 0) // 永久分享
	if err != nil {
		return ErrorResult(fmt.Sprintf("创建分享失败: %v", err)), nil
	}

	// 6. 返回结果
	return SuccessResult("转存成功", map[string]interface{}{
		"shareUrl": shareResult.Data.ShareURL,
		"title":    fmt.Sprintf("迅雷分享_%s", shareID),
		"fid":      strings.Join(fileIDs, ","),
	}), nil
}

// waitForTask 等待任务完成
func (x *XunleiPanService) waitForTask(taskID string) (*XLTaskResult, error) {
	maxRetries := 50
	retryDelay := 2 * time.Second

	for retryIndex := 0; retryIndex < maxRetries; retryIndex++ {
		result, err := x.getTaskStatus(taskID, retryIndex)
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

// getTaskStatus 获取任务状态
func (x *XunleiPanService) getTaskStatus(taskID string, retryIndex int) (*XLTaskResult, error) {
	apiURL := x.apiHost("") + "/drive/v1/task"
	params := url.Values{}
	params.Set("task_id", taskID)
	params.Set("retry_index", fmt.Sprintf("%d", retryIndex))
	apiURL = apiURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	x.setCommonHeader(req)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(result))
	}
	var data XLTaskResult
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, err
	}
	return &data, nil
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

// GetFiles 获取文件列表 - 实现 PanService 接口
func (x *XunleiPanService) GetFiles(pdirFid string) (*TransferResult, error) {
	log.Printf("开始获取迅雷网盘文件列表，目录ID: %s", pdirFid)

	// 使用现有的 GetShareList 方法获取文件列表
	shareList, err := x.GetShareList("")
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取文件列表失败: %v", err)), nil
	}

	// 转换为通用格式
	fileList := make([]interface{}, 0)
	for _, share := range shareList.Data.List {
		fileList = append(fileList, map[string]interface{}{
			"share_id": share.ShareID,
			"title":    share.Title,
		})
	}

	return SuccessResult("获取成功", fileList), nil
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

// GetUserInfo 获取用户信息 - 实现 PanService 接口
func (x *XunleiPanService) GetUserInfo(cookie string) (*UserInfo, error) {
	log.Printf("开始获取迅雷网盘用户信息")

	// 临时设置cookie
	x.SetHeader("Authorization", cookie)

	// 获取用户信息
	apiURL := x.apiHost("user") + "/v1/user/me"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	x.setCommonHeader(req)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(result))
	}

	var response struct {
		Username string `json:"name"`
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	aboutURL := x.apiHost("") + "/drive/v1/about"
	req, err = http.NewRequest("GET", aboutURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	defer resp.Body.Close()
	result, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(result))
	}

	return &UserInfo{
		Username:    response.Username,
		ServiceType: "xunlei",
	}, nil
}

// GetShareList 严格对齐 GET + query（xunleix实现）
func (x *XunleiPanService) GetShareList(pageToken string) (*XLShareListResp, error) {
	api := x.apiHost("") + "/drive/v1/share/list"
	params := url.Values{}
	params.Set("limit", "100")
	params.Set("thumbnail_size", "SIZE_SMALL")
	if pageToken != "" {
		params.Set("page_token", pageToken)
	}
	apiURL := api + "?" + params.Encode()

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	x.setCommonHeader(req)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(result))
	}
	var data XLShareListResp
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// FileBatchShare 创建分享（POST, body）
func (x *XunleiPanService) FileBatchShare(ids []string, needPassword bool, expirationDays int) (*XLBatchShareResp, error) {
	apiURL := x.apiHost("") + "/drive/v1/share/batch"
	body := map[string]interface{}{
		"file_ids":        ids,
		"need_password":   needPassword,
		"expiration_days": expirationDays,
	}
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	x.setCommonHeader(req)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(result))
	}
	var data XLBatchShareResp
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// ShareBatchDelete 取消分享（POST, body）
func (x *XunleiPanService) ShareBatchDelete(ids []string) (*XLCommonResp, error) {
	apiURL := x.apiHost("") + "/drive/v1/share/batch/delete"
	body := map[string]interface{}{
		"share_ids": ids,
	}
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	x.setCommonHeader(req)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(result))
	}
	var data XLCommonResp
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// GetShareFolder 获取分享内容（POST, body）
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
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	x.setCommonHeader(req)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(result))
	}
	var data XLShareFolderResp
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// Restore 转存（POST, body）
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
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	x.setCommonHeader(req)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(result))
	}
	var data XLRestoreResp
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, err
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
