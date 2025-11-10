package pan

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// 新增常量定义
const (
	XLUSER_CLIENT_ID = "XW5SkOhLDjnOZP7J" // 登录
	PAN_CLIENT_ID    = "Xqp0kJBXWhwaTpB6" // 获取文件列表
	CLIENT_SECRET    = "Og9Vr1L8Ee6bh0olFxFDRg"
	CLIENT_VERSION   = "1.92.9" // 更新为与xunlei_3项目相同的版本
	PACKAG_ENAME     = "pan.xunlei.com"
)

var SALTS = []string{
	"QG3/GhopO+5+T",
	"1Sv94+ANND3lDmmw",
	"q2eTxRva8b3B5d",
	"m2",
	"VIc5CZRBMU71ENfbOh0+RgWIuzLy",
	"66M8Wpw6nkBEekOtL6e",
	"N0rucK7S8W/vrRkfPto5urIJJS8dVY0S",
	"oLAR7pdUVUAp9xcuHWzrU057aUhdCJrt",
	"6lxcykBSsfI//GR9",
	"r50cz+1I4gbU/fk8",
	"tdwzrTc4SNFC4marNGTgf05flC85A",
	"qvNVUDFjfsOMqvdi2gB8gCvtaJAIqxXs",
}

// captchaSign 生成验证码签名 - 完全复制自xunlei_3项目
func (x *XunleiPanService) captchaSign(clientId string, deviceID string, timestamp string) string {
	sign := clientId + CLIENT_VERSION + PACKAG_ENAME + deviceID + timestamp
	log.Printf("urldb 签名基础字符串: %s", sign)
	for _, salt := range SALTS { // salt =
		hash := md5.Sum([]byte(sign + salt))
		sign = hex.EncodeToString(hash[:])
	}
	log.Printf("urldb 最终签名: 1.%s", sign)
	return fmt.Sprintf("1.%s", sign)
}

// getTimestamp 获取当前时间戳
func (x *XunleiPanService) getTimestamp() int64 {
	return time.Now().UnixMilli()
}

// LoginWithCredentials 使用账号密码登录
func (x *XunleiPanService) LoginWithCredentials(username, password string) (XunleiTokenData, error) {
	loginURL := "https://xluser-ssl.xunlei.com/v1/auth/signin"

	// 初始化验证码 - 完全模仿xunlei_3的CaptchaInit方法
	captchaURL := "https://xluser-ssl.xunlei.com/v1/shield/captcha/init"

	// 构造meta参数（完全模仿xunlei_3，只包含phone_number）
	meta := map[string]interface{}{
		"phone_number": "+86" + username,
	}

	// 构造验证码请求（完全模仿xunlei_3）
	captchaBody := map[string]interface{}{
		"client_id": XLUSER_CLIENT_ID,
		"action":    "POST:/v1/auth/signin",
		"device_id": x.deviceId,
		"meta":      meta,
	}

	log.Printf("发送验证码初始化请求: %+v", captchaBody)
	resp, err := x.sendCaptchaRequest(captchaURL, captchaBody)
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("获取验证码失败: %v", err)
	}

	if resp["captcha_token"] == nil {
		return XunleiTokenData{}, fmt.Errorf("获取验证码失败: 响应中没有captcha_token")
	}

	captchaToken, ok := resp["captcha_token"].(string)
	if !ok {
		return XunleiTokenData{}, fmt.Errorf("获取验证码失败: captcha_token格式错误")
	}
	log.Printf("成功获取captcha_token: %s", captchaToken)

	// 构造登录请求数据
	loginData := map[string]interface{}{
		"client_id":     XLUSER_CLIENT_ID,
		"client_secret": CLIENT_SECRET,
		"password":      password,
		"username":      "+86 " + username,
		"captcha_token": captchaToken,
	}

	// 发送登录请求
	userInfo, err := x.sendCaptchaRequest(loginURL, loginData)
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("登录请求失败: %v", err)
	}

	// 提取token信息
	accessToken, ok := userInfo["access_token"].(string)
	if !ok {
		return XunleiTokenData{}, fmt.Errorf("登录响应中没有access_token")
	}

	refreshToken, ok := userInfo["refresh_token"].(string)
	if !ok {
		return XunleiTokenData{}, fmt.Errorf("登录响应中没有refresh_token")
	}

	sub, ok := userInfo["sub"].(string)
	if !ok {
		sub = ""
	}

	// 计算过期时间
	expiresIn := int64(3600) // 默认1小时
	if exp, ok := userInfo["expires_in"].(float64); ok {
		expiresIn = int64(exp)
	}
	expiresAt := time.Now().Unix() + expiresIn - 60 // 减去60秒缓冲

	log.Printf("登录成功，获取到token")
	return XunleiTokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    expiresAt,
		Sub:          sub,
		TokenType:    "Bearer",
		UserId:       sub,
	}, nil
}

// sendCaptchaRequest 发送验证码请求 - 完全复制xunlei_3的sendRequest实现
func (x *XunleiPanService) sendCaptchaRequest(url string, data map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Printf("发送验证码请求URL: %s", url)
	log.Printf("发送验证码请求数据: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 完全复制xunlei_3的请求头设置
	reqHeaders := x.getHeadersForRequest(nil)
	// 添加特定的headers
	reqHeaders["Content-Type"] = "application/x-www-form-urlencoded"
	reqHeaders["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}

	// 根据URL确定使用哪个client_id
	if strings.Contains(url, "shield/captcha/init") {
		// 对于验证码初始化，如果数据中指定了client_id，则使用该client_id
		if clientID, ok := data["client_id"].(string); ok {
			req.Header.Set("X-Client-Id", clientID)
		} else {
			// 默认使用PAN_CLIENT_ID用于API相关的验证码
			req.Header.Set("X-Client-Id", PAN_CLIENT_ID)
		}
	} else if strings.Contains(url, "auth/") {
		// 对于认证相关的请求，使用登录相关的client_id
		req.Header.Set("X-Client-Id", XLUSER_CLIENT_ID)
	} else {
		// 对于一般的API请求，使用PAN_CLIENT_ID
		req.Header.Set("X-Client-Id", PAN_CLIENT_ID)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("验证码响应状态码: %d", resp.StatusCode)
	log.Printf("验证码响应内容: %s", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v, raw: %s", err, string(body))
	}

	log.Printf("解析后的响应: %+v", result)
	return result, nil
}

// getHeadersForRequest 获取请求头
func (x *XunleiPanService) getHeadersForRequest(accessToken *string) map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	// 这里我们简化处理，因为验证码请求不需要这些
	// if x.CaptchaToken != nil {
	// 	headers["User-Agent"] = x.buildCustomUserAgent()
	// 	headers["X-Captcha-Token"] = *x.CaptchaToken
	// } else {
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	// }

	// if accessToken != nil {
	// 	headers["Authorization"] = fmt.Sprintf("Bearer %s", *accessToken)
	// }

	// if x.DeviceID != "" {
	// 	headers["X-Device-Id"] = x.DeviceID
	// }

	return headers
}