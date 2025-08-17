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
	"net/http"
	"net/url"
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
func (x *XunleiPanService) apiHost() string {
	// if x.config != nil && x.config.ApiHost != "" {
	// 	return x.config.ApiHost
	// }
	// 推荐用官方: https://api-pan.xunlei.com
	return "https://api-pan.xunlei.com"
}

// 工具：自动补全必要 header
func (x *XunleiPanService) setCommonHeader(req *http.Request) {
	for k, v := range x.headers {
		req.Header.Set(k, v)
	}
	// 可扩展: Authorization, x-device-id, x-client-id, x-captcha-token
}

func NewXunleiPanService(config *PanConfig) *XunleiPanService {
	xunleiOnce.Do(func() {
		xunleiInstance = &XunleiPanService{
			BasePanService: NewBasePanService(config),
		}
		xunleiInstance.SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Cookie":       config.Cookie,
		})
	})
	xunleiInstance.UpdateConfig(config)
	return xunleiInstance
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

// GetShareList 严格对齐 GET + query（xunleix实现）
func (x *XunleiPanService) GetShareList(pageToken string) (*XLShareListResp, error) {
	api := x.apiHost() + "/drive/v1/share/list"
	params := url.Values{}
	params.Set("limit", "100")
	params.Set("thumbnail_size", "SIZE_SMALL")
	if pageToken != "" {
		params.Set("page_token", pageToken)
	}
	url := api + "?" + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
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
	url := x.apiHost() + "/drive/v1/share/batch"
	body := map[string]interface{}{
		"file_ids":        ids,
		"need_password":   needPassword,
		"expiration_days": expirationDays,
	}
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bs))
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
	url := x.apiHost() + "/drive/v1/share/batch/delete"
	body := map[string]interface{}{
		"share_ids": ids,
	}
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bs))
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
	url := x.apiHost() + "/drive/v1/share/detail"
	body := map[string]interface{}{
		"share_id":        shareID,
		"pass_code_token": passCodeToken,
		"parent_id":       parentID,
		"limit":           100,
		"thumbnail_size":  "SIZE_LARGE",
		"order":           "6",
	}
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bs))
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
	url := x.apiHost() + "/drive/v1/share/restore"
	body := map[string]interface{}{
		"share_id":          shareID,
		"pass_code_token":   passCodeToken,
		"file_ids":          fileIDs,
		"folder_type":       "NORMAL",
		"specify_parent_id": true,
		"parent_id":         "",
	}
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bs))
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
