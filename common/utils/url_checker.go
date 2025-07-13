package common

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// 定义网盘服务类型（使用字符串常量，避免与ServiceType冲突）
const (
	UCStr       = "uc"
	AliyunStr   = "aliyun"
	QuarkStr    = "quark"
	Pan115Str   = "115"
	Pan123Str   = "123pan"
	TianyiStr   = "tianyi"
	XunleiStr   = "xunlei"
	BaiduStr    = "baidu"
	NotFoundStr = "notfound"
)

// 检查结果结构
type CheckResult struct {
	URL    string
	Status bool
}

// 提取分享码和服务类型（重命名避免冲突）
func ExtractShareIdString(urlStr string) (string, string) {
	return extractShareID(urlStr)
}

// 提取分享码和服务类型
func extractShareID(urlStr string) (string, string) {
	netDiskPatterns := map[string]struct {
		Domains []string
		Pattern *regexp.Regexp
	}{
		UCStr: {
			Domains: []string{"drive.uc.cn"},
			Pattern: regexp.MustCompile(`https?://drive\.uc\.cn/s/([a-zA-Z0-9]+)`),
		},
		AliyunStr: {
			Domains: []string{"aliyundrive.com", "alipan.com"},
			Pattern: regexp.MustCompile(`https?://(?:www\.)?(?:aliyundrive|alipan)\.com/s/([a-zA-Z0-9]+)`),
		},
		QuarkStr: {
			Domains: []string{"pan.quark.cn"},
			Pattern: regexp.MustCompile(`https?://(?:www\.)?pan\.quark\.cn/s/([a-zA-Z0-9]+)`),
		},
		Pan115Str: {
			Domains: []string{"115.com", "115cdn.com", "anxia.com"},
			Pattern: regexp.MustCompile(`https?://(?:www\.)?(?:115|115cdn|anxia)\.com/s/([a-zA-Z0-9]+)`),
		},
		Pan123Str: {
			Domains: []string{"123684.com", "123685.com", "123912.com", "123pan.com", "123pan.cn", "123592.com"},
			Pattern: regexp.MustCompile(`https?://(?:www\.)?(?:123684|123685|123912|123pan|123pan\.cn|123592)\.com/s/([a-zA-Z0-9-]+)`),
		},
		TianyiStr: {
			Domains: []string{"cloud.189.cn"},
			Pattern: regexp.MustCompile(`https?://cloud\.189\.cn/(?:t/|web/share\?code=)([a-zA-Z0-9]+)`),
		},
		XunleiStr: {
			Domains: []string{"pan.xunlei.com"},
			Pattern: regexp.MustCompile(`https?://(?:www\.)?pan\.xunlei\.com/s/([a-zA-Z0-9-]+)`),
		},
		BaiduStr: {
			Domains: []string{"pan.baidu.com", "yun.baidu.com"},
			Pattern: regexp.MustCompile(`https?://(?:[a-z]+\.)?(?:pan|yun)\.baidu\.com/(?:s/|share/init\?surl=)([a-zA-Z0-9_-]+)(?:\?|$)`),
		},
	}

	for service, config := range netDiskPatterns {
		if containsDomain(urlStr, config.Domains) {
			match := config.Pattern.FindStringSubmatch(urlStr)
			if len(match) > 1 {
				return match[1], service
			}
		}
	}

	return "", NotFoundStr
}

// 检查域名是否包含在列表中
func containsDomain(urlStr string, domains []string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	host := u.Hostname()
	for _, domain := range domains {
		if strings.Contains(host, domain) {
			return true
		}
	}

	return false
}

// 创建HTTP客户端
func createHTTPClient() *resty.Client {
	return resty.New().
		SetTimeout(10 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)
}

// 检查UC网盘链接
func checkUC(shareID string) (bool, error) {
	client := createHTTPClient()
	url := fmt.Sprintf("https://drive.uc.cn/s/%s", shareID)

	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.101 Mobile Safari/537.36").
		SetHeader("Host", "drive.uc.cn").
		SetHeader("Referer", url).
		SetHeader("Origin", "https://drive.uc.cn").
		Get(url)

	if err != nil {
		return false, err
	}

	if resp.StatusCode() != 200 {
		return false, nil
	}

	bodyStr := resp.String()

	// 检查错误关键词
	errorKeywords := []string{"失效", "不存在", "违规", "删除", "已过期", "被取消"}
	for _, keyword := range errorKeywords {
		if strings.Contains(bodyStr, keyword) {
			return false, nil
		}
	}

	// 检查是否需要访问码
	if strings.Contains(bodyStr, "class=\"main-body\"") && strings.Contains(bodyStr, "class=\"input-wrap\"") {
		fmt.Println("发现访问码输入框，判断为有效（需密码）")
		return true, nil
	}

	// 检查是否包含文件列表或分享内容
	if strings.Contains(bodyStr, "文件") || strings.Contains(bodyStr, "分享") || strings.Contains(bodyStr, "class=\"file-list\"") {
		return true, nil
	}

	return false, nil
}

// 检查阿里云盘链接
func checkAliyun(shareID string) (bool, error) {
	client := createHTTPClient()
	apiURL := "https://api.aliyundrive.com/adrive/v3/share_link/get_share_by_anonymous"
	data := map[string]string{"share_id": shareID}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(apiURL)

	if err != nil {
		return false, err
	}

	var responseJSON map[string]interface{}
	err = json.Unmarshal(resp.Body(), &responseJSON)
	if err != nil {
		return false, err
	}

	if hasPwd, ok := responseJSON["has_pwd"].(bool); ok && hasPwd {
		return true, nil
	}

	if code, ok := responseJSON["code"].(string); ok && code == "NotFound.ShareLink" {
		return false, nil
	}

	if _, ok := responseJSON["file_infos"]; !ok {
		return false, nil
	}

	return true, nil
}

// 检查115网盘链接
func check115(shareID string) (bool, error) {
	client := createHTTPClient()
	apiURL := "https://webapi.115.com/share/snap"
	params := map[string]string{
		"share_code":   shareID,
		"receive_code": "",
	}

	resp, err := client.R().
		SetQueryParams(params).
		Get(apiURL)

	if err != nil {
		return false, err
	}

	var responseJSON map[string]interface{}
	err = json.Unmarshal(resp.Body(), &responseJSON)
	if err != nil {
		return false, err
	}

	if state, ok := responseJSON["state"].(bool); ok && state {
		return true, nil
	}

	if errorMsg, ok := responseJSON["error"].(string); ok && strings.Contains(errorMsg, "请输入访问码") {
		return true, nil
	}

	return false, nil
}

// 检查夸克网盘链接
func checkQuark(shareID string) (bool, error) {
	client := createHTTPClient()
	apiURL := "https://drive.quark.cn/1/clouddrive/share/sharepage/token"
	data := map[string]string{"pwd_id": shareID, "passcode": ""}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(apiURL)

	if err != nil {
		return false, err
	}

	var responseJSON map[string]interface{}
	err = json.Unmarshal(resp.Body(), &responseJSON)
	if err != nil {
		return false, err
	}

	if message, ok := responseJSON["message"].(string); ok && message == "ok" {
		data, ok := responseJSON["data"].(map[string]interface{})
		if !ok {
			return false, nil
		}

		stoken, ok := data["stoken"].(string)
		if !ok || stoken == "" {
			return false, nil
		}

		detailURL := fmt.Sprintf("https://drive-h.quark.cn/1/clouddrive/share/sharepage/detail?pwd_id=%s&stoken=%s&_fetch_share=1", shareID, url.QueryEscape(stoken))
		detailResp, err := client.R().Get(detailURL)
		if err != nil {
			return false, err
		}

		var detailResponseJSON map[string]interface{}
		err = json.Unmarshal(detailResp.Body(), &detailResponseJSON)
		if err != nil {
			return false, err
		}

		if status, ok := detailResponseJSON["status"].(float64); ok && status == 400 {
			return true, nil
		}

		if data, ok := detailResponseJSON["data"].(map[string]interface{}); ok {
			if share, ok := data["share"].(map[string]interface{}); ok {
				if status, ok := share["status"].(float64); ok && status == 1 {
					return true, nil
				}
			}
		}

		return false, nil
	} else if message, ok := responseJSON["message"].(string); ok && message == "需要提取码" {
		return true, nil
	}

	return false, nil
}

// 检查123网盘链接
func check123pan(shareID string) (bool, error) {
	client := createHTTPClient()
	apiURL := fmt.Sprintf("https://www.123pan.com/api/share/info?shareKey=%s", shareID)

	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0").
		Get(apiURL)

	if err != nil {
		return false, err
	}

	bodyStr := resp.String()

	if bodyStr == "" || strings.Contains(bodyStr, "分享页面不存在") {
		return false, nil
	}

	var responseJSON map[string]interface{}
	err = json.Unmarshal(resp.Body(), &responseJSON)
	if err != nil {
		return false, err
	}

	if code, ok := responseJSON["code"].(float64); ok && code != 0 {
		return false, nil
	}

	if data, ok := responseJSON["data"].(map[string]interface{}); ok {
		if hasPwd, ok := data["HasPwd"].(bool); ok && hasPwd {
			return true, nil
		}
	}

	return true, nil
}

// 检查天翼网盘链接
func checkTianyi(shareID string) (bool, error) {
	client := createHTTPClient()
	apiURL := "https://api.cloud.189.cn/open/share/getShareInfoByCodeV2.action"
	data := map[string]string{"shareCode": shareID}

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(data).
		Post(apiURL)

	if err != nil {
		return false, err
	}

	bodyStr := resp.String()

	errorKeywords := []string{"ShareInfoNotFound", "ShareNotFound", "FileNotFound", "ShareExpiredError", "ShareAuditNotPass"}
	for _, keyword := range errorKeywords {
		if strings.Contains(bodyStr, keyword) {
			return false, nil
		}
	}

	if strings.Contains(bodyStr, "needAccessCode") {
		return true, nil
	}

	return true, nil
}

// 检查迅雷网盘链接
func checkXunlei(shareID string) (bool, error) {
	client := createHTTPClient()
	tokenURL := "https://xluser-ssl.xunlei.com/v1/shield/captcha/init"
	data := map[string]interface{}{
		"client_id": "Xqp0kJBXWhwaTpB6",
		"device_id": "925b7631473a13716b791d7f28289cad",
		"action":    "get:/drive/v1/share",
		"meta": map[string]interface{}{
			"package_name":   "pan.xunlei.com",
			"client_version": "1.45.0",
			"captcha_sign":   "1.fe2108ad808a74c9ac0243309242726c",
			"timestamp":      "1645241033384",
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(tokenURL)

	if err != nil {
		return false, err
	}

	var tokenResponseJSON map[string]interface{}
	err = json.Unmarshal(resp.Body(), &tokenResponseJSON)
	if err != nil {
		return false, err
	}

	token, ok := tokenResponseJSON["captcha_token"].(string)
	if !ok || token == "" {
		return false, nil
	}

	apiURL := fmt.Sprintf("https://api-pan.xunlei.com/drive/v1/share?share_id=%s", shareID)
	resp, err = client.R().
		SetHeader("x-captcha-token", token).
		SetHeader("x-client-id", "Xqp0kJBXWhwaTpB6").
		SetHeader("x-device-id", "925b7631473a13716b791d7f28289cad").
		Get(apiURL)

	if err != nil {
		return false, err
	}

	bodyStr := resp.String()

	errorKeywords := []string{"NOT_FOUND", "SENSITIVE_RESOURCE", "EXPIRED"}
	for _, keyword := range errorKeywords {
		if strings.Contains(bodyStr, keyword) {
			return false, nil
		}
	}

	if strings.Contains(bodyStr, "PASS_CODE_EMPTY") {
		return true, nil
	}

	return true, nil
}

// 检查百度网盘链接
func checkBaidu(shareID string) (bool, error) {
	client := createHTTPClient()
	url := fmt.Sprintf("https://pan.baidu.com/s/%s", shareID)

	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36").
		Get(url)

	if err != nil {
		return false, err
	}

	bodyStr := resp.String()

	errorKeywords := []string{"分享的文件已经被取消", "分享已过期", "你访问的页面不存在", "你所访问的页面"}
	for _, keyword := range errorKeywords {
		if strings.Contains(bodyStr, keyword) {
			return false, nil
		}
	}

	if strings.Contains(bodyStr, "请输入提取码") || strings.Contains(bodyStr, "提取文件") {
		return true, nil
	}

	if strings.Contains(bodyStr, "过期时间") || strings.Contains(bodyStr, "文件列表") {
		return true, nil
	}

	return false, nil
}

// 检查URL有效性
func CheckURL(urlStr string) (CheckResult, error) {
	shareID, service := extractShareID(urlStr)
	if shareID == "" || service == NotFoundStr {
		fmt.Printf("无法识别的链接或网盘服务: %s\n", urlStr)
		return CheckResult{URL: urlStr, Status: false}, nil
	}

	checkFunctions := map[string]func(string) (bool, error){
		UCStr:     checkUC,
		AliyunStr: checkAliyun,
		QuarkStr:  checkQuark,
		Pan115Str: check115,
		Pan123Str: check123pan,
		TianyiStr: checkTianyi,
		XunleiStr: checkXunlei,
		BaiduStr:  checkBaidu,
	}

	if fn, ok := checkFunctions[service]; ok {
		result, err := fn(shareID)
		if err != nil {
			return CheckResult{URL: urlStr, Status: false}, err
		}
		return CheckResult{URL: urlStr, Status: result}, nil
	}

	fmt.Printf("未找到服务 %s 的检测函数\n", service)
	return CheckResult{URL: urlStr, Status: false}, nil
}

// 主函数
func Test() {
	urls := []string{
		// UC网盘
		"https://drive.uc.cn/s/e1ebe95d144c4?public=1", // UC网盘有效
		"https://drive.uc.cn/s/m7and23e132a1?public=1", // UC网盘无效
		// 阿里云
		"https://www.aliyundrive.com/s/hz1HHxhahsE", // aliyundrive 公开分享
		"https://www.alipan.com/s/QbaHJ71QjV1",      // alipan 公开分享
		"https://www.alipan.com/s/GMrv1QCZhNB",      // 带提取码
		"https://www.aliyundrive.com/s/p51zbVtgmy",  // 链接错误 NotFound.ShareLink
		"https://www.aliyundrive.com/s/hZnj4qLMMd9", // 空文件
		// 115
		"https://115cdn.com/s/swh88n13z72?password=r9b2#",
		"https://anxia.com/s/swhm75q3z5o?password=ayss",
		"https://115.com/s/swhsaua36a1?password=oc92", // 带访问码
		"https://115.com/s/sw313r03zx1",               // 分享的文件涉嫌违规，链接已失效
		// 夸克
		"https://pan.quark.cn/s/9803af406f13", // 公开分享
		"https://pan.quark.cn/s/f161a5364657", // 提取码
		"https://pan.quark.cn/s/9803af406f15", // 分享不存在
		"https://pan.quark.cn/s/b999385c0936", // 违规
		"https://pan.quark.cn/s/c66f71b6f7d5", // 取消分享
		// 123
		"https://www.123pan.com/s/i4uaTd-WHn0", // 公开分享
		"https://www.123912.com/s/U8f2Td-ZeOX",
		"https://www.123684.com/s/u9izjv-k3uWv",
		"https://www.123pan.com/s/A6cA-AKH11", // 外链不存在
		// 天翼
		"https://cloud.189.cn/t/viy2quQzMBne",              // 公开分享
		"https://cloud.189.cn/web/share?code=UfUjiiFRbymq", // 带密码分享长链接
		"https://cloud.189.cn/t/vENFvuVNbyqa",              // 外链不存在
		"https://cloud.189.cn/t/notexist",                  // 分享不存在
		// 百度
		"https://pan.baidu.com/s/1rIcc6X7D3rVzNSqivsRejw?pwd=0w0j", // 带提取码分享
		"https://pan.baidu.com/s/1TMhfQ5yNnlPPSGbw4RQ-LA?pwd=6j77", // 带提取码分享
		"https://pan.baidu.com/s/1J_CUxLKqC0h3Ypg4sQV0_g",          // 无法识别
		"https://pan.baidu.com/s/1HlvGfj8qVUBym24X2I9ukA",          // 分享被和谐
		"https://pan.baidu.com/s/1cgsY10lkrPGZ-zt8oVdR_w",          // 分享已过期
		"https://pan.baidu.com/s/1R_itrvmA0ZyMMaHybg7G2Q",          // 分享已删除
		"https://pan.baidu.com/s/1hqge8hI",                         // 分享链接错误
		"https://pan.baidu.com/s/1notexist",                        // 分享不存在
	}

	var wg sync.WaitGroup
	results := make(chan CheckResult, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result, err := CheckURL(u)
			if err != nil {
				fmt.Printf("检查 %s 时出错: %v\n", u, err)
				result = CheckResult{URL: u, Status: false}
			}
			results <- result
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("\n检测结果：")
	for result := range results {
		fmt.Printf("%s - %s\n", result.URL, map[bool]string{true: "有效", false: "无效"}[result.Status])
	}
}
