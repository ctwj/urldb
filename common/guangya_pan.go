package pan

// ============================================================================
// 光鸭云盘（guangyapan.com）服务实现。
//
// 蓝本：demo/cloud-save-max-main 光鸭适配器（Python，生产验证）+ 用户抓包的创建分享
// 接口。接口契约见 specs/013-guangya-pan-integration/contracts/guangya-api.md，
// 实现决策见 specs/013-guangya-pan-integration/research.md（D1~D9）。
//
// 与既有平台的三项关键差异：
//  1. Bearer 令牌鉴权（非 Cookie）：凭据 kv 串存 Cks.Ck，401 自动刷新并回写（D2）；
//  2. 转存接口不返回新 fid：固定专用目录「转存」+ 前后快照 diff 按文件名对齐（D4）；
//  3. 响应无统一数值 code：按 msg/data 容错判定（contracts 通用约定）。
//
// 日志约定与 uc_pan.go 一致：诊断 utils.Debug、里程碑 utils.Info、失败 utils.Error，
// 严禁打印令牌明文。
// ============================================================================

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

const (
	// guangyaAPIUserResV2 用户资源 v2（文件列表、创建分享）
	guangyaAPIUserResV2 = "https://api.guangyapan.com/userres/v1"
	// guangyaAPIUserRes 用户资源 v1（目录/删除/任务/转存/匿名分享访问）
	guangyaAPIUserRes = "https://api.guangyapan.com/nd.bizuserres.s/v1"
	// guangyaAPIAssets 资产（容量/会员）
	guangyaAPIAssets = "https://api.guangyapan.com/assets/v1"
	// guangyaAccountAPI 账号（token 刷新、user/me）
	guangyaAccountAPI = "https://account.guangyapan.com/v1"
	guangyaWebOrigin  = "https://www.guangyapan.com"
	guangyaClientID   = "aMe-8VSlkrbQXpUR"
	guangyaUA         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36"
	// guangyaTransferDirName 转存固定专用目录（D4）：fid 对齐范围小、清理按目录回收。
	// 命名与阿里云盘/百度一致（urldb），保持全平台统一约定（2026-08-31 用户确认，由「转存」更名）
	guangyaTransferDirName = "urldb"
	// guangyaTaskPollMax / guangyaTaskPollInterval 转存任务轮询上限与间隔（D5）
	guangyaTaskPollMax      = 20
	guangyaTaskPollInterval = time.Second
	// guangyaDiffRetryMax 转存后目录 diff 对齐重试次数（1s 间隔）
	guangyaDiffRetryMax = 15
	// guangyaPageMax 分页聚合上限（200 页）
	guangyaPageMax = 200
)

// ============================================================================
// 全局节流（D7 / FR-013）：所有光鸭调用路径统一生效，进程级互斥
// ============================================================================

var (
	guangyaThrottleMu sync.Mutex
	guangyaLastReqAt  time.Time
)

// guangyaThrottle 保证相邻请求间隔 50-100ms 随机值（参考适配器安全水位）
func guangyaThrottle() {
	guangyaThrottleMu.Lock()
	defer guangyaThrottleMu.Unlock()
	interval := time.Duration(50+rand.Intn(50)) * time.Millisecond
	if !guangyaLastReqAt.IsZero() {
		if wait := interval - time.Since(guangyaLastReqAt); wait > 0 {
			time.Sleep(wait)
		}
	}
	guangyaLastReqAt = time.Now()
}

// generateGuangyaDid 生成设备标识：md5(随机16字节hex)
func generateGuangyaDid() string {
	sum := md5.Sum([]byte(guangyaRandHex(16)))
	return hex.EncodeToString(sum[:])
}

// guangyaRandHex 生成 n 字节随机数的 hex 串（长度 2n）
func guangyaRandHex(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return hex.EncodeToString(b)
}

// generateGuangyaTraceparent 生成 W3C traceparent 头（每请求随机）
func generateGuangyaTraceparent() string {
	return fmt.Sprintf("00-%s-%s-01", guangyaRandHex(16), guangyaRandHex(8))
}

// ============================================================================
// 凭据（D2）：kv 串解析/序列化，4 条解析规则见 data-model.md
// ============================================================================

// guangyaCredential 光鸭令牌凭据运行时
type guangyaCredential struct {
	accessToken  string
	refreshToken string
	deviceID     string
}

// parseGuangyaCredential 解析凭据串：完整 kv / 裸 access_token / 缺 device_id 自动生成 / 全空报错
func parseGuangyaCredential(input string) (guangyaCredential, error) {
	var c guangyaCredential
	text := strings.TrimSpace(input)
	if text == "" {
		return c, fmt.Errorf("光鸭凭据为空")
	}
	if !strings.Contains(text, "=") {
		// 裸 access_token（无 = 视为 token 本体）
		c.accessToken = text
	} else {
		for _, chunk := range strings.Split(text, ";") {
			chunk = strings.TrimSpace(chunk)
			if chunk == "" {
				continue
			}
			k, v, ok := strings.Cut(chunk, "=")
			if !ok {
				continue
			}
			k, v = strings.TrimSpace(k), strings.TrimSpace(v)
			switch k {
			case "access_token", "accessToken":
				c.accessToken = v
			case "refresh_token", "refreshToken":
				c.refreshToken = v
			case "device_id", "deviceId":
				c.deviceID = v
			}
		}
	}
	if c.accessToken == "" {
		return c, fmt.Errorf("光鸭凭据缺少 access_token")
	}
	if c.deviceID == "" {
		c.deviceID = generateGuangyaDid()
	}
	return c, nil
}

// serialize 规范化 kv（落库/回写格式；单一真相，research D2）
func (c guangyaCredential) serialize() string {
	parts := []string{"access_token=" + c.accessToken}
	if c.refreshToken != "" {
		parts = append(parts, "refresh_token="+c.refreshToken)
	}
	parts = append(parts, "device_id="+c.deviceID)
	return strings.Join(parts, ";")
}

// maskGuangyaMobile 手机号脱敏（138****0000；+86 前缀剥离；不足 7 位原样返回）
func maskGuangyaMobile(mobile string) string {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, mobile)
	if len(digits) == 13 && strings.HasPrefix(digits, "86") {
		digits = digits[2:] // 剥离国家码（86 + 11 位手机号）
	}
	if len(digits) < 7 {
		return mobile
	}
	return digits[:3] + "****" + digits[len(digits)-4:]
}

// ============================================================================
// 响应通用判定（contracts 通用约定：无统一数值 code，按 msg/data 容错）
// ============================================================================

func guangyaClean(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

// guangyaDeepGet 依次取首个非 nil 键
func guangyaDeepGet(payload map[string]interface{}, keys ...string) interface{} {
	if payload == nil {
		return nil
	}
	for _, k := range keys {
		if v, ok := payload[k]; ok && v != nil {
			return v
		}
	}
	return nil
}

// guangyaExtractMessage 提取错误消息
func guangyaExtractMessage(payload map[string]interface{}) string {
	for _, scope := range []map[string]interface{}{payload, guangyaAsMap(guangyaExtractData(payload))} {
		if scope == nil {
			continue
		}
		for _, k := range []string{"message", "msg", "errorMessage", "error_message", "detail"} {
			if s := guangyaClean(scope[k]); s != "" {
				return s
			}
		}
	}
	return ""
}

// guangyaExtractData 取 data 字段（无则原样）
func guangyaExtractData(payload map[string]interface{}) interface{} {
	if payload == nil {
		return nil
	}
	if d, ok := payload["data"]; ok {
		return d
	}
	return payload
}

func guangyaAsMap(v interface{}) map[string]interface{} {
	m, _ := v.(map[string]interface{})
	return m
}

// guangyaExtractItems 提取列表条目（data 为 list 或 data/list|items|records|rows|fileList|files）
func guangyaExtractItems(payload map[string]interface{}) []map[string]interface{} {
	data := guangyaExtractData(payload)
	appendFrom := func(v interface{}) []map[string]interface{} {
		arr, ok := v.([]interface{})
		if !ok {
			return nil
		}
		out := make([]map[string]interface{}, 0, len(arr))
		for _, item := range arr {
			if m, ok := item.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
		return out
	}
	if items := appendFrom(data); items != nil {
		return items
	}
	if dm := guangyaAsMap(data); dm != nil {
		for _, k := range []string{"list", "items", "records", "rows", "fileList", "files"} {
			if items := appendFrom(dm[k]); items != nil {
				return items
			}
		}
	}
	if pm := guangyaAsMap(payload); pm != nil {
		for _, k := range []string{"list", "items", "records", "rows"} {
			if items := appendFrom(pm[k]); items != nil {
				return items
			}
		}
	}
	return nil
}

func guangyaToBool(v interface{}) bool {
	switch x := v.(type) {
	case bool:
		return x
	case float64:
		return x != 0
	case string:
		switch strings.ToLower(strings.TrimSpace(x)) {
		case "1", "true", "yes", "y", "folder", "dir", "directory":
			return true
		}
	}
	return false
}

func guangyaToInt(v interface{}) int64 {
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int:
		return int64(x)
	case int64:
		return x
	case string:
		n, _ := strconv.ParseInt(strings.TrimSpace(x), 10, 64)
		return n
	}
	return 0
}

// guangyaIsSuccess 成功判定（contracts 通用约定）
func guangyaIsSuccess(payload map[string]interface{}) bool {
	if payload == nil {
		return false
	}
	if b, ok := payload["success"].(bool); ok && b {
		return true
	}
	if b, ok := payload["ok"].(bool); ok && b {
		return true
	}
	for _, k := range []string{"msg", "message"} {
		if strings.EqualFold(guangyaClean(payload[k]), "success") {
			return true
		}
	}
	for _, k := range []string{"code", "status", "result"} {
		v, ok := payload[k]
		if !ok || v == nil {
			continue
		}
		if n := guangyaToInt(v); n == 0 || n == 200 {
			return true
		}
		switch strings.ToLower(guangyaClean(v)) {
		case "0", "200", "ok", "success":
			return true
		}
	}
	if d, ok := payload["data"]; ok && d != nil && guangyaExtractMessage(payload) == "" {
		return true
	}
	return false
}

// guangyaHasExplicitFailure 显式失败判定（code/status/result 明确非 0/200）
func guangyaHasExplicitFailure(payload map[string]interface{}) bool {
	for _, k := range []string{"code", "status", "result"} {
		v, ok := payload[k]
		if !ok || v == nil {
			continue
		}
		if n := guangyaToInt(v); n != 0 && n != 200 {
			return true
		}
		text := strings.ToLower(guangyaClean(v))
		if text != "" && text != "0" && text != "200" && text != "ok" && text != "success" && text != "true" {
			return true
		}
	}
	return false
}

// guangyaRaiseIfFailed 失败时返回带上游消息的 error
func guangyaRaiseIfFailed(payload map[string]interface{}, fallback string) error {
	if guangyaIsSuccess(payload) {
		return nil
	}
	msg := guangyaExtractMessage(payload)
	if msg == "" {
		msg = fallback
	}
	return fmt.Errorf("%s", msg)
}

// ============================================================================
// HTTP 客户端：登录态 guangyaClient + 匿名包级函数
// ============================================================================

// guangyaHTTPError 携带状态码的 HTTP 错误（401 触发刷新重试）
type guangyaHTTPError struct {
	Status int
	Body   string
}

func (e *guangyaHTTPError) Error() string {
	body := e.Body
	if len(body) > 200 {
		body = body[:200]
	}
	return fmt.Sprintf("HTTP %d: %s", e.Status, body)
}

type guangyaClient struct {
	mu        sync.Mutex // 保护 cred/expiresAt（并发刷新）
	cred      guangyaCredential
	expiresAt *time.Time
	http      *http.Client
	cksRepo   repo.CksRepository // 令牌刷新回写（D2；可空=不回写）
	account   *entity.Cks
}

func newGuangyaClient(cred guangyaCredential) *guangyaClient {
	return &guangyaClient{
		cred: cred,
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

// setGuangyaRequestHeaders 设置公共请求头（contracts 通用约定）
func setGuangyaRequestHeaders(req *http.Request, deviceID, token string) {
	h := req.Header
	h.Set("accept", "application/json, text/plain, */*")
	h.Set("content-type", "application/json")
	h.Set("did", deviceID)
	h.Set("dt", "4")
	h.Set("origin", guangyaWebOrigin)
	h.Set("referer", guangyaWebOrigin+"/")
	h.Set("user-agent", guangyaUA)
	h.Set("traceparent", generateGuangyaTraceparent())
	if token != "" {
		h.Set("authorization", "Bearer "+token)
	}
}

// doJSON 登录态请求：节流 → 发送 → 401 刷新重试一次 → 解析 JSON dict
func (c *guangyaClient) doJSON(endpoint string, payload interface{}) (map[string]interface{}, error) {
	result, err := c.doJSONOnce(endpoint, payload, true)
	if err == nil {
		return result, nil
	}
	httpErr, ok := err.(*guangyaHTTPError)
	if !ok || httpErr.Status != http.StatusUnauthorized {
		return nil, err
	}
	// 401 → 刷新令牌重试一次（D2 状态机）
	if refreshErr := c.refreshAccessToken(); refreshErr != nil {
		return nil, fmt.Errorf("令牌已过期，请重新获取（刷新失败: %v）", refreshErr)
	}
	return c.doJSONOnce(endpoint, payload, false)
}

func (c *guangyaClient) doJSONOnce(endpoint string, payload interface{}, allowRefresh bool) (map[string]interface{}, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// 临期主动刷新（expiresAt 已知时）
	c.mu.Lock()
	cred := c.cred
	exp := c.expiresAt
	c.mu.Unlock()
	if allowRefresh && exp != nil && cred.refreshToken != "" && time.Now().After(*exp) {
		if err := c.refreshAccessToken(); err != nil {
			utils.Debug("[光鸭:CLIENT] 临期刷新失败（继续用当前令牌尝试）: %v", err)
		}
		c.mu.Lock()
		cred = c.cred
		c.mu.Unlock()
	}

	guangyaThrottle()
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	setGuangyaRequestHeaders(req, cred.deviceID, cred.accessToken)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &guangyaHTTPError{Status: resp.StatusCode, Body: string(respBody)}
	}
	var out map[string]interface{}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("响应解析失败: %v body=%s", err, headSnippet(respBody))
	}
	return out, nil
}

// refreshAccessToken 刷新令牌并回写账号凭据（contracts #11；x-* 头按契约原样）
func (c *guangyaClient) refreshAccessToken() error {
	c.mu.Lock()
	cred := c.cred
	c.mu.Unlock()
	if cred.refreshToken == "" {
		return fmt.Errorf("无可用的 refresh_token")
	}

	body, _ := json.Marshal(map[string]interface{}{
		"client_id":     guangyaClientID,
		"grant_type":    "refresh_token",
		"refresh_token": cred.refreshToken,
	})
	guangyaThrottle()
	req, err := http.NewRequest("POST", guangyaAccountAPI+"/auth/token", bytes.NewReader(body))
	if err != nil {
		return err
	}
	h := req.Header
	h.Set("accept", "*/*")
	h.Set("content-type", "application/json")
	h.Set("origin", guangyaWebOrigin)
	h.Set("referer", guangyaWebOrigin+"/")
	h.Set("user-agent", guangyaUA)
	h.Set("x-client-id", guangyaClientID)
	h.Set("x-client-version", "0.0.1")
	h.Set("x-device-id", cred.deviceID)
	h.Set("x-device-model", "chrome%2F147.0.0.0")
	h.Set("x-device-name", "PC-Chrome")
	h.Set("x-device-sign", "wdi10."+cred.deviceID+guangyaRandHex(16))
	h.Set("x-net-work-type", "NONE")
	h.Set("x-os-version", "MacIntel")
	h.Set("x-platform-version", "1")
	h.Set("x-protocol-version", "301")
	h.Set("x-provider-name", "NONE")
	h.Set("x-sdk-version", "9.0.2")
	h.Set("x-action", "401")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &guangyaHTTPError{Status: resp.StatusCode, Body: string(respBody)}
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return fmt.Errorf("刷新响应解析失败: %v", err)
	}
	accessToken := guangyaClean(guangyaDeepGet(payload, "access_token"))
	if accessToken == "" {
		return fmt.Errorf("刷新响应缺少 access_token: %s", guangyaExtractMessage(payload))
	}
	newCred := guangyaCredential{
		accessToken:  accessToken,
		refreshToken: guangyaClean(payload["refresh_token"]),
		deviceID:     cred.deviceID,
	}
	if newCred.refreshToken == "" {
		newCred.refreshToken = cred.refreshToken
	}
	expiresAt := time.Time{}
	if sec := guangyaToInt(payload["expires_in"]); sec > 0 {
		expiresAt = time.Now().Add(time.Duration(sec) * time.Second)
	}

	c.mu.Lock()
	c.cred = newCred
	if sec := guangyaToInt(payload["expires_in"]); sec > 0 {
		c.expiresAt = &expiresAt
	} else {
		c.expiresAt = nil
	}
	c.mu.Unlock()

	// 回写账号凭据（单一真相，research D2；落库失败仅记日志不阻断）
	if c.cksRepo != nil && c.account != nil && c.account.ID > 0 {
		if err := c.cksRepo.UpdateCk(c.account.ID, newCred.serialize()); err != nil {
			utils.Error("[光鸭:CLIENT] 令牌刷新回写失败（内存已更新） accountID=%d err=%v", c.account.ID, err)
		} else {
			c.account.Ck = newCred.serialize()
			utils.Info("[光鸭:CLIENT] 令牌已刷新并回写 accountID=%d", c.account.ID)
		}
	}
	return nil
}

// --- 登录态 API ---

// fsFiles 个人盘目录列表（分页从 0 起，contracts #6）
func (c *guangyaClient) fsFiles(parentID string, page, pageSize int) (map[string]interface{}, error) {
	if parentID == "0" || parentID == "/" || parentID == "root" {
		parentID = ""
	}
	return c.doJSON(guangyaAPIUserResV2+"/file/get_file_list", map[string]interface{}{
		"parentId": parentID, "page": page, "pageSize": pageSize, "orderBy": 0, "sortType": 0,
	})
}

// fsCreateDir 建目录（contracts #7）
func (c *guangyaClient) fsCreateDir(dirName, parentID string) (map[string]interface{}, error) {
	if parentID == "0" || parentID == "/" || parentID == "root" {
		parentID = ""
	}
	return c.doJSON(guangyaAPIUserRes+"/file/create_dir", map[string]interface{}{
		"dirName": dirName, "parentId": parentID,
	})
}

// fsDelete 删文件（contracts #8）
func (c *guangyaClient) fsDelete(fileIDs []string) (map[string]interface{}, error) {
	return c.doJSON(guangyaAPIUserRes+"/file/delete_file", map[string]interface{}{"fileIds": fileIDs})
}

// restoreShare 转存分享到自己盘（contracts #4；不返回新 fid，靠目录 diff）
func (c *guangyaClient) restoreShare(accessToken string, fileIDs []string, parentID string) (map[string]interface{}, error) {
	if parentID == "0" || parentID == "/" || parentID == "root" {
		parentID = ""
	}
	return c.doJSON(guangyaAPIUserRes+"/restore_share", map[string]interface{}{
		"accessToken": accessToken, "fileIds": fileIDs, "parentId": parentID,
	})
}

// taskStatus 任务状态查询（contracts #5）
func (c *guangyaClient) taskStatus(taskID string) (map[string]interface{}, error) {
	return c.doJSON(guangyaAPIUserRes+"/get_task_status", map[string]interface{}{"taskId": taskID})
}

// shareFile 创建分享（contracts #9，D9 载荷原样复制，禁止改字段）
func (c *guangyaClient) shareFile(fileIDs []string, title string) (map[string]interface{}, error) {
	return c.doJSON(guangyaAPIUserResV2+"/share_file", buildGuangyaSharePayload(fileIDs, title))
}

// buildGuangyaSharePayload 创建分享请求载荷（抓包原样，D9；抽出便于单测逐字段校验）
func buildGuangyaSharePayload(fileIDs []string, title string) map[string]interface{} {
	return map[string]interface{}{
		"fileIds":          fileIDs,
		"title":            title,
		"validateDuration": 0,
		"shareType":        1,
		"autoFillCode":     true,
		"trafficLimit":     "0",
		"maxRestoreCount":  0,
		"downloadType":     1,
		"enableShareCode":  false,
		"shareCode":        "",
	}
}

// getAssets 容量/会员（contracts #10）
func (c *guangyaClient) getAssets() (map[string]interface{}, error) {
	return c.doJSON(guangyaAPIAssets+"/get_assets", map[string]interface{}{})
}

// accountUserInfo user/me 用户信息（contracts #12；不稳定，仅作昵称来源）
func (c *guangyaClient) accountUserInfo() (map[string]interface{}, error) {
	return c.doJSON(guangyaAccountAPI+"/user/me", map[string]interface{}{})
}

// --- 匿名 API（无需登录态，公开分享访问；每次请求用随机 did）---

func guangyaPublicPost(endpoint string, payload interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	guangyaThrottle()
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	setGuangyaRequestHeaders(req, generateGuangyaDid(), "")
	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &guangyaHTTPError{Status: resp.StatusCode, Body: string(respBody)}
	}
	var out map[string]interface{}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("响应解析失败: %v body=%s", err, headSnippet(respBody))
	}
	return out, nil
}

// guangyaShareSummary 分享摘要（匿名，存在性校验，contracts #1）
func guangyaShareSummary(shareID string) (map[string]interface{}, error) {
	return guangyaPublicPost(guangyaAPIUserRes+"/get_share_summary", map[string]interface{}{"shareId": shareID})
}

// guangyaShareAccessToken 提取码换分享访问令牌（匿名，contracts #2）
func guangyaShareAccessToken(shareID, code string) (map[string]interface{}, error) {
	return guangyaPublicPost(guangyaAPIUserRes+"/get_share_access_token", map[string]interface{}{
		"shareId": shareID, "code": code,
	})
}

// guangyaShareFilesList 分享文件列表（匿名，分页从 1 起，contracts #3）
func guangyaShareFilesList(accessToken, parentID string, page, pageSize int) (map[string]interface{}, error) {
	return guangyaPublicPost(guangyaAPIUserRes+"/get_share_page_files_list", map[string]interface{}{
		"accessToken": accessToken, "parentId": parentID, "page": page, "pageSize": pageSize, "orderBy": 0, "sortType": 0,
	})
}

// ============================================================================
// URL 解析（data-model.md URL 契约）
// ============================================================================

// isGuangyaURL 判断是否光鸭链接（guangyapan.com 子串匹配，含 www./app. 变体）
func isGuangyaURL(rawURL string) bool {
	return strings.Contains(strings.ToLower(strings.TrimSpace(rawURL)), "guangyapan.com")
}

// IsGuangyaURL 判断是否光鸭链接（导出供 link_check_service 分流，internal.md §2）
func IsGuangyaURL(rawURL string) bool {
	return isGuangyaURL(rawURL)
}

// parseGuangyaShareURL 从光鸭分享链接解析分享标识与提取码。
// 提取码优先级：URL 查询参数(code/pwd/passcode/accessCode) > 调用方补充（config.Code）。
// 注意：分享标识大小写敏感，禁止对 URL 整体 ToLower。
func parseGuangyaShareURL(rawURL string) (shareID, passcode string) {
	normalized := strings.TrimSpace(rawURL)
	normalized = strings.ReplaceAll(normalized, "？", "?")
	normalized = strings.ReplaceAll(normalized, "＆", "&")
	normalized = strings.Join(strings.Fields(normalized), "")

	parsed, err := url.Parse(normalized)
	if err != nil {
		return "", ""
	}
	query := parsed.Query()
	// 查询键大小写不敏感匹配（精确名优先，再按 EqualFold 兜底）
	getCI := func(keys ...string) string {
		for _, key := range keys {
			if v := strings.TrimSpace(query.Get(key)); v != "" {
				return v
			}
		}
		for k, vs := range query {
			for _, key := range keys {
				if strings.EqualFold(k, key) && len(vs) > 0 {
					if v := strings.TrimSpace(vs[0]); v != "" {
						return v
					}
				}
			}
		}
		return ""
	}

	passcode = getCI("pwd", "code", "passcode", "accessCode")
	shareID = getCI("shareId", "share_id", "id", "sid")

	// 路径模式 /s/ /share/ /link/ /download/（保留原始大小写）
	if shareID == "" && parsed.Path != "" {
		path := parsed.Path
		for _, prefix := range []string{"/s/", "/share/", "/link/", "/download/"} {
			if idx := strings.Index(path, prefix); idx != -1 {
				token := strings.Trim(path[idx+len(prefix):], "/")
				if token != "" {
					shareID = token
					break
				}
			}
		}
	}
	return shareID, passcode
}

// extractGuangyaPasscodeFromText 定义于文件末尾（guangyaPasscodePatterns 一处维护）

// ============================================================================
// GuangyaService：PanService + Sharer 实现（镜像 uc_pan.go 模式）
// ============================================================================

// GuangyaService 光鸭云盘服务
type GuangyaService struct {
	*BasePanService
	configMutex sync.RWMutex
	client      *guangyaClient
	cksRepo     repo.CksRepository
	account     entity.Cks
}

// NewGuangyaService 创建光鸭网盘服务
func NewGuangyaService(config *PanConfig) *GuangyaService {
	service := &GuangyaService{
		BasePanService: NewBasePanService(config),
	}
	cookie := ""
	if config != nil {
		cookie = config.Cookie
	}
	if cred, err := parseGuangyaCredential(cookie); err == nil {
		service.client = newGuangyaClient(cred)
	} else {
		utils.Debug("[光鸭] 创建服务（凭据为空或无效，等待 GetUserInfo/UpdateConfig 注入）: %v", err)
	}
	utils.Debug("[光鸭] 创建光鸭服务 cookieLen=%d", len(cookie))
	return service
}

// GetServiceType 获取服务类型
func (g *GuangyaService) GetServiceType() ServiceType {
	return Guangya
}

// UpdateConfig 更新配置（凭据随之重建）
func (g *GuangyaService) UpdateConfig(config *PanConfig) {
	if config == nil {
		return
	}
	g.configMutex.Lock()
	defer g.configMutex.Unlock()
	g.config = config
	if config.Cookie != "" {
		if cred, err := parseGuangyaCredential(config.Cookie); err == nil {
			if g.client != nil {
				g.client.mu.Lock()
				g.client.cred = cred
				g.client.mu.Unlock()
			} else {
				g.client = newGuangyaClient(cred)
			}
		} else {
			utils.Debug("[光鸭] UpdateConfig 凭据解析失败（保留旧凭据）: %v", err)
		}
	}
}

// SetCKSRepository 设置账号仓储（令牌刷新回写钩子，research D2）
func (g *GuangyaService) SetCKSRepository(cksRepo repo.CksRepository, account entity.Cks) {
	g.cksRepo = cksRepo
	g.account = account
	if g.client != nil {
		g.client.cksRepo = cksRepo
		g.client.account = &account
	}
}

// clientFor 按传入凭据取可用 client（GetUserInfo 传 ck；Transfer 用服务内凭据）
func (g *GuangyaService) clientFor(ck *string) (*guangyaClient, error) {
	if ck != nil && *ck != "" {
		cred, err := parseGuangyaCredential(*ck)
		if err != nil {
			return nil, err
		}
		c := newGuangyaClient(cred)
		c.cksRepo = g.cksRepo
		if g.account.ID > 0 {
			c.account = &g.account
		}
		return c, nil
	}
	g.configMutex.RLock()
	client := g.client
	g.configMutex.RUnlock()
	if client == nil {
		return nil, fmt.Errorf("光鸭凭据未配置")
	}
	return client, nil
}

// ============================================================================
// GetUserInfo：账号校验 + 容量/会员（US1；get_assets 主路径 + 文件列表探活回退）
// ============================================================================

func (g *GuangyaService) GetUserInfo(ck *string) (*UserInfo, error) {
	client, err := g.clientFor(ck)
	if err != nil {
		return nil, err
	}
	utils.Debug("[光鸭] GetUserInfo 请求 get_assets")

	assetsPayload, assetsErr := client.getAssets()
	assetsData := guangyaAsMap(guangyaExtractData(assetsPayload))

	if assetsErr != nil || assetsData == nil || !guangyaIsSuccess(assetsPayload) {
		// 回退：个人盘列表探活（user/assets 不稳时的轻量可用性验证，research D2）
		utils.Debug("[光鸭] get_assets 失败，回退文件列表探活: assetsErr=%v success=%v", assetsErr, assetsPayload != nil && guangyaIsSuccess(assetsPayload))
		probePayload, probeErr := client.fsFiles("", 0, 1)
		if probeErr != nil || !guangyaIsSuccess(probePayload) {
			if assetsErr != nil {
				return nil, fmt.Errorf("光鸭账号校验失败: %v", assetsErr)
			}
			return nil, fmt.Errorf("光鸭账号校验失败: %s", guangyaExtractMessage(assetsPayload))
		}
		if probeErr == nil && assetsErr == nil && !guangyaIsSuccess(assetsPayload) {
			return nil, fmt.Errorf("光鸭账号校验失败: %s", guangyaExtractMessage(assetsPayload))
		}
		assetsData = map[string]interface{}{} // 探活通过但容量未知
		utils.Warn("[光鸭] get_assets 不可用，探活通过但容量/会员未知（显示 0）")
	}

	// 昵称三级回退：user/me → 手机号掩码 → 占位（internal.md）
	nickname := ""
	if infoPayload, infoErr := client.accountUserInfo(); infoErr == nil && guangyaIsSuccess(infoPayload) {
		info := guangyaAsMap(guangyaExtractData(infoPayload))
		if info == nil {
			info = infoPayload
		}
		nickname = guangyaClean(guangyaDeepGet(info, "nickname", "nickName", "username"))
		if nickname == "" {
			if phone := guangyaClean(guangyaDeepGet(info, "phoneNumber", "phone", "mobile")); phone != "" {
				nickname = maskGuangyaMobile(phone)
			}
		}
	}
	if nickname == "" {
		nickname = "光鸭云盘用户"
	}

	totalSpace := guangyaToInt(assetsData["totalSpaceSize"])
	usedSpace := guangyaToInt(assetsData["usedSpaceSize"])
	vipStatus := guangyaToInt(assetsData["vipStatus"]) > 0 || guangyaToInt(assetsData["svipStatus"]) > 0

	// 规范化凭据回传（AddCks 据此把规范化 kv 写入 Ck，research D2 单一真相）
	client.mu.Lock()
	normalized := client.cred.serialize()
	client.mu.Unlock()

	utils.Debug("[光鸭] GetUserInfo 成功 nickname=%s total=%d used=%d vip=%v", nickname, totalSpace, usedSpace, vipStatus)
	return &UserInfo{
		Username:    nickname,
		VIPStatus:   vipStatus,
		UsedSpace:   usedSpace,
		TotalSpace:  totalSpace,
		ServiceType: "guangya",
		ExtraData:   normalized,
	}, nil
}

// ============================================================================
// Transfer：转存主链路（US2；摘要→令牌→列表→转存→轮询→diff→广告→再分享）
// ============================================================================

// Transfer 转存分享链接（流程见 data-model.md 转存链路数据流）
func (g *GuangyaService) Transfer(shareID string) (*TransferResult, error) {
	g.configMutex.RLock()
	config := g.config
	client := g.client
	g.configMutex.RUnlock()

	isType := 0
	rawURL := ""
	passcode := ""
	if config != nil {
		isType = config.IsType
		rawURL = config.URL
		passcode = config.Code
	}
	if passcode == "" && rawURL != "" {
		_, passcode = parseGuangyaShareURL(rawURL)
	}
	utils.Info("[光鸭] 开始处理分享 shareID=%s isType=%d(0=转存,1=校验) passcodeLen=%d", shareID, isType, len(passcode))

	// 1) 分享摘要校验（匿名；完整标识失败时数字段兜底重试一次，research D8）
	if err := g.ensureShareExists(shareID); err != nil {
		utils.Error("[光鸭] 分享校验失败 shareID=%s err=%v", shareID, err)
		return ErrorResult(fmt.Sprintf("%v", err)), nil
	}

	// 2) 提取码换访问令牌（匿名）
	accessToken, err := g.fetchShareAccessToken(shareID, passcode)
	if err != nil {
		utils.Error("[光鸭] 获取分享访问令牌失败 shareID=%s err=%v", shareID, err)
		return ErrorResult(fmt.Sprintf("%v", err)), nil
	}

	// 3) 分享文件列表（匿名，分页聚合）
	items, err := listGuangyaShareItems(accessToken, "")
	if err != nil {
		utils.Error("[光鸭] 获取分享文件列表失败 shareID=%s err=%v", shareID, err)
		return ErrorResult(fmt.Sprintf("获取分享文件列表失败: %v", err)), nil
	}
	if len(items) == 0 {
		return ErrorResult("分享文件列表为空"), nil
	}

	title := guangyaItemName(items[0])
	utils.Debug("[光鸭] 获取分享详情成功 shareID=%s title=%s fileCount=%d", shareID, title, len(items))

	// 4) 检验模式（isType==1，只读不转存）
	if isType == 1 {
		utils.Info("[光鸭] 校验模式完成（不转存）shareID=%s title=%s", shareID, title)
		return SuccessResult("检验成功", map[string]interface{}{
			"title":    title,
			"shareUrl": rawURL,
		}), nil
	}

	if client == nil {
		return ErrorResult("光鸭凭据未配置"), nil
	}

	// 5) 确保专用目录「转存」存在（D4）
	dirFid, err := g.ensureTransferDir(client)
	if err != nil {
		utils.Error("[光鸭] 确保转存目录失败 err=%v", err)
		return ErrorResult(fmt.Sprintf("创建转存目录失败: %v", err)), nil
	}

	// 6) 转存前快照专用目录
	beforeItems, err := listGuangyaFsItems(client, dirFid, 0)
	if err != nil {
		utils.Error("[光鸭] 转存前目录快照失败 err=%v", err)
		return ErrorResult(fmt.Sprintf("获取转存目录失败: %v", err)), nil
	}
	beforeFids := make(map[string]bool, len(beforeItems))
	fileNames := make([]string, 0, len(items))
	shareFids := make([]string, 0, len(items))
	for _, it := range beforeItems {
		if fid := guangyaItemFid(it); fid != "" {
			beforeFids[fid] = true
		}
	}
	for _, it := range items {
		if fid := guangyaItemFid(it); fid != "" {
			shareFids = append(shareFids, fid)
			fileNames = append(fileNames, guangyaItemName(it))
		}
	}

	// 7) 转存 + 轮询
	restorePayload, err := client.restoreShare(accessToken, shareFids, dirFid)
	if err != nil {
		utils.Error("[光鸭] 转存请求失败 shareID=%s err=%v", shareID, err)
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}
	if err := guangyaRaiseIfFailed(restorePayload, "转存失败"); err != nil {
		utils.Error("[光鸭] 转存失败 shareID=%s err=%v", shareID, err)
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}
	if taskID := guangyaExtractTaskID(restorePayload); taskID != "" {
		if err := waitGuangyaTask(client, taskID); err != nil {
			utils.Error("[光鸭] 转存任务未完成 shareID=%s taskID=%s err=%v", shareID, taskID, err)
			return ErrorResult(fmt.Sprintf("%v", err)), nil
		}
	}

	// 8) 目录 diff 按文件名对齐新 fid（重试等待目录可见）
	newFids, err := alignGuangyaNewFids(client, dirFid, beforeFids, fileNames)
	if err != nil {
		utils.Error("[光鸭] 转存后文件对齐失败 shareID=%s err=%v", shareID, err)
		return ErrorResult(fmt.Sprintf("%v", err)), nil
	}
	utils.Debug("[光鸭] 转存完成对齐 shareID=%s newFids=%v", shareID, newFids)

	// 9) 广告处理（删来源广告 + 条件插入，research D6）
	newFids = g.cleanupAdFiles(client, dirFid, fileNames, newFids)
	if err := g.addAdIfNeeded(client, dirFid); err != nil {
		utils.Debug("[光鸭] 广告插入失败（不阻断）err=%v", err)
	}

	// 10) 创建分享（永久 + 自带提取码，D9 载荷原样）
	sharePayload, err := client.shareFile(newFids, title)
	if err != nil {
		utils.Error("[光鸭] 创建分享失败 err=%v", err)
		return ErrorResult(fmt.Sprintf("创建分享失败: %v", err)), nil
	}
	if err := guangyaRaiseIfFailed(sharePayload, "创建分享失败"); err != nil {
		utils.Error("[光鸭] 创建分享失败 err=%v", err)
		return ErrorResult(fmt.Sprintf("创建分享失败: %v", err)), nil
	}
	shareData := guangyaAsMap(guangyaExtractData(sharePayload))
	shareURL := guangyaClean(shareData["shareUrl"])
	if shareURL == "" {
		return ErrorResult("创建分享成功但未获取到分享链接"), nil
	}
	code := guangyaClean(shareData["code"])
	fid := strings.Join(newFids, ",")

	utils.Info("[光鸭] 转存成功 shareID=%s newShareUrl=%s title=%s fid=%s", shareID, shareURL, title, fid)
	return SuccessResult("转存成功", map[string]interface{}{
		"shareUrl": shareURL,
		"title":    title,
		"fid":      fid,
		"code":     code,
	}), nil
}

// Share 对系统已存文件按 fid 重新生成光鸭分享（Sharer，FR-015/US4）。
// fid 支持逗号分隔多值（Transfer 产出形态）。
func (g *GuangyaService) Share(fid string) (*TransferResult, error) {
	if strings.TrimSpace(fid) == "" {
		return &TransferResult{Success: false, Message: "fid 为空"}, nil
	}
	g.configMutex.RLock()
	client := g.client
	g.configMutex.RUnlock()
	if client == nil {
		return &TransferResult{Success: false, Message: "光鸭凭据未配置"}, nil
	}
	fids := make([]string, 0, 4)
	for _, part := range strings.Split(fid, ",") {
		if p := strings.TrimSpace(part); p != "" {
			fids = append(fids, p)
		}
	}
	sharePayload, err := client.shareFile(fids, "资源分享")
	if err != nil {
		return &TransferResult{Success: false, Message: fmt.Sprintf("创建分享失败: %v", err)}, nil
	}
	if err := guangyaRaiseIfFailed(sharePayload, "创建分享失败"); err != nil {
		return &TransferResult{Success: false, Message: fmt.Sprintf("创建分享失败: %v", err)}, nil
	}
	shareData := guangyaAsMap(guangyaExtractData(sharePayload))
	shareURL := guangyaClean(shareData["shareUrl"])
	if shareURL == "" {
		return &TransferResult{Success: false, Message: "创建分享成功但未获取到分享链接"}, nil
	}
	utils.Info("[光鸭:SHARE] 重新分享成功 fid=%s url=%s", fid, shareURL)
	return &TransferResult{Success: true, ShareURL: shareURL, Fid: fid}, nil
}

// ============================================================================
// GetFiles / DeleteFiles（供广告清理与自动清理调度）
// ============================================================================

// GetFiles 获取目录文件列表
func (g *GuangyaService) GetFiles(pdirFid string) (*TransferResult, error) {
	if pdirFid == "" {
		pdirFid = "0"
	}
	client, err := g.clientFor(nil)
	if err != nil {
		return ErrorResult(fmt.Sprintf("%v", err)), nil
	}
	items, err := listGuangyaFsItems(client, pdirFid, 0)
	if err != nil {
		utils.Error("[光鸭] GetFiles 请求失败 pdirFid=%s err=%v", pdirFid, err)
		return ErrorResult(fmt.Sprintf("获取光鸭文件列表失败: %v", err)), nil
	}
	list := make([]interface{}, 0, len(items))
	for _, it := range items {
		list = append(list, it)
	}
	utils.Debug("[光鸭] GetFiles 成功 pdirFid=%s count=%d", pdirFid, len(items))
	return SuccessResult("获取成功", list), nil
}

// DeleteFiles 批量删除文件（元素支持逗号分隔多 fid；空列表幂等成功）
func (g *GuangyaService) DeleteFiles(fileList []string) (*TransferResult, error) {
	if len(fileList) == 0 {
		return SuccessResult("删除成功", nil), nil
	}
	client, err := g.clientFor(nil)
	if err != nil {
		return ErrorResult(fmt.Sprintf("%v", err)), nil
	}
	fids := make([]string, 0, len(fileList))
	for _, entry := range fileList {
		for _, part := range strings.Split(entry, ",") {
			if p := strings.TrimSpace(part); p != "" {
				fids = append(fids, p)
			}
		}
	}
	if len(fids) == 0 {
		return SuccessResult("删除成功", nil), nil
	}
	payload, err := client.fsDelete(fids)
	if err != nil {
		utils.Error("[光鸭] 删除文件失败 err=%v", err)
		return ErrorResult(fmt.Sprintf("删除光鸭文件失败: %v", err)), nil
	}
	if err := guangyaRaiseIfFailed(payload, "删除失败"); err != nil {
		utils.Error("[光鸭] 删除文件失败 err=%v", err)
		return ErrorResult(fmt.Sprintf("删除光鸭文件失败: %v", err)), nil
	}
	utils.Debug("[光鸭] DeleteFiles 成功 count=%d", len(fids))
	return SuccessResult("删除成功", nil), nil
}

// ============================================================================
// 匿名校验器（US3；结论映射见 data-model.md）
// ============================================================================

// CheckGuangyaLink 匿名校验光鸭分享链接（无需账号）。
// 返回 status: valid / invalid / undetermined 与可读原因。
func CheckGuangyaLink(rawURL string) (status, reason string) {
	shareID, passcode := parseGuangyaShareURL(rawURL)
	if shareID == "" {
		return "undetermined", "无法解析分享标识"
	}

	// 分享摘要：不存在/已失效 → invalid
	summaryOK := false
	summaryPayload, err := guangyaShareSummary(shareID)
	if err == nil && guangyaIsSuccess(summaryPayload) {
		summaryOK = true
	}
	if !summaryOK {
		// 数字段兜底重试一次（research D8）
		if idx := strings.Index(shareID, "_"); idx > 0 {
			numeric := shareID[:idx]
			retryPayload, retryErr := guangyaShareSummary(numeric)
			if retryErr == nil && guangyaIsSuccess(retryPayload) {
				summaryOK = true
				shareID = numeric
			}
		}
	}
	if !summaryOK {
		if err != nil {
			return "undetermined", "" // 网络错误不翻转状态
		}
		return "invalid", "分享已取消"
	}

	// 提取码换访问令牌：失败 → invalid（区分提取码错误）
	tokenPayload, err := guangyaShareAccessToken(shareID, passcode)
	if err != nil {
		return "undetermined", ""
	}
	if extractGuangyaShareAccessToken(tokenPayload) == "" {
		if passcode != "" {
			return "invalid", "提取码错误"
		}
		return "invalid", "分享不可访问（需要提取码）"
	}
	return "valid", ""
}

// ============================================================================
// 内部工具
// ============================================================================

// ensureShareExists 分享摘要校验（完整标识优先，数字段兜底，research D8）
func (g *GuangyaService) ensureShareExists(shareID string) error {
	payload, err := guangyaShareSummary(shareID)
	if err == nil && guangyaIsSuccess(payload) {
		return nil
	}
	if idx := strings.Index(shareID, "_"); idx > 0 {
		numeric := shareID[:idx]
		retryPayload, retryErr := guangyaShareSummary(numeric)
		if retryErr == nil && guangyaIsSuccess(retryPayload) {
			return nil
		}
	}
	if err != nil {
		return fmt.Errorf("分享校验请求失败: %v", err)
	}
	return fmt.Errorf("分享不存在或已失效")
}

// fetchShareAccessToken 提取码换分享访问令牌
func (g *GuangyaService) fetchShareAccessToken(shareID, passcode string) (string, error) {
	payload, err := guangyaShareAccessToken(shareID, passcode)
	if err != nil {
		return "", fmt.Errorf("获取分享访问令牌失败: %v", err)
	}
	token := extractGuangyaShareAccessToken(payload)
	if token == "" {
		if passcode != "" {
			return "", fmt.Errorf("提取码错误")
		}
		return "", fmt.Errorf("分享不可访问（需要提取码）")
	}
	return token, nil
}

// extractGuangyaShareAccessToken 从响应提取访问令牌（data 为串或对象）
func extractGuangyaShareAccessToken(payload map[string]interface{}) string {
	data := guangyaExtractData(payload)
	if s, ok := data.(string); ok {
		return strings.TrimSpace(s)
	}
	if m := guangyaAsMap(data); m != nil {
		if v := guangyaClean(guangyaDeepGet(m, "access_token", "accessToken", "token")); v != "" {
			return v
		}
	}
	return guangyaClean(guangyaDeepGet(payload, "access_token", "accessToken", "token"))
}

// ensureTransferDir 确保专用目录「转存」存在，返回目录 fid（D4）
func (g *GuangyaService) ensureTransferDir(client *guangyaClient) (string, error) {
	items, err := listGuangyaFsItems(client, "", 0)
	if err != nil {
		return "", err
	}
	for _, it := range items {
		if guangyaItemIsDir(it) && guangyaItemName(it) == guangyaTransferDirName {
			fid := guangyaItemFid(it)
			if fid == "" {
				fid = "0"
			}
			return fid, nil
		}
	}
	payload, err := client.fsCreateDir(guangyaTransferDirName, "")
	if err != nil {
		return "", err
	}
	if err := guangyaRaiseIfFailed(payload, "创建目录失败"); err != nil {
		return "", err
	}
	newID := guangyaExtractTaskID(payload)
	if newID == "" {
		// 响应无 fid → 回查根目录按名匹配（contracts #7）
		items, err := listGuangyaFsItems(client, "", 0)
		if err != nil {
			return "", err
		}
		for _, it := range items {
			if guangyaItemIsDir(it) && guangyaItemName(it) == guangyaTransferDirName {
				newID = guangyaItemFid(it)
				break
			}
		}
	}
	if newID == "" {
		return "", fmt.Errorf("创建目录后未找到目标目录")
	}
	return newID, nil
}

// listGuangyaFsItems 聚合个人盘目录条目（分页从 0 起；maxItems=0 不限，上限 200 页）
func listGuangyaFsItems(client *guangyaClient, parentID string, maxItems int) ([]map[string]interface{}, error) {
	var all []map[string]interface{}
	const pageSize = 50
	for page := 0; page < guangyaPageMax; page++ {
		payload, err := client.fsFiles(parentID, page, pageSize)
		if err != nil {
			return nil, err
		}
		if !guangyaIsSuccess(payload) {
			return nil, fmt.Errorf("%s", guangyaExtractMessage(payload))
		}
		batch := guangyaExtractItems(payload)
		if len(batch) == 0 {
			break
		}
		all = append(all, batch...)
		if maxItems > 0 && len(all) >= maxItems {
			return all[:maxItems], nil
		}
		if len(batch) < pageSize {
			break
		}
	}
	return all, nil
}

// listGuangyaShareItems 聚合分享文件条目（匿名，分页从 1 起）
func listGuangyaShareItems(accessToken, parentID string) ([]map[string]interface{}, error) {
	var all []map[string]interface{}
	const pageSize = 50
	for page := 1; page <= guangyaPageMax; page++ {
		payload, err := guangyaShareFilesList(accessToken, parentID, page, pageSize)
		if err != nil {
			return nil, err
		}
		if !guangyaIsSuccess(payload) {
			return nil, fmt.Errorf("%s", guangyaExtractMessage(payload))
		}
		batch := guangyaExtractItems(payload)
		if len(batch) == 0 {
			break
		}
		all = append(all, batch...)
		if len(batch) < pageSize {
			break
		}
	}
	return all, nil
}

// guangyaItemFid 条目 fid 提取（多键名容错；空归一为 "0"）
func guangyaItemFid(item map[string]interface{}) string {
	v := guangyaDeepGet(item, "fid", "fileId", "id", "resId")
	if v == nil {
		dm := guangyaAsMap(guangyaExtractData(item))
		if dm != nil {
			v = guangyaDeepGet(dm, "fid", "fileId", "id", "resId")
		}
	}
	fid := guangyaClean(v)
	if fid == "" {
		return "0"
	}
	return fid
}

// guangyaItemName 条目名称提取
func guangyaItemName(item map[string]interface{}) string {
	v := guangyaDeepGet(item, "file_name", "fileName", "name", "title")
	if v == nil {
		dm := guangyaAsMap(guangyaExtractData(item))
		if dm != nil {
			v = guangyaDeepGet(dm, "file_name", "fileName", "name", "title")
		}
	}
	return guangyaClean(v)
}

// guangyaItemIsDir 条目目录标记提取
func guangyaItemIsDir(item map[string]interface{}) bool {
	for _, scope := range []map[string]interface{}{item, guangyaAsMap(guangyaExtractData(item))} {
		if scope == nil {
			continue
		}
		if v := guangyaDeepGet(scope, "dir", "isDir", "isFolder", "folder"); v != nil {
			return guangyaToBool(v)
		}
	}
	for _, scope := range []map[string]interface{}{item, guangyaAsMap(guangyaExtractData(item))} {
		if scope == nil {
			continue
		}
		if v := guangyaDeepGet(scope, "fileType", "type"); v != nil {
			switch strings.ToLower(guangyaClean(v)) {
			case "0", "dir", "folder", "directory":
				return true
			}
		}
		if v := guangyaDeepGet(scope, "resType"); v != nil {
			return strings.ToLower(guangyaClean(v)) == "2"
		}
	}
	return false
}

// guangyaExtractTaskID 从响应提取任务/对象 ID
func guangyaExtractTaskID(payload map[string]interface{}) string {
	for _, scope := range []map[string]interface{}{guangyaAsMap(guangyaExtractData(payload)), payload} {
		if scope == nil {
			continue
		}
		if v := guangyaClean(guangyaDeepGet(scope, "taskId", "task_id", "id")); v != "" {
			return v
		}
	}
	return ""
}

// waitGuangyaTask 轮询任务至完成（1s×20，D5）
func waitGuangyaTask(client *guangyaClient, taskID string) error {
	for i := 0; i < guangyaTaskPollMax; i++ {
		payload, err := client.taskStatus(taskID)
		if err != nil {
			return fmt.Errorf("转存任务查询失败: %v", err)
		}
		if guangyaTaskFailed(payload) {
			msg := guangyaExtractMessage(payload)
			if msg == "" {
				msg = "转存任务失败"
			}
			return fmt.Errorf("%s", msg)
		}
		if guangyaTaskDone(payload) {
			return nil
		}
		time.Sleep(guangyaTaskPollInterval)
	}
	return fmt.Errorf("转存任务超时")
}

// guangyaTaskDone / guangyaTaskFailed 任务状态判定（contracts #5）。
// 优先按显式 status 字段判定（status=1 的进行中任务 msg 也可能是 "success"，
// 消息回退仅在响应无 status 字段时生效）。
func guangyaTaskDone(payload map[string]interface{}) bool {
	status := strings.ToLower(guangyaClean(firstGuangyaTaskStatus(payload)))
	if status != "" {
		switch status {
		case "2", "3", "4", "done", "success", "completed", "finish", "finished":
			return true
		}
		return false
	}
	msg := strings.ToLower(guangyaExtractMessage(payload))
	return strings.Contains(msg, "完成") || strings.Contains(msg, "成功") ||
		strings.Contains(msg, "finished") || strings.Contains(msg, "success")
}

func guangyaTaskFailed(payload map[string]interface{}) bool {
	status := strings.ToLower(guangyaClean(firstGuangyaTaskStatus(payload)))
	if status != "" {
		switch status {
		case "5", "-1", "failed", "error":
			return true
		}
		return false
	}
	msg := strings.ToLower(guangyaExtractMessage(payload))
	return strings.Contains(msg, "失败") || strings.Contains(msg, "error") || strings.Contains(msg, "failed")
}

func firstGuangyaTaskStatus(payload map[string]interface{}) interface{} {
	for _, scope := range []map[string]interface{}{guangyaAsMap(guangyaExtractData(payload)), payload} {
		if scope == nil {
			continue
		}
		if v := guangyaDeepGet(scope, "status", "taskStatus", "state"); v != nil {
			return v
		}
	}
	return nil
}

// alignGuangyaNewFids 转存后目录 diff 按文件名对齐新 fid（多文件全对齐；失败可读错误）
func alignGuangyaNewFids(client *guangyaClient, dirFid string, beforeFids map[string]bool, fileNames []string) ([]string, error) {
	var aligned []string
	for attempt := 0; attempt < guangyaDiffRetryMax; attempt++ {
		items, err := listGuangyaFsItems(client, dirFid, 1000)
		if err != nil {
			return nil, fmt.Errorf("转存后获取目录失败: %v", err)
		}
		aligned = alignGuangyaFidsFromItems(items, beforeFids, fileNames)
		if len(aligned) > 0 && allNonEmpty(aligned) {
			return aligned, nil
		}
		time.Sleep(time.Second)
	}
	if len(aligned) == 0 || !anyNonEmpty(aligned) {
		return nil, fmt.Errorf("转存后未找到新文件")
	}
	// 部分对齐：返回已对齐项（同名文件干扰等场景，尽力交付）
	utils.Warn("[光鸭] 文件对齐部分完成 want=%d got=%d", len(fileNames), countNonEmpty(aligned))
	return filterEmpty(aligned), nil
}

// alignGuangyaFidsFromItems 纯函数：从目录条目中剔除转存前已有 fid，按文件名顺序对齐新 fid
// （转存接口不返回新 fid，D4；同名多条按未使用优先，每 fid 只用一次）
func alignGuangyaFidsFromItems(items []map[string]interface{}, beforeFids map[string]bool, fileNames []string) []string {
	nameToFids := make(map[string][]string)
	for _, it := range items {
		fid := guangyaItemFid(it)
		name := guangyaItemName(it)
		if fid == "" || fid == "0" || beforeFids[fid] || name == "" {
			continue
		}
		nameToFids[name] = append(nameToFids[name], fid)
	}
	used := make(map[string]bool)
	aligned := make([]string, 0, len(fileNames))
	for _, name := range fileNames {
		target := ""
		for _, fid := range nameToFids[name] {
			if !used[fid] {
				target = fid
				break
			}
		}
		if target != "" {
			used[target] = true
		}
		aligned = append(aligned, target)
	}
	return aligned
}

func allNonEmpty(list []string) bool {
	for _, v := range list {
		if v == "" {
			return false
		}
	}
	return len(list) > 0
}

func anyNonEmpty(list []string) bool { return countNonEmpty(list) > 0 }

func countNonEmpty(list []string) int {
	n := 0
	for _, v := range list {
		if v != "" {
			n++
		}
	}
	return n
}

func filterEmpty(list []string) []string {
	out := make([]string, 0, len(list))
	for _, v := range list {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

// cleanupAdFiles 删除转存目录内命中广告关键词的文件（复用 pan_ad.go 共享配置，D6）。
// 返回剔除广告后的 newFids（按 fileNames 位置对应）。
func (g *GuangyaService) cleanupAdFiles(client *guangyaClient, dirFid string, fileNames, newFids []string) []string {
	deleted := 0
	kept := make([]string, 0, len(newFids))
	for i, fid := range newFids {
		name := ""
		if i < len(fileNames) {
			name = fileNames[i]
		}
		if name != "" && containsAdKeywords(name) && fid != "" {
			utils.Debug("[光鸭] 命中广告关键词，删除 fileName=%s fid=%s", name, fid)
			if _, err := g.DeleteFiles([]string{fid}); err == nil {
				deleted++
				continue
			}
			utils.Debug("[光鸭] 删除广告文件失败（继续）fileName=%s", name)
		}
		kept = append(kept, fid)
	}
	if deleted > 0 {
		utils.Debug("[光鸭] 广告清理完成 dirFid=%s deleted=%d", dirFid, deleted)
	}
	return kept
}

// addAdIfNeeded 插入自定义广告：仅当广告分享链接为光鸭平台时转存进产出目录，否则跳过（D6）
func (g *GuangyaService) addAdIfNeeded(client *guangyaClient, dirFid string) error {
	autoInsertAdStr, err := getAdSystemConfigValue(entity.ConfigKeyAutoInsertAd)
	if err != nil {
		return err
	}
	if autoInsertAdStr == "" {
		return nil
	}
	adURLs := splitAdURLs(autoInsertAdStr)
	if len(adURLs) == 0 {
		return nil
	}

	// 选第一个光鸭平台的广告链接；没有则跳过（跨盘转存不可行）
	var guangyaAdURL string
	for _, u := range adURLs {
		if isGuangyaURL(u) {
			guangyaAdURL = u
			break
		}
	}
	if guangyaAdURL == "" {
		utils.Debug("[光鸭] 广告链接非光鸭平台，跳过插入（adURLCount=%d）", len(adURLs))
		return nil
	}

	shareID, passcode := parseGuangyaShareURL(guangyaAdURL)
	if shareID == "" {
		return fmt.Errorf("广告分享链接无法解析")
	}
	summaryPayload, err := guangyaShareSummary(shareID)
	if err != nil || !guangyaIsSuccess(summaryPayload) {
		return fmt.Errorf("广告分享不可用")
	}
	accessToken, err := g.fetchShareAccessToken(shareID, passcode)
	if err != nil {
		return err
	}
	adItems, err := listGuangyaShareItems(accessToken, "")
	if err != nil || len(adItems) == 0 {
		return fmt.Errorf("广告文件列表为空")
	}
	adFids := make([]string, 0, len(adItems))
	for _, it := range adItems {
		if fid := guangyaItemFid(it); fid != "" && fid != "0" {
			adFids = append(adFids, fid)
		}
	}
	if len(adFids) == 0 {
		return fmt.Errorf("广告文件为空")
	}
	restorePayload, err := client.restoreShare(accessToken, adFids, dirFid)
	if err != nil {
		return err
	}
	if err := guangyaRaiseIfFailed(restorePayload, "广告转存失败"); err != nil {
		return err
	}
	if taskID := guangyaExtractTaskID(restorePayload); taskID != "" {
		if err := waitGuangyaTask(client, taskID); err != nil {
			return err
		}
	}
	utils.Debug("[光鸭] 广告插入成功 adShareID=%s dirFid=%s", shareID, dirFid)
	return nil
}

// guangyaPasscodePatterns 提取码文本模式（含全角括号/访问码变体）
var guangyaPasscodePatterns = []string{
	`[（(]提取码[：:]\s*([a-zA-Z0-9]{4,8})[)）]`,
	`[（(]访问码[：:]\s*([a-zA-Z0-9]{4,8})[)）]`,
	`提取码[：:]\s*([a-zA-Z0-9]{4,8})`,
	`访问码[：:]\s*([a-zA-Z0-9]{4,8})`,
}

// extractGuangyaPasscodeFromText 从消息正文提取"提取码：xxxx"（含全角括号/访问码变体）
func extractGuangyaPasscodeFromText(text string) string {
	compact := strings.Join(strings.Fields(text), "")
	for _, pattern := range guangyaPasscodePatterns {
		if m := regexp.MustCompile(pattern).FindStringSubmatch(compact); len(m) > 1 {
			return m[1]
		}
	}
	return ""
}
