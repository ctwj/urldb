package pan

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
	"github.com/dustinxie/ecc"
)

// ============================================================================
// 阿里云盘服务（refresh_token 授权路线，web 接口免签名）。
// 依据：specs/008-aliyun-pan-transfer/research.md（R1~R14）、contracts/alipan-api-contract.md。
// 与夸克/UC 同构（统一 PanService 接口），但凭证走 refresh_token 换 access_token：
//   - Cks.Ck       = refresh_token（主凭证，运营人员填入）
//   - Cks.Extra    = JSON(AlipanExtraData)（运行期 access_token/drive_id/urldb 目录等）
// 日志约定：utils.Debug/Info/Error，严禁打印 refresh_token/access_token 明文（FR-001）。
// ============================================================================

const (
	alipanTokenURL    = "https://auth.alipan.com/v2/account/token" // refresh_token 换 access_token
	alipanAPIBase     = "https://api.aliyundrive.com"
	alipanUrldbFolder = "urldb"            // FR-014 固定转存目录名
	alipanMinInterval = 800 * time.Millisecond // per-account 最小请求间隔（防风控，FR-015）
	alipanMaxRetry    = 3                  // 风控退避重试上限（SC-006）
	// ECC 动态签名（移植 OpenList aliyundrive）：secp256k1 密钥对 + sign(secpAppID:deviceID:userID:0)
	alipanSecpAppID        = "5dde4e1bdf9e4966b387ba58f4b3fdc3"
	alipanCreateSessionURL = alipanAPIBase + "/users/v1/users/device/create_session"
)

// AlipanExtraData 阿里云盘运行期数据，JSON 序列化后存入 Cks.Extra（data-model.md §2.1）
type AlipanExtraData struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`   // 最新轮换值（与 Cks.Ck 冗余，便于单字段读写）
	DriveID       string `json:"drive_id"`        // 网盘 ID（去硬编码，research R2）
	ExpiresAt     int64  `json:"expires_at"`      // access_token 过期 unix 秒
	UrldbFolderID string `json:"urldb_folder_id"` // urldb 目录 file_id（懒初始化，research R5）
	DeviceID      string `json:"device_id"`       // 设备指纹（x-device-id）
	UserID        string `json:"user_id"`         // 阿里云盘 user_id（从 access_token JWT 解析，动态签名需要）
	PrivateKey    string `json:"private_key"`     // secp256k1 私钥 hex（动态签名用，移植 OpenList）
}

// alipanLimiter per-account 简单限速器（自实现，避免引入 x/time 升级 Go 版本破坏 Docker 构建）。
type alipanLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     time.Time
}

func newAlipanLimiter(interval time.Duration) *alipanLimiter {
	return &alipanLimiter{interval: interval}
}

// Wait 阻塞至满足最小请求间隔
func (l *alipanLimiter) Wait() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if !l.last.IsZero() {
		if d := l.interval - now.Sub(l.last); d > 0 {
			time.Sleep(d)
		}
	}
	l.last = time.Now()
}

// 全局 per-account 限速器注册表（跨 service 实例共享，按 refresh_token 索引）
var (
	alipanLimitersMu sync.Mutex
	alipanLimiters   = make(map[string]*alipanLimiter)
)

func getAlipanLimiter(key string) *alipanLimiter {
	if key == "" {
		key = "_default"
	}
	alipanLimitersMu.Lock()
	defer alipanLimitersMu.Unlock()
	l, ok := alipanLimiters[key]
	if !ok {
		l = newAlipanLimiter(alipanMinInterval)
		alipanLimiters[key] = l
	}
	return l
}

// AlipanService 阿里云盘服务
type AlipanService struct {
	*BasePanService
	configMutex sync.RWMutex

	// 凭证与运行期数据（SetCKSRepository 注入）
	refreshToken string
	extra        AlipanExtraData
	cksRepo      repo.CksRepository
	cksEntity    entity.Cks
	hasRepo      bool

	tokenMu sync.Mutex
	limiter *alipanLimiter
}

// NewAlipanService 创建阿里云盘服务（每次新建实例；per-account 限速器在 SetCKSRepository 时按账号绑定）
func NewAlipanService(config *PanConfig) *AlipanService {
	s := &AlipanService{
		BasePanService: NewBasePanService(config),
	}
	s.SetHeaders(map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Content-Type":    "application/json;charset=UTF-8",
		"Origin":          "https://www.alipan.com",
		"Referer":         "https://www.alipan.com/",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"X-Canary":        "client=web,app=adrive,version=v6.8.12", // 网页版最新版本（createShare 等需最新版本校验）
	})
	s.UpdateConfig(config)
	return s
}

// GetServiceType 获取服务类型
func (a *AlipanService) GetServiceType() ServiceType { return Alipan }

// UpdateConfig 更新配置（线程安全）
func (a *AlipanService) UpdateConfig(config *PanConfig) {
	if config == nil {
		return
	}
	a.configMutex.Lock()
	defer a.configMutex.Unlock()
	a.config = config
}

func (a *AlipanService) configValue() *PanConfig {
	a.configMutex.RLock()
	defer a.configMutex.RUnlock()
	return a.config
}

// SetCKSRepository 注入账号凭证仓储（research R4）
// 从 cks.Ck(refresh_token) 与 cks.Extra(AlipanExtraData) 解析凭证，并绑定 per-account 限速器。
func (a *AlipanService) SetCKSRepository(cksRepo repo.CksRepository, cks entity.Cks) {
	a.cksRepo = cksRepo
	a.cksEntity = cks
	a.hasRepo = cksRepo != nil
	a.refreshToken = cks.Ck
	a.extra = AlipanExtraData{}
	if cks.Extra != "" {
		_ = json.Unmarshal([]byte(cks.Extra), &a.extra)
	}
	if a.extra.RefreshToken == "" && a.refreshToken != "" {
		a.extra.RefreshToken = a.refreshToken
	}
	a.limiter = getAlipanLimiter(a.refreshToken)
	utils.Debug("[Alipan] SetCKSRepository accountID=%d driveID=%s hasAT=%v", cks.ID, a.extra.DriveID, a.extra.AccessToken != "")
}

// ============================================================================
// 令牌刷新与统一请求（research R1/R11）
// ============================================================================

// refreshToken 用 refresh_token 换 access_token（轮换，research R1）
// 刷新成功后回写 Cks.Extra；refresh_token 失效时标记账号（FR-016）。
func (a *AlipanService) refreshAccessToken() error {
	a.tokenMu.Lock()
	defer a.tokenMu.Unlock()

	rt := a.refreshToken
	if rt == "" {
		rt = a.extra.RefreshToken
	}
	if rt == "" {
		return fmt.Errorf("refresh_token 为空")
	}

	a.SetHeader("Authorization", "") // token 端点匿名，清空旧 Authorization
	respData, err := a.HTTPPost(alipanTokenURL, map[string]interface{}{
		"refresh_token": rt,
		"grant_type":    "refresh_token",
	}, nil)
	if err != nil {
		return fmt.Errorf("刷新 token 请求失败: %v", err)
	}

	var r struct {
		AccessToken    string `json:"access_token"`
		RefreshToken   string `json:"refresh_token"`
		ExpiresIn      int64  `json:"expires_in"`
		DefaultDriveID string `json:"default_drive_id"`
		Code           string `json:"code"`
		Message        string `json:"message"`
	}
	if err := json.Unmarshal(respData, &r); err != nil {
		return fmt.Errorf("解析 token 响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}
	if r.Code != "" || r.AccessToken == "" {
		a.markInvalid(fmt.Sprintf("refresh_token 失效: %s %s", r.Code, r.Message))
		return fmt.Errorf("refresh_token 失效: %s %s", r.Code, r.Message)
	}

	a.extra.AccessToken = r.AccessToken
	a.extra.RefreshToken = r.RefreshToken
	a.refreshToken = r.RefreshToken
	a.extra.ExpiresAt = time.Now().Unix() + r.ExpiresIn - 60 // 提前 60s 续期
	// 从 access_token JWT 解析 userId（动态签名 sign 需要）
	if uid := parseAlipanUserID(r.AccessToken); uid != "" {
		a.extra.UserID = uid
	}
	// token 响应的 default_drive_id 是真正的默认盘（对应网页"全部文件"）。
	// 对齐 aligo：BaseAligo.default_drive_id = token.default_drive_id（而非 /v2/user/get 的字段）。
	if r.DefaultDriveID != "" {
		if a.extra.DriveID != "" && a.extra.DriveID != r.DefaultDriveID {
			a.extra.UrldbFolderID = "" // drive 变了，清空旧 urldb 缓存，在新 drive 重新定位
		}
		a.extra.DriveID = r.DefaultDriveID
	}
	a.limiter = getAlipanLimiter(a.refreshToken)
	a.saveExtra()
	// token 刷新后设备会话失效，主动重建（createShare 等敏感接口需有效会话 + 动态签名）
	if serr := a.createSession(); serr != nil {
		utils.Debug("[Alipan] 刷新后 createSession 失败（不阻断，后续请求遇设备错误仍会重试）: %v", serr)
	}
	utils.Debug("[Alipan] token 刷新成功 expires_in=%d driveID=%s", r.ExpiresIn, a.extra.DriveID)
	return nil
}

// ensureAccessToken 确保有效 token（过期自动刷新）
func (a *AlipanService) ensureAccessToken() error {
	if a.extra.AccessToken != "" && time.Now().Unix() < a.extra.ExpiresAt {
		return nil
	}
	return a.refreshAccessToken()
}

// ensureDeviceID 确保账号有一个稳定的设备指纹（写入 Extra），用于 x-device-id 请求头
func (a *AlipanService) ensureDeviceID() string {
	if a.extra.DeviceID != "" {
		return a.extra.DeviceID
	}
	a.extra.DeviceID = newAlipanDeviceID()
	a.saveExtra()
	return a.extra.DeviceID
}

// newAlipanDeviceID 生成 32 位 hex 设备指纹（对齐 aligo 的 uuid.uuid4().hex，无连字符）
func newAlipanDeviceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano()) // 极端兜底
	}
	return fmt.Sprintf("%x", b)
}

// ============================================================================
// ECC 动态签名工具（移植 OpenList aliyundrive help.go / util.go sign）
// ============================================================================

// newAlipanPrivateKey 生成 secp256k1 私钥
func newAlipanPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(ecc.P256k1(), rand.Reader)
}

func alipanPrivateKeyFromHex(h string) (*ecdsa.PrivateKey, error) {
	data, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	curve := ecc.P256k1()
	x, y := curve.ScalarBaseMult(data)
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: curve, X: x, Y: y},
		D:         new(big.Int).SetBytes(data),
	}, nil
}

func alipanPrivateKeyToHex(k *ecdsa.PrivateKey) string {
	return hex.EncodeToString(k.D.Bytes())
}

func alipanPublicKeyToHex(pub *ecdsa.PublicKey) string {
	x := pub.X.Bytes()
	for len(x) < 32 {
		x = append([]byte{0}, x...)
	}
	y := pub.Y.Bytes()
	for len(y) < 32 {
		y = append([]byte{0}, y...)
	}
	return hex.EncodeToString(append(x, y...))
}

// parseAlipanUserID 从 access_token(JWT) 的 payload 解析 userId（动态签名 sign 需要）
func parseAlipanUserID(accessToken string) string {
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return ""
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}
	var claims struct {
		UserID string `json:"userId"`
	}
	if json.Unmarshal(payload, &claims) == nil {
		return claims.UserID
	}
	return ""
}

// ensurePrivateKey 确保账号有 secp256k1 私钥（持久化到 Extra），返回私钥
func (a *AlipanService) ensurePrivateKey() *ecdsa.PrivateKey {
	if a.extra.PrivateKey != "" {
		if k, err := alipanPrivateKeyFromHex(a.extra.PrivateKey); err == nil {
			return k
		}
	}
	k, err := newAlipanPrivateKey()
	if err != nil {
		return nil
	}
	a.extra.PrivateKey = alipanPrivateKeyToHex(k)
	a.saveExtra()
	return k
}

// alipanSign 用私钥对 secpAppID:deviceID:userID:0 的 sha256 做 ECC 签名（移植 OpenList sign）
func (a *AlipanService) alipanSign() string {
	singdata := fmt.Sprintf("%s:%s:%s:%d", alipanSecpAppID, a.ensureDeviceID(), a.extra.UserID, 0)
	hash := sha256.Sum256([]byte(singdata))
	pk := a.ensurePrivateKey()
	if pk == nil {
		return ""
	}
	data, _ := ecc.SignBytes(pk, hash[:], ecc.RecID|ecc.LowerS)
	return hex.EncodeToString(data)
}

// markInvalid 标记账号失效并写日志（FR-016；不主动推送，由运维巡检）
func (a *AlipanService) markInvalid(reason string) {
	utils.Error("[Alipan] 账号标记失效 accountID=%d reason=%s", a.cksEntity.ID, reason)
	if !a.hasRepo || a.cksEntity.ID == 0 {
		return
	}
	a.cksEntity.IsValid = false
	if err := a.cksRepo.UpdateWithAllFields(&a.cksEntity); err != nil {
		utils.Error("[Alipan] 标记失效回写失败 accountID=%d err=%v", a.cksEntity.ID, err)
	}
}

// saveExtra 回写运行期数据到 Cks.Extra（research R4）
func (a *AlipanService) saveExtra() {
	if !a.hasRepo || a.cksEntity.ID == 0 {
		return
	}
	b, _ := json.Marshal(a.extra)
	a.cksEntity.Extra = string(b)
	if a.extra.RefreshToken != "" {
		a.cksEntity.Ck = a.extra.RefreshToken // 同步轮换后的 refresh_token
	}
	if err := a.cksRepo.UpdateWithAllFields(&a.cksEntity); err != nil {
		utils.Error("[Alipan] saveExtra 回写失败 accountID=%d err=%v", a.cksEntity.ID, err)
	}
}

// alipanRequest 统一请求：token 自动续期 + per-account 限速 + 风控退避重试（research R6/R11）
// 解析响应体 code：token 类错误→刷新重试；风控类→指数退避重试；其余→直接返回错误。
func (a *AlipanService) alipanRequest(method, url string, body interface{}, extraHeaders map[string]string) ([]byte, error) {
	if err := a.ensureAccessToken(); err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt <= alipanMaxRetry; attempt++ {
		if a.limiter != nil {
			a.limiter.Wait()
		}
		a.SetHeader("Authorization", "Bearer "+a.extra.AccessToken)
		a.SetHeader("x-device-id", a.ensureDeviceID())
		a.SetHeader("x-signature", a.alipanSign()) // ECC 动态签名（移植 OpenList）
		for k, v := range extraHeaders {
			a.SetHeader(k, v)
		}

		var respData []byte
		var err error
		if method == "GET" {
			respData, err = a.HTTPGet(url, nil)
		} else {
			respData, err = a.HTTPPost(url, body, nil)
		}

		if err == nil {
			var e struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}
			_ = json.Unmarshal(respData, &e)
			if e.Code == "" {
				return respData, nil // 成功
			}
			if isAlipanTokenErr(e.Code) {
				utils.Debug("[Alipan] token 错误，刷新后重试 code=%s", e.Code)
				if rerr := a.refreshAccessToken(); rerr != nil {
					return nil, rerr
				}
				lastErr = fmt.Errorf("%s: %s", e.Code, e.Message)
				continue
			}
			if isAlipanRateLimitCode(e.Code) {
				utils.Debug("[Alipan] 风控限流，退避重试 code=%s attempt=%d", e.Code, attempt)
				backoffSleep(attempt)
				lastErr = fmt.Errorf("%s: %s", e.Code, e.Message)
				continue
			}
			return nil, fmt.Errorf("%s: %s", e.Code, e.Message)
		}

		// HTTP 层错误（非 2xx）
		if isAlipanDeviceErr(err.Error()) {
			// 设备会话失效（DeviceSessionSignatureInvalid/UserDeviceOffline/not found device），createSession 后重试
			utils.Debug("[Alipan] 设备会话失效，createSession 重试: %v", err)
			if cerr := a.createSession(); cerr != nil {
				return nil, cerr
			}
			lastErr = err
			continue
		}
		if isAlipanRateLimitErr(err.Error()) {
			utils.Debug("[Alipan] HTTP 风控，退避重试 attempt=%d err=%v", attempt, err)
			backoffSleep(attempt)
			lastErr = err
			continue
		}
		return nil, err
	}
	return nil, fmt.Errorf("阿里云盘请求重试耗尽: %v", lastErr)
}

func isAlipanTokenErr(code string) bool {
	c := strings.ToLower(code)
	return strings.Contains(c, "accesstokeninvalid") ||
		strings.Contains(c, "accesstokenexpired") ||
		strings.Contains(c, "tokenexpired") ||
		strings.Contains(c, "tokeninvalid")
}

func isAlipanRateLimitCode(code string) bool {
	c := strings.ToLower(code)
	return strings.Contains(c, "toomanyrequests") ||
		strings.Contains(c, "trafficlimit") ||
		strings.Contains(c, "ratelimit") ||
		strings.Contains(c, "requestfrequencylimit")
}

func isAlipanRateLimitErr(s string) bool {
	c := strings.ToLower(s)
	return strings.Contains(c, "429") ||
		strings.Contains(c, "toomany") ||
		strings.Contains(c, "rate limit") ||
		strings.Contains(c, "traffic")
}

func backoffSleep(attempt int) {
	d := time.Duration(1<<attempt) * time.Second // 1s, 2s, 4s
	if d > 8*time.Second {
		d = 8 * time.Second
	}
	time.Sleep(d)
}

// isAlipanDeviceErr 判断是否设备会话失效（需 createSession 重建），对齐 aligo request 的 400/401 设备错误处理
func isAlipanDeviceErr(s string) bool {
	c := strings.ToLower(s)
	return strings.Contains(c, "devicesessionsignatureinvalid") ||
		strings.Contains(c, "userdeviceoffline") ||
		strings.Contains(c, "not found device")
}

// createSession 建立设备会话（移植 OpenList：用账号自己的 secp256k1 公钥 + refreshToken + 动态签名）
func (a *AlipanService) createSession() error {
	if err := a.ensureAccessToken(); err != nil {
		return err
	}
	pk := a.ensurePrivateKey()
	a.SetHeader("Authorization", "Bearer "+a.extra.AccessToken)
	a.SetHeader("x-device-id", a.ensureDeviceID())
	a.SetHeader("x-signature", a.alipanSign())
	body := map[string]interface{}{
		"deviceName":   "urldb",
		"modelName":    "SM-G9810",
		"nonce":        0,
		"pubKey":       alipanPublicKeyToHex(&pk.PublicKey),
		"refreshToken": a.refreshToken,
	}
	if _, err := a.HTTPPost(alipanCreateSessionURL, body, nil); err != nil {
		utils.Warn("[Alipan] createSession 失败 deviceID=%s userID=%s err=%v", a.extra.DeviceID, a.extra.UserID, err)
		return fmt.Errorf("createSession 失败: %v", err)
	}
	utils.Info("[Alipan] createSession 成功 deviceID=%s userID=%s", a.extra.DeviceID, a.extra.UserID)
	return nil
}

// ============================================================================
// GetUserInfo 添加账号 / 刷新容量入口（research R13）
// ============================================================================

// GetUserInfo 传入的 *ck 为 refresh_token（Cks.Ck 或表单输入）。
// 用 refresh_token 换 access_token，调 /v2/user/get，返回容量/会员/昵称，
// 并经 UserInfo.ExtraData 携带 AlipanExtraData（含 access_token/drive_id）供首次持久化。
func (a *AlipanService) GetUserInfo(ck *string) (*UserInfo, error) {
	if ck != nil && *ck != "" {
		a.refreshToken = *ck
		a.extra.RefreshToken = *ck
		a.limiter = getAlipanLimiter(a.refreshToken)
	}
	if err := a.refreshAccessToken(); err != nil {
		return nil, err
	}

	a.SetHeader("Authorization", "Bearer "+a.extra.AccessToken)
	respData, err := a.HTTPPost(alipanAPIBase+"/v2/user/get", map[string]interface{}{}, nil)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	var r struct {
		Code            string `json:"code"`
		Message         string `json:"message"`
		NickName        string `json:"nick_name"`
		VipStatus       string `json:"vip_status"`
		Vip             string `json:"vip"`
		DefaultDriveID  string `json:"default_drive_id"`
		ResourceDriveID string `json:"resource_drive_id"`
		BackupDriveID   string `json:"backup_drive_id"`
	}
	if err := json.Unmarshal(respData, &r); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v bodyHead=%s", err, headSnippet(respData))
	}
	if r.Code != "" {
		return nil, fmt.Errorf("获取用户信息错误: %s %s", r.Code, r.Message)
	}

	// drive 由 token 刷新返回的 default_drive_id 决定（对齐 aligo：BaseAligo.default_drive_id = token.default_drive_id）。
	// user/get 这里仅记录各 drive 供参考，不覆盖 extra.DriveID。
	utils.Info("[Alipan] drive 列表（user/get）default=%s resource=%s backup=%s 当前选用(token default)=%s",
		r.DefaultDriveID, r.ResourceDriveID, r.BackupDriveID, a.extra.DriveID)

	// 容量走专用接口 driveCapacityDetails（/v2/user/get 不返回可靠容量；OpenList 验证路径）
	total, used, capErr := a.getDriveCapacity()
	if capErr != nil {
		utils.Debug("[Alipan] 获取容量失败（继续，置0）: %v", capErr)
	}

	vip := r.VipStatus == "vip" || r.Vip == "vip"
	extraJSON, _ := json.Marshal(a.extra)
	utils.Debug("[Alipan] GetUserInfo 成功 nick=%s driveID=%s total=%d used=%d vip=%v",
		r.NickName, a.extra.DriveID, total, used, vip)

	return &UserInfo{
		Username:    r.NickName,
		VIPStatus:   vip,
		UsedSpace:   used,
		TotalSpace:  total,
		ServiceType: "alipan",
		ExtraData:   string(extraJSON),
	}, nil
}

// getDriveCapacity 调专用接口获取账号容量（字节数，顶层 drive_total_size/drive_used_size）
func (a *AlipanService) getDriveCapacity() (total, used int64, err error) {
	respData, err := a.alipanRequest("POST", alipanAPIBase+"/adrive/v1/user/driveCapacityDetails", map[string]interface{}{}, nil)
	if err != nil {
		return 0, 0, err
	}
	var r struct {
		DriveTotalSize int64 `json:"drive_total_size"`
		DriveUsedSize  int64 `json:"drive_used_size"`
	}
	if err := json.Unmarshal(respData, &r); err != nil {
		return 0, 0, fmt.Errorf("解析容量响应失败: %v bodyHead=%s", err, headSnippet(respData))
	}
	return r.DriveTotalSize, r.DriveUsedSize, nil
}

// GetUserInfoByEntity 从 entity.Cks 获取用户信息（research R13）
func (a *AlipanService) GetUserInfoByEntity(cks entity.Cks) (*UserInfo, error) {
	if cks.Ck == "" {
		return nil, fmt.Errorf("阿里云盘账号 refresh_token 为空")
	}
	a.SetCKSRepository(a.cksRepo, cks)
	ck := cks.Ck
	return a.GetUserInfo(&ck)
}

// ============================================================================
// urldb 目录定位 / 创建（research R5；FR-014）
// ============================================================================

// ensureUrldbFolder 确保账号网盘根目录下存在 urldb 文件夹，返回其 file_id（缓存到 Extra）
func (a *AlipanService) ensureUrldbFolder() (string, error) {
	if a.extra.UrldbFolderID != "" {
		return a.extra.UrldbFolderID, nil
	}
	if a.extra.DriveID == "" {
		return "", fmt.Errorf("drive_id 为空，无法定位 urldb 目录")
	}

	// 1. 先在根目录查找
	if fid := a.findFolderInRoot(alipanUrldbFolder); fid != "" {
		a.extra.UrldbFolderID = fid
		a.saveExtra()
		return fid, nil
	}

	// 2. 不存在则创建（check_name_mode refuse 避免重名）
	createResp, err := a.alipanRequest("POST", alipanAPIBase+"/adrive/v2/file/createWithFolders", map[string]interface{}{
		"check_name_mode": "refuse",
		"drive_id":        a.extra.DriveID,
		"parent_file_id":  "root",
		"name":            alipanUrldbFolder,
		"type":            "folder",
	}, nil)
	var cr struct {
		FileID  string `json:"file_id"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err == nil {
		_ = json.Unmarshal(createResp, &cr)
	}
	if err != nil || cr.FileID == "" {
		// 创建失败或无 file_id（含 refuse 已存在、并发创建等）→ fallback 重新列根查找
		if err != nil {
			utils.Debug("[Alipan] 创建 urldb 目录返回错误，重新查找: %v", err)
		}
		if fid := a.findFolderInRoot(alipanUrldbFolder); fid != "" {
			a.extra.UrldbFolderID = fid
			a.saveExtra()
			return fid, nil
		}
		if err != nil {
			return "", fmt.Errorf("创建 urldb 目录失败: %v", err)
		}
		return "", fmt.Errorf("创建 urldb 目录无 file_id: %s %s", cr.Code, cr.Message)
	}
	a.extra.UrldbFolderID = cr.FileID
	a.saveExtra()
	utils.Debug("[Alipan] 创建 urldb 目录 file_id=%s", cr.FileID)
	return cr.FileID, nil
}

// findFolderInRoot 在账号根目录下查找指定名称的文件夹/文件，返回 file_id（未找到或出错返回空）
func (a *AlipanService) findFolderInRoot(name string) string {
	respData, err := a.alipanRequest("POST", alipanAPIBase+"/adrive/v3/file/list", map[string]interface{}{
		"drive_id":        a.extra.DriveID,
		"parent_file_id":  "root",
		"limit":           200,
		"order_by":        "name",
		"order_direction": "ASC",
		"fields":          "*",
	}, nil)
	if err != nil {
		return ""
	}
	var lr struct {
		Items []struct {
			FileID string `json:"file_id"`
			Name   string `json:"name"`
		} `json:"items"`
	}
	_ = json.Unmarshal(respData, &lr)
	for _, it := range lr.Items {
		if it.Name == name {
			return it.FileID
		}
	}
	return ""
}

// ============================================================================
// Transfer 转存分享链接（research R6；FR-004/005/006/007/014）
// 链路：匿名取分享 → (校验模式直接返回) → share_token(含提取码) → urldb 目录 → batch copy → 永久再分享
// ============================================================================

func (a *AlipanService) Transfer(shareID string) (*TransferResult, error) {
	config := a.configValue()
	isType := 0
	if config != nil {
		isType = config.IsType
	}
	utils.Info("[Alipan] 开始转存 shareID=%s isType=%d(0=转存,1=校验)", shareID, isType)

	if err := a.ensureAccessToken(); err != nil {
		return ErrorResult(fmt.Sprintf("获取 access_token 失败: %v", err)), nil
	}

	// 1. 匿名获取分享文件
	shareInfo, err := a.getShareByAnonymous(shareID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取分享信息失败: %v", err)), nil
	}

	// 校验模式：仅返回标题，不转存（FR-007）
	if isType == 1 {
		shareURL := ""
		if config != nil {
			shareURL = config.URL
		}
		utils.Info("[Alipan] 校验模式完成（不转存）shareID=%s title=%s", shareID, shareInfo.ShareName)
		return SuccessResult("检验成功", map[string]interface{}{
			"title":    shareInfo.ShareName,
			"shareUrl": shareURL,
		}), nil
	}

	// 2. share_token（含提取码，FR-005）
	sharePwd := ""
	if config != nil {
		sharePwd = config.Code
	}
	shareToken, err := a.getShareToken(shareID, sharePwd)
	if err != nil {
		return ErrorResult(err.Error()), nil
	}

	// 3. urldb 目录（FR-014）
	urldbID, err := a.ensureUrldbFolder()
	if err != nil {
		return ErrorResult(fmt.Sprintf("定位 urldb 目录失败: %v", err)), nil
	}

	// 4. 批量转存到 urldb 目录
	fileIDs := make([]string, 0, len(shareInfo.FileInfos))
	for _, f := range shareInfo.FileInfos {
		if f.FileID != "" {
			fileIDs = append(fileIDs, f.FileID)
		}
	}
	if len(fileIDs) == 0 {
		return ErrorResult("分享内无可转存文件"), nil
	}
	newFileIDs, err := a.batchCopy(shareID, fileIDs, urldbID, shareToken)
	if err != nil {
		// batchCopy 已返回具体原因（如"账号容量不足..."），直接透传，避免"转存失败: 转存失败:"嵌套
		return ErrorResult(err.Error()), nil
	}
	utils.Info("[Alipan] 转存到 urldb 完成 driveID=%s urldbFolderID=%s 转存文件数=%d 新fileIDs=%v", a.extra.DriveID, urldbID, len(newFileIDs), newFileIDs)

	// 5. 对【本账号 urldb 目录下的转存文件】创建永久再分享（FR-006/Q2）
	// 注意：newFileIDs 是 batchCopy 转存后生成的新文件 ID（位于本账号 urldb 目录），不是原分享的文件 ID
	shareRes, createErr := a.createShare(newFileIDs)
	shareURL := ""
	shareTitle := shareInfo.ShareName
	if createErr != nil {
		// createShare 受阿里云盘签名校验限制（403"升级版本"，需动态 x-signature），降级：
		// 文件已转存到本账号 urldb 目录（持久备份），用原分享 url 入库，保证转存流程可用。
		utils.Warn("[Alipan] 创建新分享失败，降级用原分享 url（文件已转存到 urldb 备份）: %v", createErr)
		if cfg := a.configValue(); cfg != nil {
			shareURL = cfg.URL
		}
	} else {
		shareURL = shareRes.ShareURL
		shareTitle = shareRes.ShareTitle
	}

	fid := strings.Join(newFileIDs, ",")
	utils.Info("[Alipan] 转存完成 shareID=%s shareUrl=%s title=%s fid=%s 新分享创建=%v", shareID, shareURL, shareTitle, fid, createErr == nil)
	return SuccessResult("转存成功", map[string]interface{}{
		"shareUrl": shareURL,
		"title":    shareTitle,
		"fid":      fid,
	}), nil
}

// getShareByAnonymous 匿名获取分享文件列表
func (a *AlipanService) getShareByAnonymous(shareID string) (*alipanShareInfo, error) {
	respData, err := a.alipanRequest("POST", alipanAPIBase+"/adrive/v2/share_link/get_share_by_anonymous", map[string]interface{}{
		"share_id": shareID,
	}, nil)
	if err != nil {
		return nil, err
	}
	var r alipanShareInfo
	_ = json.Unmarshal(respData, &r)
	if r.ShareName == "" && len(r.FileInfos) == 0 {
		return nil, fmt.Errorf("分享不存在或已失效")
	}
	return &r, nil
}

// getShareToken 获取 share_token（含提取码；提取码错误转译，FR-005）
func (a *AlipanService) getShareToken(shareID, sharePwd string) (string, error) {
	respData, err := a.alipanRequest("POST", alipanAPIBase+"/v2/share_link/get_share_token", map[string]interface{}{
		"share_id":  shareID,
		"share_pwd": sharePwd,
	}, nil)
	if err != nil {
		return "", err
	}
	var r struct {
		ShareToken string `json:"share_token"`
		Code       string `json:"code"`
		Message    string `json:"message"`
	}
	_ = json.Unmarshal(respData, &r)
	if r.ShareToken == "" {
		cl := strings.ToLower(r.Code + " " + r.Message)
		if strings.Contains(cl, "pwd") || strings.Contains(cl, "密码") || strings.Contains(cl, "提取码") {
			return "", fmt.Errorf("提取码错误: %s", r.Message)
		}
		return "", fmt.Errorf("获取 share_token 失败: %s %s", r.Code, r.Message)
	}
	return r.ShareToken, nil
}

// batchCopy 批量转存分享文件到 urldb 目录（去硬编码 drive_id，research R6）
func (a *AlipanService) batchCopy(shareID string, fileIDs []string, toParentFolderID, shareToken string) ([]string, error) {
	requests := make([]map[string]interface{}, 0, len(fileIDs))
	for i, fid := range fileIDs {
		requests = append(requests, map[string]interface{}{
			"headers": map[string]string{"Content-Type": "application/json"},
			"method":  "POST",
			"id":      fmt.Sprintf("%d", i),
			"url":     "/file/copy",
			"body": map[string]interface{}{
				"share_id":          shareID,
				"file_id":           fid,
				"to_drive_id":       a.extra.DriveID,
				"to_parent_file_id": toParentFolderID,
				"auto_rename":       true,
			},
		})
	}
	respData, err := a.alipanRequest("POST", alipanAPIBase+"/adrive/v2/batch", map[string]interface{}{
		"requests": requests,
		"resource": "file",
	}, map[string]string{"X-Share-Token": shareToken})
	if err != nil {
		return nil, err
	}
	var r struct {
		Responses []struct {
			Status int `json:"status"`
			Body   struct {
				Code    string `json:"code"`
				Message string `json:"message"`
				FileID  string `json:"file_id"`
			} `json:"body"`
		} `json:"responses"`
	}
	_ = json.Unmarshal(respData, &r)
	newIDs := make([]string, 0, len(r.Responses))
	for _, resp := range r.Responses {
		if resp.Body.Code != "" {
			cl := strings.ToLower(resp.Body.Code + " " + resp.Body.Message)
			if strings.Contains(cl, "space") || strings.Contains(cl, "capacity") || strings.Contains(cl, "容量") || strings.Contains(cl, "exceeded the limit") {
				return nil, fmt.Errorf("账号容量不足（资源盘已达上限），请清理空间或更换账号")
			}
			return nil, fmt.Errorf("转存文件失败: %s %s", resp.Body.Code, resp.Body.Message)
		}
		if resp.Body.FileID != "" {
			newIDs = append(newIDs, resp.Body.FileID)
		}
	}
	if len(newIDs) == 0 {
		return nil, fmt.Errorf("转存完成但未获得新文件 ID")
	}
	return newIDs, nil
}

// createShare 创建永久分享（expiration:"" = 永久，FR-006）
func (a *AlipanService) createShare(fileIDs []string) (*alipanShareResult, error) {
	// createShare 需有效设备会话，主动建一次（幂等；createSession 内部带动态签名）
	if serr := a.createSession(); serr != nil {
		utils.Warn("[Alipan] createShare 前 createSession 失败: %v", serr)
	}
	respData, err := a.alipanRequest("POST", alipanAPIBase+"/adrive/v2/share_link/create", map[string]interface{}{
		"drive_id":     a.extra.DriveID,
		"file_id_list": fileIDs,
		"expiration":   "",
		"share_pwd":    "",
	}, map[string]string{
		// 补齐浏览器特征头（对齐网页版 cURL），规避 createShare 的非浏览器识别
		"sec-fetch-site":     "cross-site",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"sec-ch-ua-platform": `"Windows"`,
		"sec-ch-ua-mobile":   "?0",
		"dnt":                "1",
	})
	if err != nil {
		return nil, err
	}
	var r alipanShareResult
	_ = json.Unmarshal(respData, &r)
	if r.ShareURL == "" {
		return nil, fmt.Errorf("创建分享未返回链接")
	}
	return &r, nil
}

// ============================================================================
// GetFiles / DeleteFiles（research R8；FR-008）
// ============================================================================

// GetFiles 获取文件列表
func (a *AlipanService) GetFiles(pdirFid string) (*TransferResult, error) {
	if err := a.ensureAccessToken(); err != nil {
		return ErrorResult(fmt.Sprintf("获取 access_token 失败: %v", err)), nil
	}
	if pdirFid == "" {
		pdirFid = "root"
	}
	respData, err := a.alipanRequest("POST", alipanAPIBase+"/adrive/v3/file/list", map[string]interface{}{
		"drive_id":        a.extra.DriveID,
		"parent_file_id":  pdirFid,
		"limit":           100,
		"order_by":        "updated_at",
		"order_direction": "DESC",
		"fields":          "*",
	}, nil)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取文件列表失败: %v", err)), nil
	}
	var r struct {
		Items   []interface{} `json:"items"`
		Message string        `json:"message"`
	}
	_ = json.Unmarshal(respData, &r)
	return SuccessResult("获取成功", r.Items), nil
}

// DeleteFiles 删除文件（清理用；FR-008）。fid 可能是逗号分隔（转存时 join 存储）。
// 错误信息确保能被 cleanup_service.isFileNotExist 宽松匹配（research R8）。
func (a *AlipanService) DeleteFiles(fileList []string) (*TransferResult, error) {
	if len(fileList) == 0 {
		return ErrorResult("文件列表为空"), nil
	}
	if err := a.ensureAccessToken(); err != nil {
		return ErrorResult(fmt.Sprintf("获取 access_token 失败: %v", err)), nil
	}
	allIDs := make([]string, 0, len(fileList))
	for _, f := range fileList {
		for _, id := range strings.Split(f, ",") {
			id = strings.TrimSpace(id)
			if id != "" {
				allIDs = append(allIDs, id)
			}
		}
	}
	if len(allIDs) == 0 {
		return ErrorResult("文件列表为空"), nil
	}

	_, err := a.alipanRequest("POST", alipanAPIBase+"/adrive/v3/file/delete", map[string]interface{}{
		"drive_id":     a.extra.DriveID,
		"file_id_list": allIDs,
	}, nil)
	if err != nil {
		msg := err.Error()
		cl := strings.ToLower(msg)
		if strings.Contains(cl, "not found") || strings.Contains(cl, "not exist") || strings.Contains(cl, "不存在") {
			msg = "文件不存在"
		}
		return ErrorResult(fmt.Sprintf("删除文件失败: %s", msg)), nil
	}
	utils.Debug("[Alipan] 删除文件成功 count=%d", len(allIDs))
	return SuccessResult("删除成功", nil), nil
}

// ============================================================================
// 阿里云盘响应结构体
// ============================================================================

type alipanShareInfo struct {
	ShareName string `json:"share_name"`
	FileInfos []struct {
		FileID string `json:"file_id"`
	} `json:"file_infos"`
}

type alipanShareResult struct {
	ShareURL   string   `json:"share_url"`
	ShareTitle string   `json:"share_title"`
	FileIDList []string `json:"file_id_list"`
}
