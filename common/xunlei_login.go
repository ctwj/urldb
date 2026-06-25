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
// 迅雷客户端身份配置（Profile）
// refresh_token 绑定 client_id，下载管家（手机迅雷 APP）与迅雷浏览器两 APP 的
// client_id / 签名盐值 / 设备签名密钥 / 网盘 API host 均不同，按账号 client_type
// 选择对应 profile。
// ============================================================================

// xunleiProfile 迅雷客户端身份参数集
type xunleiProfile struct {
	Name          string   // "android" | "browser"
	ClientID      string
	ClientSecret  string
	ClientVersion string
	PackageName   string
	UserAgent     string   // 网盘 API User-Agent
	CoreLoginUA   string   // CoreLogin(v3) 专用 UA
	AppID         string   // 设备签名 APPID
	AppKey        string   // 设备签名 APPKey
	Algorithms    []string // 验证码签名盐值（顺序敏感）
	PanAPIHost    string   // 网盘 API host
}

// xlProfileAndroid 下载管家（com.xunlei.downloadprovider）—— 手机迅雷 APP
var xlProfileAndroid = xunleiProfile{
	Name:          "android",
	ClientID:      "Xp6vsxz_7IYVw2BB",
	ClientSecret:  "Xp6vsy4tN9toTVdMSpomVdXpRmES",
	ClientVersion: "8.31.0.9726",
	PackageName:   "com.xunlei.downloadprovider",
	UserAgent:     "ANDROID-com.xunlei.downloadprovider/8.31.0.9726 netWorkType/5G appid/40 deviceName/Xiaomi_M2004j7ac deviceModel/M2004J7AC OSVersion/12 protocolVersion/301 platformVersion/10 sdkVersion/512000 Oauth2Client/0.9 (Linux 4_14_186-perf-gddfs8vbb238b) (JAVA 0)",
	CoreLoginUA:   "android-ok-http-client/xl-acc-sdk/version-5.0.12.512000",
	AppID:         "40",
	AppKey:        "34a062aaa22f906fca4fefe9fb3a3021",
	Algorithms: []string{
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
	},
	PanAPIHost: "https://api-pan.xunlei.com",
}

// xlProfileBrowser 迅雷浏览器（com.xunlei.browser）
var xlProfileBrowser = xunleiProfile{
	Name:          "browser",
	ClientID:      "ZUBzD9J_XPXfn7f7",
	ClientSecret:  "yESVmHecEe6F0aou69vl-g",
	ClientVersion: "1.40.0.7208",
	PackageName:   "com.xunlei.browser",
	UserAgent:     "ANDROID-com.xunlei.browser/1.40.0.7208 networkType/WIFI appid/22062 deviceName/Xiaomi_M2004j7ac deviceModel/M2004J7AC OSVersion/13 protocolVersion/301 platformversion/10 sdkVersion/509300 Oauth2Client/0.9 (Linux 4_9_337-perf-sn-uotan-gd9d488809c3d) (JAVA 0) ",
	CoreLoginUA:   "android-ok-http-client/xl-acc-sdk/version-5.0.12.509300",
	AppID:         "22062",
	AppKey:        "a5d7416858147a4ab99573872ffccef8",
	Algorithms: []string{
		"Cw4kArmKJ/aOiFTxnQ0ES+D4mbbrIUsFn",
		"HIGg0Qfbpm5ThZ/RJfjoao4YwgT9/M",
		"u/PUD",
		"OlAm8tPkOF1qO5bXxRN2iFttuDldrg",
		"FFIiM6sFhWhU7tIMVUKOF7CUv/KzgwwV8FE",
		"yN",
		"4m5mglrIHksI6wYdq",
		"LXEfS7",
		"T+p+C+F2yjgsUtiXWU/cMNYEtJI4pq7GofW",
		"14BrGIEMXkbvFvZ49nDUfVCRcHYFOJ1BP1Y",
		"kWIH3Row",
		"RAmRTKNCjucPWC",
	},
	PanAPIHost: "https://x-api-pan.xunlei.com",
}

// xlProfileByType 按 client_type 选择 profile，默认 android（向后兼容）
func xlProfileByType(t string) xunleiProfile {
	if t == "browser" {
		return xlProfileBrowser
	}
	return xlProfileAndroid
}

// xlSignProvider signin/token 的 provider（两身份相同）
const xlSignProvider = "access_end_point_token"

// xlSignWithTimestamp 给定 profile 与时间戳计算验证码签名（纯函数，便于单测固化算法）。
//
//	s = ClientID + ClientVersion + PackageName + DeviceID + timestamp
//	对 Algorithms 每项盐值依次 MD5Hex 迭代，返回 "1." + 最终哈希。
func xlSignWithTimestamp(p xunleiProfile, deviceID, timestamp string) string {
	s := p.ClientID + p.ClientVersion + p.PackageName + deviceID + timestamp
	for _, a := range p.Algorithms {
		sum := md5.Sum([]byte(s + a))
		s = hex.EncodeToString(sum[:])
	}
	return "1." + s
}

// xlGenerateDeviceSign 生成设备签名（参考 OpenList generateDeviceSign）。
// 算法：SHA1(DeviceID+PackageName+AppID+AppKey) → MD5 → "div101."+DeviceID+MD5。
func xlGenerateDeviceSign(p xunleiProfile, deviceID string) string {
	base := deviceID + p.PackageName + p.AppID + p.AppKey
	sha1Sum := sha1.Sum([]byte(base))
	sha1Str := hex.EncodeToString(sha1Sum[:])
	md5Sum := md5.Sum([]byte(sha1Str))
	return "div101." + deviceID + hex.EncodeToString(md5Sum[:])
}

// deriveDeviceID 按账号派生稳定的 32 位设备标识（各账号独立，R-05）。
func deriveDeviceID(username, password string) string {
	sum := md5.Sum([]byte(username + password))
	return hex.EncodeToString(sum[:])
}

var xlEmailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// captchaSign 用当前 profile 生成验证码签名，返回 (当前毫秒时间戳, "1."+签名)。
// timestamp 与 sign 必须实时生成，严禁硬编码（FR-002 / SC-006）。
func (x *XunleiPanService) captchaSign() (timestamp, sign string) {
	timestamp = fmt.Sprint(time.Now().UnixMilli())
	sign = xlSignWithTimestamp(x.profile, x.deviceId, timestamp)
	return
}

// ============================================================================
// 验证码令牌（captcha_token）—— 两阶段刷新
// ============================================================================

// captchaInit 调用 /v1/shield/captcha/init 换取验证码令牌（两阶段共用）。
func (x *XunleiPanService) captchaInit(action string, meta map[string]string) (string, error) {
	existingToken := ""
	if x.extra.Captcha != nil {
		existingToken = x.extra.Captcha.CaptchaToken
	}
	body := map[string]interface{}{
		"action":        action,
		"captcha_token": existingToken,
		"client_id":     x.profile.ClientID,
		"device_id":     x.deviceId,
		"meta":          meta,
		"redirect_uri":  "xlaccsdk01://xunlei.com/callback?state=harbor",
	}
	resp, err := x.sendXunleiAuthRequest("https://xluser-ssl.xunlei.com/v1/shield/captcha/init", body, nil)
	if err != nil {
		return "", fmt.Errorf("captcha init 请求失败: %v", err)
	}
	// 风控：响应含非空 url 表示需要浏览器/短信验证
	if u, ok := resp["url"].(string); ok && u != "" {
		return "", fmt.Errorf("迅雷账号触发安全验证（need verify），请在迅雷官方端登录该账号一次后重试: %s", u)
	}
	token, _ := resp["captcha_token"].(string)
	if token == "" {
		return "", fmt.Errorf("captcha init 未返回 captcha_token: %v", resp)
	}
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

// refreshCaptchaTokenInLogin 登录阶段验证码刷新：meta 仅含账号标识，不带签名。
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
// 账号密码登录（CoreLogin + signin/token）
// 保留供未来迅雷放宽风控时复用；当前前端不再调用（被 review 拦截）。
// ============================================================================

// coreLogin 调用 v3 登录接口换取 sessionID。
func (x *XunleiPanService) coreLogin(username, password, creditkey string) (sessionID string, err error) {
	body := map[string]interface{}{
		"protocolVersion": "301",
		"sequenceNo":      "1000012",
		"platformVersion": "10",
		"isCompressed":    "0",
		"appid":           x.profile.AppID,
		"clientVersion":   x.profile.ClientVersion,
		"peerID":          "00000000000000000000000000000000",
		"appName":         "ANDROID-" + x.profile.PackageName,
		"sdkVersion":      "512000",
		"devicesign":      xlGenerateDeviceSign(x.profile, x.deviceId),
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
	resp, err := x.sendXunleiAuthRequest("https://xluser-ssl.xunlei.com/xluser.core.login/v3/login", body, map[string]string{"User-Agent": x.profile.CoreLoginUA})
	if err != nil {
		return "", fmt.Errorf("CoreLogin 请求失败: %v", err)
	}
	if reviewURL, _ := resp["reviewurl"].(string); reviewURL != "" {
		ck, _ := resp["creditkey"].(string)
		devSign := xlGenerateDeviceSign(x.profile, x.deviceId)
		return "", fmt.Errorf("迅雷账号触发安全验证（review），请用手机迅雷 APP 登录该账号一次后重试，或打开链接完成验证: %s&deviceid=%s (creditkey=%s)", reviewURL, devSign, ck)
	}
	if errMsg, _ := resp["error"].(string); errMsg != "" && errMsg != "success" {
		return "", fmt.Errorf("CoreLogin 失败: %v", resp)
	}
	sessionID, _ = resp["sessionID"].(string)
	if sessionID == "" {
		return "", fmt.Errorf("CoreLogin 未返回 sessionID: %v", resp)
	}
	return sessionID, nil
}

// LoginWithCredentials 账号密码登录（保留，前端不再调用）。
func (x *XunleiPanService) LoginWithCredentials(username, password, creditkey string) (XunleiTokenData, error) {
	x.deviceId = deriveDeviceID(username, password)
	sessionID, err := x.coreLogin(username, password, creditkey)
	if err != nil {
		return XunleiTokenData{}, err
	}
	captchaToken, err := x.refreshCaptchaTokenInLogin("POST:/v1/auth/signin/token", username)
	if err != nil {
		return XunleiTokenData{}, fmt.Errorf("获取验证码令牌失败: %v", err)
	}
	body := map[string]interface{}{
		"client_id":     x.profile.ClientID,
		"client_secret": x.profile.ClientSecret,
		"provider":      xlSignProvider,
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
	log.Printf("迅雷账号 %s 登录成功（账号密码 CoreLogin）", username)
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

// ============================================================================
// refresh_token 登录（当前主用方式，跳过 CoreLogin/review）
// ============================================================================

// LoginByRefreshToken 用 refresh_token 直接登录。
func (x *XunleiPanService) LoginByRefreshToken(refreshToken string) (XunleiTokenData, error) {
	if refreshToken == "" {
		return XunleiTokenData{}, fmt.Errorf("refresh_token 为空")
	}
	// refresh_token 登录无 username，按 token 派生稳定 deviceId
	x.deviceId = deriveDeviceID(refreshToken, "xlrefresh")
	x.SetHeader("x-device-id", x.deviceId)
	return x.GetAccessTokenByRefreshToken(refreshToken)
}

// ============================================================================
// 通用认证请求
// ============================================================================

// sendXunleiAuthRequest 发送迅雷认证类 POST 请求（captcha init / CoreLogin / signin/token / token 刷新）。
// extraHeaders 用于按接口附加特殊头（如 CoreLogin 的 UA、signin/token 的 X-Captcha-Token）。
func (x *XunleiPanService) sendXunleiAuthRequest(reqURL string, data map[string]interface{}, extraHeaders map[string]string) (map[string]interface{}, error) {
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", x.profile.UserAgent)
	req.Header.Set("X-Client-Id", x.profile.ClientID)
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
