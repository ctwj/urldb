package pan

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// BaiduPanService 百度网盘服务 - 基于BaiduPanFilesTransfers项目改写
type BaiduPanService struct {
	*BasePanService
	bdstoken string
	entity   entity.Cks
	cksRepo  repo.CksRepository
	client   *http.Client
}

// NewBaiduPanService 创建百度网盘服务
func NewBaiduPanService(config *PanConfig) *BaiduPanService {
	service := &BaiduPanService{
		BasePanService: NewBasePanService(config),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// 设置与Python项目完全一致的请求头
	service.SetHeaders(map[string]string{
		"Host":               "pan.baidu.com",
		"Connection":         "keep-alive",
		"Upgrade-Insecure-Requests": "1",
		"Accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Sec-Fetch-Dest":     "document",
		"Sec-Fetch-Site":     "same-site",
		"Sec-Fetch-Mode":     "navigate",
		"Referer":            "https://pan.baidu.com",
		"Accept-Encoding":    "gzip, deflate, br",
		"Accept-Language":    "zh-CN,zh;q=0.9,en;q=0.8",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	})

	return service
}

// GetUserInfo 获取用户信息 - 基于Python项目的get_templatevariable方法
func (b *BaiduPanService) GetUserInfo(cookie *string) (*UserInfo, error) {
	// 设置Cookie
	b.SetHeader("Cookie", *cookie)

	// 首先调用gettemplatevariable获取基础token信息
	userInfoURL := "https://pan.baidu.com/api/gettemplatevariable"
	queryParams := map[string]string{
		"fields":     `["bdstoken","token","uk","isdocuser","servertime"]`,
		"clienttype": "0",
		"app_id":     "38824127",
		"web":        "1",
	}

	utils.Info("调用百度网盘用户信息API - URL: %s, Cookie长度: %d", userInfoURL, len(*cookie))

	// 构建查询参数
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		utils.Info("创建请求失败: %v", err)
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置查询参数
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// 设置请求头
	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := b.client.Do(req)
	if err != nil {
		utils.Info("百度网盘API调用失败: %v", err)
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	defer resp.Body.Close()

	// 处理gzip压缩
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("解压gzip失败: %v", err)
		}
		defer reader.(*gzip.Reader).Close()
	}

	// 读取响应
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	utils.Info("百度网盘API调用成功，响应长度: %d", len(body))

	// 解析响应
	var result struct {
		Errno int `json:"errno"`
		Result struct {
			Bdstoken   string `json:"bdstoken"`
			Token      string `json:"token"`
			UK         interface{} `json:"uk"`  // 可能是数字或字符串
			IsDocUser  int    `json:"isdocuser"`
			ServerTime int64  `json:"servertime"`
		} `json:"result"`
	}

	utils.Info("API响应内容: %s", string(body))

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	if result.Errno != 0 {
		return nil, fmt.Errorf("API返回错误: %d", result.Errno)
	}

	utils.Info("成功获取用户信息token: %+v", result)

	// 保存bdstoken
	b.bdstoken = result.Result.Bdstoken

	// 处理UK字段，可能是数字或字符串
	var uk string
	switch v := result.Result.UK.(type) {
	case string:
		uk = v
	case float64:
		uk = fmt.Sprintf("%.0f", v)
	case int:
		uk = fmt.Sprintf("%d", v)
	default:
		uk = ""
	}

	// 获取容量信息（必须在bdstoken设置之后调用）
	usedSpace, totalSpace, err := b.getQuotaInfo()
	if err != nil {
		utils.Info("获取容量信息失败: %v", err)
		// 如果获取容量失败，使用默认值0
		usedSpace = 0
		totalSpace = 0
	}

	// 返回用户信息
	return &UserInfo{
		Username:    uk,
		VIPStatus:   result.Result.IsDocUser == 1,
		UsedSpace:   usedSpace,
		TotalSpace:  totalSpace,
		ServiceType: "baidu",
		ExtraData:   fmt.Sprintf(`{"bdstoken":"%s","token":"%s","uk":"%s"}`, result.Result.Bdstoken, result.Result.Token, uk),
	}, nil
}

// SetCKSRepository 设置CKSRepository
func (b *BaiduPanService) SetCKSRepository(cksRepo repo.CksRepository, cks entity.Cks) {
	b.cksRepo = cksRepo
	b.entity = cks
}

// GetUserInfoByEntity 通过实体获取用户信息
func (b *BaiduPanService) GetUserInfoByEntity(cks entity.Cks) (*UserInfo, error) {
	cookie := cks.Ck
	return b.GetUserInfo(&cookie)
}

// Transfer 转存分享链接
func (b *BaiduPanService) Transfer(shareID string) (*TransferResult, error) {
	// 测试基本的链接解析功能
	testURL := "https://pan.baidu.com/s/1x0jzc7gzdMqtQB3GB2dcDA?pwd=4u76"

	// 解析链接和提取码
	linkURL, passCode := b.parseURLAndCode(testURL)
	if linkURL == "" {
		return &TransferResult{Success: false, Message: "链接解析失败"}, nil
	}

	// 验证链接
	result, err := b.verifyLink(linkURL, passCode)
	if err != nil {
		return &TransferResult{Success: false, Message: fmt.Sprintf("链接验证失败: %v", err)}, nil
	}

	return &TransferResult{
		Success: true,
		Message: "链接解析和验证成功",
		Data: map[string]interface{}{
			"linkURL":  linkURL,
			"passCode": passCode,
			"result":   result,
		},
	}, nil
}

// GetFiles 获取文件列表
func (b *BaiduPanService) GetFiles(pdirFid string) (*TransferResult, error) {
	// 暂时返回未实现
	return &TransferResult{Success: false, Message: "GetFiles方法尚未实现"}, nil
}

// DeleteFiles 删除文件
func (b *BaiduPanService) DeleteFiles(fileList []string) (*TransferResult, error) {
	// 暂时返回未实现
	return &TransferResult{Success: false, Message: "DeleteFiles方法尚未实现"}, nil
}

// GetServiceType 获取服务类型
func (b *BaiduPanService) GetServiceType() ServiceType {
	return BaiduPan
}

// UpdateConfig 更新配置
func (b *BaiduPanService) UpdateConfig(config *PanConfig) {
	// 暂时不实现
}

// parseURLAndCode 解析URL和提取码 - 基于Python的parse_url_and_code方法
func (b *BaiduPanService) parseURLAndCode(urlCode string) (string, string) {
	// 首先检查是否包含pwd参数（格式：https://pan.baidu.com/s/xxx?pwd=4u76）
	if strings.Contains(urlCode, "?pwd=") {
		// 解析URL参数
		parts := strings.Split(urlCode, "?pwd=")
		if len(parts) == 2 {
			url := parts[0]
			code := parts[1]
			// 处理提取码长度限制
			if len(code) > 4 {
				code = code[:4]
			}
			return url, code
		}
	}

	// 标准化链接格式
	normalized := b.normalizeLink(urlCode)

	// 以空格分割URL和提取码（格式：https://pan.baidu.com/s/xxx 4u76）
	parts := strings.SplitN(normalized, " ", 2)
	if len(parts) < 2 {
		// 如果没有空格分隔，可能是没有提取码的链接
		return normalized, ""
	}

	url := strings.TrimSpace(parts[0])
	code := strings.TrimSpace(parts[1])

	// 处理URL长度限制
	if len(url) > 47 {
		url = url[:47]
	}

	// 处理提取码长度限制
	if len(code) > 4 {
		code = code[len(code)-4:]
	}

	return url, code
}

// normalizeLink 标准化链接 - 基于Python的normalize_link方法
func (b *BaiduPanService) normalizeLink(link string) string {
	// 替换常见的短链接格式
	if strings.HasPrefix(link, "https://pan.baidu.com/s/") {
		// 提取/s/后面的部分
		parts := strings.Split(link, "/s/")
		if len(parts) >= 2 {
			// 移除查询参数部分，只保留surl
			surlPart := strings.Split(parts[1], "?")[0]
			return "https://pan.baidu.com/share/init?surl=" + surlPart
		}
	}

	if strings.HasPrefix(link, "https://pan.baidu.com/share/") {
		// 如果已经是share格式，直接返回
		return link
	}

	return link
}

// verifyLink 验证链接 - 基于Python的verify_link方法
func (b *BaiduPanService) verifyLink(url, password string) ([]string, error) {
	// 对于有提取码的链接先验证提取码
	if password != "" {
		bdclnd, err := b.verifyPassCode(url, password)
		if err != nil {
			return nil, err
		}

		// 如果验证失败，返回错误
		if bdclnd != "0" {
			return nil, fmt.Errorf("提取码验证失败")
		}
	}

	// 获取转存参数
	params, err := b.getTransferParams(url)
	if err != nil {
		return nil, err
	}

	return params, nil
}

// verifyPassCode 验证提取码 - 基于Python的verify_pass_code方法
func (b *BaiduPanService) verifyPassCode(shareURL, passCode string) (string, error) {
	// 先标准化链接
	normalizedURL := b.normalizeLink(shareURL)
	surl := b.extractSurl(normalizedURL)
	verifyURL := "https://pan.baidu.com/share/verify"

	utils.Info("验证提取码 - 原链接: %s, 标准化链接: %s, 提取码: %s, surl: %s, bdstoken: %s", shareURL, normalizedURL, passCode, surl, b.bdstoken)

	// 构建表单数据
	formData := url.Values{}
	formData.Set("surl", surl)
	formData.Set("bdstoken", b.bdstoken)
	formData.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))
	formData.Set("channel", "chunlei")
	formData.Set("web", "1")
	formData.Set("clienttype", "0")
	formData.Set("pwd", passCode)
	formData.Set("vcode", "")
	formData.Set("vcode_str", "")

	req, err := http.NewRequest("POST", verifyURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "0", err
	}

	// 设置表单请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 设置其他请求头
	for key, value := range b.headers {
		if key != "Content-Type" {
			req.Header.Set(key, value)
		}
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return "0", err
	}
	defer resp.Body.Close()

	// 处理gzip压缩
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return "0", fmt.Errorf("解压gzip失败: %v", err)
		}
		defer reader.(*gzip.Reader).Close()
	}

	// 读取响应
	respData, err := io.ReadAll(reader)
	if err != nil {
		return "0", err
	}

	if resp.StatusCode != http.StatusOK {
		return "0", fmt.Errorf("HTTP请求失败: %d, %s", resp.StatusCode, string(respData))
	}

	utils.Info("验证提取码API响应: %s", string(respData))

	var result struct {
		Errno int `json:"errno"`
		Randsk string `json:"randsk"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return "0", err
	}

	if result.Errno != 0 {
		return fmt.Sprintf("%d", result.Errno), nil
	}

	return result.Randsk, nil
}

// extractSurl 从分享链接中提取surl - 基于Python的extract_surl方法
func (b *BaiduPanService) extractSurl(shareURL string) string {
	// 使用正则表达式提取surl
	re := regexp.MustCompile(`surl=([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(shareURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// getTransferParams 获取转存参数 - 基于Python的get_transfer_params方法
func (b *BaiduPanService) getTransferParams(shareURL string) ([]string, error) {
	req, err := http.NewRequest("GET", shareURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 处理gzip压缩
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("解压gzip失败: %v", err)
		}
		defer reader.(*gzip.Reader).Close()
	}

	// 读取响应
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	utils.Info("获取分享页面响应成功，长度: %d", len(body))
	utils.Info("响应内容前100字符: %s", string(body[:100]))

	// 解析响应内容，提取必要参数 - 基于Python的parse_response方法
	return b.parseResponse(string(body))
}

// parseResponse 解析响应内容 - 基于Python的parse_response方法
func (b *BaiduPanService) parseResponse(html string) ([]string, error) {
	// 使用正则表达式提取必要参数
	shareidRegex := regexp.MustCompile(`"shareid":\s*"(\d+)"`)
	ukRegex := regexp.MustCompile(`"uk":\s*"(\d+)"`)
	fsIdListRegex := regexp.MustCompile(`"fs_idlist":\s*\[([^\]]+)\]`)

	shareidMatch := shareidRegex.FindStringSubmatch(html)
	ukMatch := ukRegex.FindStringSubmatch(html)
	fsIdListMatch := fsIdListRegex.FindStringSubmatch(html)

	if len(shareidMatch) < 2 || len(ukMatch) < 2 {
		return nil, fmt.Errorf("无法解析分享参数")
	}

	shareid := shareidMatch[1]
	uk := ukMatch[1]
	var fsIdList string = ""

	if len(fsIdListMatch) > 1 {
		fsIdList = fsIdListMatch[1]
	}

	return []string{shareid, uk, fsIdList}, nil
}

// getQuotaInfo 获取百度网盘容量信息
func (b *BaiduPanService) getQuotaInfo() (int64, int64, error) {
	// 百度网盘容量查询API
	quotaURL := "https://pan.baidu.com/api/quota"
	queryParams := map[string]string{
		"clienttype": "0",
		"app_id":     "38824127",
		"web":        "1",
		"bdstoken":   b.bdstoken,
	}

	utils.Info("调用百度网盘容量查询API - URL: %s, bdstoken: %s", quotaURL, b.bdstoken)

	// 发送请求
	data, err := b.HTTPGet(quotaURL, queryParams)
	if err != nil {
		return 0, 0, fmt.Errorf("调用容量API失败: %v", err)
	}

	utils.Info("百度网盘容量API响应: %s", string(data))

	// 解析响应
	var result struct {
		Errno int `json:"errno"`
		Total int64 `json:"total"`     // 总容量（字节）
		Used  int64 `json:"used"`      // 已使用容量（字节）
		Free  int64 `json:"free"`      // 免费容量（字节）
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return 0, 0, fmt.Errorf("解析容量信息失败: %v", err)
	}

	if result.Errno != 0 {
		return 0, 0, fmt.Errorf("容量API返回错误: %d", result.Errno)
	}

	utils.Info("成功获取容量信息 - 总容量: %d, 已使用: %d, 免费容量: %d",
		result.Total, result.Used, result.Free)

	return result.Used, result.Total, nil
}