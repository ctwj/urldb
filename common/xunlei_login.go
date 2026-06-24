package pan

import (
	"crypto/md5"
	"crypto/sha1"
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
// 参考 OpenList thunder driver（安卓客户端参数，已验证可用）。
// ============================================================================

const (
	// XLClientID 迅雷客户端 ID（迅雷下载管家 com.xunlei.downloadprovider）
	// 注意：refresh_token 绑定 client_id，必须与抓取 token 的 APP 一致（手机迅雷 = 下载管家）
	XLClientID = "Xp6vsxz_7IYVw2BB"
	// XLClientSecret 客户端密钥
	XLClientSecret = "Xp6vsy4tN9toTVdMSpomVdXpRmES"
	// XLClientVersion 客户端版本号
	XLClientVersion = "8.31.0.9726"
	// XLPackageName 安卓包名
	XLPackageName = "com.xunlei.downloadprovider"
	// XLUserAgent 安卓客户端 User-Agent（网盘 API）
	XLUserAgent = "ANDROID-com.xunlei.downloadprovider/8.31.0.9726 netWorkType/5G appid/40 deviceName/Xiaomi_M2004j7ac deviceModel/M2004J7AC OSVersion/12 protocolVersion/301 platformVersion/10 sdkVersion/512000 Oauth2Client/0.9 (Linux 4_14_186-perf-gddfs8vbb238b) (JAVA 0)"
	// XLCoreLoginUA CoreLogin(v3) 专用 UA
	XLCoreLoginUA = "android-ok-http-client/xl-acc-sdk/version-5.0.12.512000"
	// XLAppID/XLAppKey 设备签名所需（参考 OpenList thunder/util.go:40-44）
	XLAppID  = "40"
	XLAppKey = "34a062aaa22f906fca4fefe9fb3a3021"
	// XLSignProvider signin/token 的 provider
	XLSignProvider = "access_end_point_token"
)

// Algorithms 验证码签名盐值数组（参考 OpenList thunder/driver.go:48-59，迅雷下载管家版）。
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

// generateDeviceSign 生成设备签名（参考 OpenList thunder/util.go:265-284 generateDeviceSign）。
// 算法：SHA1(DeviceID+PackageName+AppID+AppKey) → MD5 → "div101."+DeviceID+MD5。
// 用于 CoreLogin(v3)，是设备信任与绕过 review 的关键。
func generateDeviceSign(deviceID, packageName string) string {
	base := deviceID + packageName + XLAppID + XLAppKey
	sha1Sum := sha1.Sum([]byte(base))
	sha1Str := hex.EncodeToString(sha1Sum[:])
	md5Sum := md5.Sum([]byte(sha1Str))
	md5Str := hex.EncodeToString(md5Sum[:])
	return "div101." + deviceID + md5Str
}

var xlEmailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// ============================================================================
// 验证码令牌（captcha_token）—— 两阶段刷新
// 参考 OpenList thunder/util.go:92-166（RefreshCaptchaTokenInLogin / AtLogin）。
// ============================================================================

// captchaInit 调用 /v1/shield/captcha/init 换取验证码令牌（登录阶段与登录后阶段共用）。
//
//	action: 触发验证码的业务动作，如 "POST:/v1/auth/signin/token"、"GET:/drive/v1/share"。
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
	resp, err := x.sendXunleiAuthRequest("https://xluser-ssl.xunlei.com/v1/shield/captcha/init", body, nil)
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
// 登录（方案B：CoreLogin + signin/token）
// 方案A（/v1/auth/signin 纯密码）会被风控 review 拦截（captcha_invalid/4002），
// 故采用 OpenList 验证可用的方案B：CoreLogin(v3, devicesign) → signin/token。
// ============================================================================

// coreLogin 调用 v3 登录接口换取 sessionID（参考 OpenList thunder/driver.go:686-729）。
// devicesign 是设备信任关键；creditkey 为空时为首次登录。
// 触发 review 时返回明确错误（含 reviewurl），提示用户在官方端验证（FR-009）。
func (x *XunleiPanService) coreLogin(username, password, creditkey string) (sessionID string, err error) {
	body := map[string]interface{}{
		"protocolVersion": "301",
		"sequenceNo":      "1000012",
		"platformVersion": "10",
		"isCompressed":    "0",
		"appid":           XLAppID,
		"clientVersion":   XLClientVersion,
		"peerID":          "00000000000000000000000000000000",
		"appName":         "ANDROID-" + XLPackageName,
		"sdkVersion":      "512000",
		"devicesign":      generateDeviceSign(x.deviceId, XLPackageName),
		"netWorkType":     "WIFI",
		"providerName":    "NONE",
		"deviceModel":     "M2004J7AC",
		"deviceName":      "Xiaomi_M2004j7ac",
		"OSVersion":       "12",
		"creditkey":       creditkey,
		"hl":              "zh-CN",
		"userName":        username,
		"passWord":        password,
		"verifyKey":       "",
		"verifyCode":      "",
		"isMd5Pwd":        "0",
	}
	resp, err := x.sendXunleiAuthRequest("https://xluser-ssl.xunlei.com/xluser.core.login/v3/login", body, map[string]string{"User-Agent": XLCoreLoginUA})
	if err != nil {
		return "", fmt.Errorf("CoreLogin 请求失败: %v", err)
	}

	// review：返回 reviewurl 表示需要安全验证（短信验证）。
	// 处理方式（参考 OpenList）：拼上 devicesign 作为 deviceid 参数（验证页据此关联设备），
	// 并返回 creditkey 供用户验证后重新登录时填入。
	if reviewURL, _ := resp["reviewurl"].(string); reviewURL != "" {
		ck, _ := resp["creditkey"].(string)
		devSign := generateDeviceSign(x.deviceId, XLPackageName)
		return "", fmt.Errorf("迅雷账号触发安全验证（review，需短信验证）。请按以下步骤处理：\n"+
			"1. 在浏览器（建议手机浏览器）打开链接完成短信验证:\n   %s&deviceid=%s\n"+
			"2. 验证成功后，重新添加账号，并在账号 JSON 中加入 creditkey 字段:\n"+
			`   {"username":"%s","password":"<密码>","creditkey":"%s"}`+"\n"+
			"（首次新设备验证一次后，该 deviceID 即被信任；creditkey 仅验证流程使用）",
			reviewURL, devSign, username, ck)
	}
	// 业务错误
	if errMsg, _ := resp["error"].(string); errMsg != "" && errMsg != "success" {
		return "", fmt.Errorf("CoreLogin 失败: %v", resp)
	}

	sessionID, _ = resp["sessionID"].(string)
	if sessionID == "" {
		return "", fmt.Errorf("CoreLogin 未返回 sessionID: %v", resp)
	}
	return sessionID, nil
}

// LoginWithCredentials 使用账号密码登录（方案B），返回令牌数据。
// creditkey 用于通过 review（首次新设备短信验证后获得，参考 OpenList）。
// 流程：按账号派生设备标识 → CoreLogin 拿 sessionID → 登录阶段验证码刷新 → signin/token 换 access_token。
func (x *XunleiPanService) LoginWithCredentials(username, password, creditkey string) (XunleiTokenData, error) {
	// 按账号派生稳定的设备标识（R-05）
	x.deviceId = deriveDeviceID(username, password)

	// 1. CoreLogin(v3) 拿 sessionID（devicesign 提供设备信任；creditkey 绕过 review）
	sessionID, err := x.coreLogin(username, password, creditkey)
	if err != nil {
		return XunleiTokenData{}, err
	}

	// 2. 登录阶段验证码刷新（signin/token 需要，通过 X-Captcha-Token 头传递）
	captchaToken, err := x.refreshCaptchaTokenInLogin("POST:/v1/auth/signin/token", username)
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("获取验证码令牌失败: %v", err)
	}

	// 3. POST /v1/auth/signin/token：用 sessionID 换 access_token
	body := map[string]interface{}{
		"client_id":     XLClientID,
		"client_secret": XLClientSecret,
		"provider":      XLSignProvider,
		"signin_token":  sessionID,
	}
	resp, err := x.sendXunleiAuthRequest("https://xluser-ssl.xunlei.com/v1/auth/signin/token", body, map[string]string{"X-Captcha-Token": captchaToken})
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("signin/token 请求失败: %v", err)
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

	log.Printf("迅雷账号 %s 登录成功（方案B CoreLogin）", username)
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

// LoginByRefreshToken 用 refresh_token 直接登录（跳过 CoreLogin/review）。
// 适用于从手机迅雷 APP 抓取到 refresh_token 的场景：避免新设备触发短信验证（FR-009 替代路径）。
func (x *XunleiPanService) LoginByRefreshToken(refreshToken string) (XunleiTokenData, error) {
	if refreshToken == "" {
		return XunleiTokenData{}, fmt.Errorf("refresh_token 为空")
	}
	// refresh_token 登录无 username，按 token 派生稳定 deviceId（用于后续 captcha 签名）
	x.deviceId = deriveDeviceID(refreshToken, "xlrefresh")
	x.SetHeader("x-device-id", x.deviceId)
	return x.GetAccessTokenByRefreshToken(refreshToken)
}

// sendXunleiAuthRequest 发送迅雷认证类 POST 请求（captcha init / CoreLogin / signin/token / token 刷新）。
// extraHeaders 用于按接口附加特殊头（如 CoreLogin 的 UA、signin/token 的 X-Captcha-Token）。
func (x *XunleiPanService) sendXunleiAuthRequest(reqURL string, data map[string]interface{}, extraHeaders map[string]string) (map[string]interface{}, error) {
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", XLUserAgent)
	req.Header.Set("X-Client-Id", XLClientID)
	req.Header.Set("X-Device-Id", x.deviceId)
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

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
