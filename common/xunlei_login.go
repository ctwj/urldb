package pan

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// 迅雷客户端常量与签名算法
// 参考 OpenList thunder driver（安卓客户端参数，已验证可用），
// 替换原失效的网页版参数（XW5SkOhLDjnOZP7J / Xqp0kJBXWhwaTpB6 / 旧 SALTS）。
// ============================================================================

const (
	// XLClientID 迅雷客户端 ID
	XLClientID = "Xp6vsxz_7IYVw2BB"
	// XLClientSecret 客户端密钥
	XLClientSecret = "Xp6vsy4tN9toTVdMSpomVdXpRmES"
	// XLClientVersion 客户端版本号
	XLClientVersion = "8.31.0.9726"
	// XLPackageName 安卓包名
	XLPackageName = "com.xunlei.downloadprovider"
	// XLUserAgent 安卓客户端 User-Agent
	XLUserAgent = "ANDROID-com.xunlei.downloadprovider/8.31.0.9726 netWorkType/5G appid/40 deviceName/Xiaomi_M2004j7ac deviceModel/M2004J7AC OSVersion/12 protocolVersion/301 platformVersion/10 sdkVersion/512000 Oauth2Client/0.9 (Linux 4_14_186-perf-gddfs8vbb238b) (JAVA 0)"
)

// Algorithms 验证码签名盐值数组（参考 OpenList thunder/driver.go:48-59）。
// 顺序敏感，禁止调整；用于 xlCaptchaSign 的多轮 MD5 迭代。
var Algorithms = []string{
	"9uJNVj/wLmdwKrJaVj/omlQ",
	"Oz64Lp0GigmChHMf/6TNfxx7O9PyopcczMsnf",
	"Eb+L7Ce+Ej48u",
	"jKY0",
	"ASr0zCl6v8W4aidjPK5KHd1Lq3t+vBFf41dqv5+fnOd",
	"wQlozdg6r1qxh0eRmt3QgNXOvSZO6q/GXK",
	"gmirk+ciAvIgA/cxUUCema47jr/YToixTT+Q6O",
	"5IiCoM9B1/788ntB",
	"P07JH0h6qoM6TSUAK2aL9T5s2QBVeY9JWvalf",
	"+oK0AN",
}

// xlCaptchaSign 生成迅雷验证码签名（参考 OpenList thunder/util.go:116-127 GetCaptchaSign）。
//
//	算法：base = ClientID + ClientVersion + PackageName + DeviceID + timestamp(毫秒)
//	  对 Algorithms 每项盐值依次 MD5Hex 迭代，sign = "1." + 最终哈希。
//
// 返回 (当前毫秒时间戳, "1."+sign)。timestamp 与 sign 必须实时生成，
// 严禁硬编码（FR-002 / SC-006）。
// xlCaptchaSignWithTimestamp 给定时间戳计算签名（纯函数，便于单元测试固化算法）。
//
//	s = ClientID + ClientVersion + PackageName + DeviceID + timestamp
//	对 Algorithms 每项盐值依次 MD5Hex 迭代，返回 "1." + 最终哈希。
func xlCaptchaSignWithTimestamp(deviceID, timestamp string) string {
	s := XLClientID + XLClientVersion + XLPackageName + deviceID + timestamp
	for _, a := range Algorithms {
		sum := md5.Sum([]byte(s + a))
		s = hex.EncodeToString(sum[:])
	}
	return "1." + s
}

// xlCaptchaSign 生成验证码签名，返回 (当前毫秒时间戳, "1."+签名)。
// timestamp 与 sign 必须实时生成，严禁硬编码（FR-002 / SC-006）。
func xlCaptchaSign(deviceID string) (timestamp, sign string) {
	timestamp = fmt.Sprint(time.Now().UnixMilli())
	sign = xlCaptchaSignWithTimestamp(deviceID, timestamp)
	return
}

// deriveDeviceID 按账号派生稳定的 32 位设备标识（参考 OpenList thunder/driver.go:60-65）。
// 各账号独立，避免全账号共用固定标识导致的干扰（R-05）。
func deriveDeviceID(username, password string) string {
	sum := md5.Sum([]byte(username + password))
	return hex.EncodeToString(sum[:])
}

var xlEmailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// ============================================================================
// 验证码令牌（captcha_token）—— 两阶段刷新
// 参考 OpenList thunder/util.go:92-166（RefreshCaptchaTokenInLogin / AtLogin）。
// ============================================================================

// captchaInit 调用 /v1/shield/captcha/init 换取验证码令牌（登录阶段与登录后阶段共用）。
//
//	action: 触发验证码的业务动作，如 "POST:/v1/auth/signin"、"GET:/drive/v1/share"。
//	meta  : 业务元数据；登录阶段仅含账号标识，登录后阶段含动态签名等。
//
// 成功返回新 captcha_token 并缓存；识别到风控（响应含非空 url）时返回明确错误（FR-009）。
func (x *XunleiPanService) captchaInit(action string, meta map[string]string) (string, error) {
	existingToken := ""
	if x.extra.Captcha != nil {
		existingToken = x.extra.Captcha.CaptchaToken
	}
	body := map[string]interface{}{
		"action":        action,
		"captcha_token": existingToken,
		"client_id":     XLClientID,
		"device_id":     x.deviceId,
		"meta":          meta,
		"redirect_uri":  "xlaccsdk01://xunlei.com/callback?state=harbor",
	}
	resp, err := x.sendXunleiAuthRequest("https://xluser-ssl.xunlei.com/v1/shield/captcha/init", body)
	if err != nil {
		return "", fmt.Errorf("captcha init 请求失败: %v", err)
	}

	// 风控：响应含非空 url 表示需要浏览器/短信验证（OpenList: need verify）
	if u, ok := resp["url"].(string); ok && u != "" {
		return "", fmt.Errorf("迅雷账号触发安全验证（need verify），请在迅雷官方端登录该账号一次后重试: %s", u)
	}

	token, _ := resp["captcha_token"].(string)
	if token == "" {
		return "", fmt.Errorf("captcha init 未返回 captcha_token: %v", resp)
	}

	// 缓存 captcha_token 与过期时间（提前 10 秒过期）
	expiresIn := int64(3600)
	if exp, ok := resp["expires_in"].(float64); ok {
		expiresIn = int64(exp)
	}
	if x.extra.Captcha == nil {
		x.extra.Captcha = &CaptchaData{}
	}
	x.extra.Captcha.CaptchaToken = token
	x.extra.Captcha.ExpiresAt = time.Now().Unix() + expiresIn - 10

	// 持久化（若有 repo；登录后阶段刷新时需要保存）
	if x.cksRepo != nil {
		x.persistExtra()
	}
	return token, nil
}

// refreshCaptchaTokenInLogin 登录阶段验证码刷新：meta 仅含账号标识，不带签名（无 user_id）。
func (x *XunleiPanService) refreshCaptchaTokenInLogin(action, username string) (string, error) {
	meta := map[string]string{}
	if xlEmailRegex.MatchString(username) {
		meta["email"] = username
	} else if len(username) >= 11 && len(username) <= 18 {
		meta["phone_number"] = username
	} else {
		meta["username"] = username
	}
	return x.captchaInit(action, meta)
}

// ============================================================================
// 登录（方案A：/v1/auth/signin）
// 若该接口被服务端限制，可在 LoginWithCredentials 内扩展方案B（CoreLogin + signin/token）。
// ============================================================================

// LoginWithCredentials 使用账号密码登录，返回令牌数据。
// 流程：按账号派生设备标识 → 登录阶段验证码刷新 → POST /v1/auth/signin → 解析 token。
func (x *XunleiPanService) LoginWithCredentials(username, password string) (XunleiTokenData, error) {
	// 登录前按账号派生稳定的设备标识（R-05），保证各账号独立
	x.deviceId = deriveDeviceID(username, password)

	captchaToken, err := x.refreshCaptchaTokenInLogin("POST:/v1/auth/signin", username)
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("获取验证码令牌失败: %v", err)
	}

	loginData := map[string]interface{}{
		"client_id":     XLClientID,
		"client_secret": XLClientSecret,
		"username":      "+86 " + username,
		"password":      password,
		"captcha_token": captchaToken,
	}
	resp, err := x.sendXunleiAuthRequest("https://xluser-ssl.xunlei.com/v1/auth/signin", loginData)
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("登录请求失败: %v", err)
	}

	accessToken, _ := resp["access_token"].(string)
	if accessToken == "" {
		return XunleiTokenData{}, fmt.Errorf("登录响应无 access_token: %v", resp)
	}
	refreshToken, _ := resp["refresh_token"].(string)
	sub, _ := resp["sub"].(string)
	expiresIn := int64(3600)
	if exp, ok := resp["expires_in"].(float64); ok {
		expiresIn = int64(exp)
	}

	log.Printf("迅雷账号 %s 登录成功", username)
	return XunleiTokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    time.Now().Unix() + expiresIn - 60,
		Sub:          sub,
		TokenType:    "Bearer",
		UserId:       sub,
	}, nil
}

// sendXunleiAuthRequest 发送迅雷认证类 POST 请求（captcha init / signin / token 刷新）。
// 这类请求需要特定 UA 与客户端头，且响应可能为非 200 的业务错误，故独立于 BasePanService 处理。
func (x *XunleiPanService) sendXunleiAuthRequest(reqURL string, data map[string]interface{}) (map[string]interface{}, error) {
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", XLUserAgent)
	req.Header.Set("X-Client-Id", XLClientID)
	req.Header.Set("X-Device-Id", x.deviceId)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("迅雷认证请求[%s] 响应[%d]: %s", reqURL, resp.StatusCode, string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("响应解析失败: %v, raw: %s", err, string(body))
	}
	return result, nil
}
