package pan

// 光鸭云盘单元测试（specs/013-guangya-pan-integration tasks T005/T012）
// 覆盖：凭据 kv 解析四规则/序列化往返/手机号掩码、URL 与提取码解析、
// 响应判定容错、创建分享载荷逐字段、diff 按文件名对齐。

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// ============================================================================
// 凭据解析（data-model.md 四条规则）
// ============================================================================

func TestParseGuangyaCredential_FullKV(t *testing.T) {
	c, err := parseGuangyaCredential("access_token=eyJhbGci; refresh_token=rt-xyz; device_id=abc123")
	if err != nil {
		t.Fatalf("完整 kv 解析失败: %v", err)
	}
	if c.accessToken != "eyJhbGci" || c.refreshToken != "rt-xyz" || c.deviceID != "abc123" {
		t.Fatalf("解析值不符: %+v", c)
	}
}

func TestParseGuangyaCredential_BareToken(t *testing.T) {
	c, err := parseGuangyaCredential("  eyJhbGciOiJIUzI1NiJ9.payload.sig ")
	if err != nil {
		t.Fatalf("裸 token 解析失败: %v", err)
	}
	if c.accessToken != "eyJhbGciOiJIUzI1NiJ9.payload.sig" {
		t.Fatalf("accessToken 不符: %s", c.accessToken)
	}
	if c.refreshToken != "" {
		t.Fatalf("裸 token 不应有 refresh_token: %s", c.refreshToken)
	}
	if c.deviceID == "" {
		t.Fatal("裸 token 应自动生成 device_id")
	}
}

func TestParseGuangyaCredential_AutoDeviceID(t *testing.T) {
	c, err := parseGuangyaCredential("access_token=tok;refresh_token=rt")
	if err != nil {
		t.Fatalf("缺 device_id 解析失败: %v", err)
	}
	if len(c.deviceID) != 32 {
		t.Fatalf("自动生成的 device_id 应为 32 位 hex（md5），实际 %d 位: %s", len(c.deviceID), c.deviceID)
	}
	// 序列化后 device_id 应保留（规范化回写）
	if !strings.Contains(c.serialize(), "device_id="+c.deviceID) {
		t.Fatalf("序列化应包含生成的 device_id: %s", c.serialize())
	}
}

func TestParseGuangyaCredential_Errors(t *testing.T) {
	if _, err := parseGuangyaCredential(""); err == nil {
		t.Fatal("空串应报错")
	}
	if _, err := parseGuangyaCredential("   "); err == nil {
		t.Fatal("空白串应报错")
	}
	if _, err := parseGuangyaCredential("refresh_token=rt;device_id=d1"); err == nil {
		t.Fatal("缺少 access_token 应报错")
	}
}

func TestGuangyaCredential_SerializeRoundTrip(t *testing.T) {
	orig := "access_token=tok1;refresh_token=rt1;device_id=dev1"
	c, err := parseGuangyaCredential(orig)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	again, err := parseGuangyaCredential(c.serialize())
	if err != nil {
		t.Fatalf("序列化往返解析失败: %v", err)
	}
	if again != c {
		t.Fatalf("往返不一致: %+v vs %+v", again, c)
	}
}

func TestMaskGuangyaMobile(t *testing.T) {
	if got := maskGuangyaMobile("+86 138 1234 5678"); got != "138****5678" {
		t.Fatalf("掩码不符: %s", got)
	}
	if got := maskGuangyaMobile("12345"); got != "12345" {
		t.Fatalf("短号应原样返回: %s", got)
	}
}

// ============================================================================
// URL / 提取码解析（data-model.md URL 契约）
// ============================================================================

func TestParseGuangyaShareURL_Standard(t *testing.T) {
	shareID, code := parseGuangyaShareURL("https://www.guangyapan.com/s/1938809643918233637_aoupdQ-VQOcjiwM3?code=yvy1")
	if shareID != "1938809643918233637_aoupdQ-VQOcjiwM3" {
		t.Fatalf("shareID 应为完整标识（含 _后缀），实际: %s", shareID)
	}
	if code != "yvy1" {
		t.Fatalf("code 不符: %s", code)
	}
}

func TestParseGuangyaShareURL_PwdParam(t *testing.T) {
	shareID, code := parseGuangyaShareURL("https://guangyapan.com/s/abc_def?pwd=wxyz")
	if shareID != "abc_def" {
		t.Fatalf("shareID 不符: %s", shareID)
	}
	if code != "wxyz" {
		t.Fatalf("pwd 提取不符: %s", code)
	}
}

func TestParseGuangyaShareURL_FullWidthAndDirty(t *testing.T) {
	shareID, code := parseGuangyaShareURL("https://www.guangyapan.com/s/tok_123？code=abcd")
	if shareID != "tok_123" {
		t.Fatalf("全角？应归一，shareID: %q", shareID)
	}
	if code != "abcd" {
		t.Fatalf("全角？后 code 提取不符: %q", code)
	}
}

func TestParseGuangyaShareURL_QueryShareID(t *testing.T) {
	shareID, code := parseGuangyaShareURL("https://www.guangyapan.com/share?shareId=sh123&code=zz99")
	if shareID != "sh123" {
		t.Fatalf("查询参数 shareId 优先: %q", shareID)
	}
	if code != "zz99" {
		t.Fatalf("code 不符: %q", code)
	}
}

func TestParseGuangyaShareURL_NoCode(t *testing.T) {
	shareID, code := parseGuangyaShareURL("https://www.guangyapan.com/s/only_token")
	if shareID != "only_token" || code != "" {
		t.Fatalf("无提取码场景不符: %q %q", shareID, code)
	}
}

func TestIsGuangyaURL(t *testing.T) {
	cases := map[string]bool{
		"https://www.guangyapan.com/s/x": true,
		"https://app.guangyapan.com/s/x": true,
		"https://pan.quark.cn/s/x":       false,
		"":                               false,
	}
	for url, want := range cases {
		if got := isGuangyaURL(url); got != want {
			t.Fatalf("isGuangyaURL(%q)=%v, want %v", url, got, want)
		}
	}
}

func TestExtractGuangyaPasscodeFromText(t *testing.T) {
	cases := map[string]string{
		"提取码：abcd":        "abcd",
		"（提取码: wxyz）":     "wxyz",
		"访问码：1234":         "1234",
		"资源链接 提取码:efgh 勿删": "efgh",
		"没有验证码信息":          "",
	}
	for text, want := range cases {
		if got := extractGuangyaPasscodeFromText(text); got != want {
			t.Fatalf("extract(%q)=%q, want %q", text, got, want)
		}
	}
}

// ============================================================================
// 响应判定（contracts 通用约定）
// ============================================================================

func mustPayload(t *testing.T, raw string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("测试载荷解析失败: %v raw=%s", err, raw)
	}
	return m
}

func TestGuangyaIsSuccess(t *testing.T) {
	success := []string{
		`{"msg":"success","data":{"shareUrl":"https://x"}}`,
		`{"code":0,"data":{"list":[]}}`,
		`{"status":200}`,
		`{"data":{"total":1}}`,
	}
	for _, raw := range success {
		if !guangyaIsSuccess(mustPayload(t, raw)) {
			t.Fatalf("应判定成功: %s", raw)
		}
	}
	fail := []string{
		`{"msg":"分享不存在"}`,
		`{"code":404,"msg":"fail"}`,
		`{}`,
	}
	for _, raw := range fail {
		if guangyaIsSuccess(mustPayload(t, raw)) {
			t.Fatalf("应判定失败: %s", raw)
		}
	}
}

func TestGuangyaExtractMessage(t *testing.T) {
	if got := guangyaExtractMessage(mustPayload(t, `{"msg":"分享不存在或已失效"}`)); got != "分享不存在或已失效" {
		t.Fatalf("msg 提取不符: %q", got)
	}
	if got := guangyaExtractMessage(mustPayload(t, `{"data":{"errorMessage":"提取码错误"}}`)); got != "提取码错误" {
		t.Fatalf("data.errorMessage 提取不符: %q", got)
	}
}

func TestGuangyaExtractItems(t *testing.T) {
	items := guangyaExtractItems(mustPayload(t, `{"msg":"success","data":{"list":[{"fid":"1"},{"fid":"2"}]}}`))
	if len(items) != 2 || items[0]["fid"] != "1" {
		t.Fatalf("data.list 提取不符: %+v", items)
	}
	arrItems := guangyaExtractItems(mustPayload(t, `{"msg":"success","data":[{"fid":"a"}]}`))
	if len(arrItems) != 1 || arrItems[0]["fid"] != "a" {
		t.Fatalf("data 数组提取不符: %+v", arrItems)
	}
}

func TestGuangyaItemFieldExtraction(t *testing.T) {
	item := map[string]interface{}{"fileId": "f1", "fileName": " movie.mp4 ", "dir": true}
	if got := guangyaItemFid(item); got != "f1" {
		t.Fatalf("fileId 容错提取不符: %q", got)
	}
	if got := guangyaItemName(item); got != "movie.mp4" {
		t.Fatalf("fileName 提取应去空白: %q", got)
	}
	if !guangyaItemIsDir(item) {
		t.Fatal("dir=true 应判定为目录")
	}
	fileItem := map[string]interface{}{"fid": "f2", "file_name": "a.txt", "dir": false}
	if guangyaItemIsDir(fileItem) {
		t.Fatal("dir=false 应判定为文件")
	}
	if got := guangyaItemFid(map[string]interface{}{}); got != "0" {
		t.Fatalf("空条目 fid 应归一为 0: %q", got)
	}
}

func TestGuangyaTaskStatusJudgement(t *testing.T) {
	done := mustPayload(t, `{"msg":"success","data":{"status":2}}`)
	if !guangyaTaskDone(done) || guangyaTaskFailed(done) {
		t.Fatal("status=2 应判定完成")
	}
	failed := mustPayload(t, `{"msg":"转存失败","data":{"status":5}}`)
	if !guangyaTaskFailed(failed) || guangyaTaskDone(failed) {
		t.Fatal("status=5 应判定失败")
	}
	running := mustPayload(t, `{"msg":"success","data":{"status":1}}`)
	if guangyaTaskDone(running) || guangyaTaskFailed(running) {
		t.Fatal("status=1 应判定进行中")
	}
}

// ============================================================================
// 创建分享载荷（D9 抓包原样，逐字段校验）
// ============================================================================

func TestBuildGuangyaSharePayload(t *testing.T) {
	payload := buildGuangyaSharePayload([]string{"1938809478796865571"}, "notify.svg")
	want := map[string]interface{}{
		"fileIds":          []interface{}{"1938809478796865571"},
		"title":            "notify.svg",
		"validateDuration": float64(0),
		"shareType":        float64(1),
		"autoFillCode":     true,
		"trafficLimit":     "0",
		"maxRestoreCount":  float64(0),
		"downloadType":     float64(1),
		"enableShareCode":  false,
		"shareCode":        "",
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("载荷序列化失败: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("载荷解析失败: %v", err)
	}
	// 逐字段（trafficLimit 必须是字符串型——抓包实证）
	for k, v := range want {
		if !reflect.DeepEqual(got[k], v) {
			t.Fatalf("载荷字段 %s = %#v, want %#v", k, got[k], v)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("载荷字段数不符: %d want %d", len(got), len(want))
	}
}

// ============================================================================
// diff 按文件名对齐（D4：转存接口不返回新 fid）
// ============================================================================

func item(fid, name string) map[string]interface{} {
	return map[string]interface{}{"fid": fid, "file_name": name, "dir": true}
}

func TestAlignGuangyaFidsFromItems(t *testing.T) {
	before := map[string]bool{"old1": true}
	items := []map[string]interface{}{
		item("old1", "旧文件"),   // 转存前已有 → 剔除
		item("new1", "剧集A"),    // 新增
		item("new2", "剧集B"),    // 新增
		item("new3", "广告.txt"), // 新增（对齐范围外不取）
	}
	got := alignGuangyaFidsFromItems(items, before, []string{"剧集A", "剧集B"})
	want := []string{"new1", "new2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("对齐不符: %v want %v", got, want)
	}
}

func TestAlignGuangyaFidsFromItems_SameNameTwice(t *testing.T) {
	// 同名文件出现两次：按未使用优先，各 fid 只用一次
	items := []map[string]interface{}{
		item("f1", "同名"),
		item("f2", "同名"),
	}
	got := alignGuangyaFidsFromItems(items, map[string]bool{}, []string{"同名", "同名"})
	if !reflect.DeepEqual(got, []string{"f1", "f2"}) {
		t.Fatalf("同名对齐不符: %v", got)
	}
}

func TestAlignGuangyaFidsFromItems_PartialMiss(t *testing.T) {
	items := []map[string]interface{}{
		item("f1", "存在的"),
	}
	got := alignGuangyaFidsFromItems(items, map[string]bool{}, []string{"存在的", "缺失的"})
	// 部分对齐：命中项返回，缺失项为空串（调用方重试/部分交付）
	if !reflect.DeepEqual(got, []string{"f1", ""}) {
		t.Fatalf("部分对齐不符: %v", got)
	}
}

func TestAlignGuangyaFidsFromItems_AllOld(t *testing.T) {
	items := []map[string]interface{}{
		item("old1", "文件"),
	}
	got := alignGuangyaFidsFromItems(items, map[string]bool{"old1": true}, []string{"文件"})
	if !reflect.DeepEqual(got, []string{""}) {
		t.Fatalf("全部为旧文件时对齐应为空: %v", got)
	}
}

// ============================================================================
// 面包屑工具
// ============================================================================

func TestGuangyaToBoolToInt(t *testing.T) {
	if !guangyaToBool("folder") || !guangyaToBool(true) || !guangyaToBool(float64(1)) {
		t.Fatal("guangyaToBool 判定不符")
	}
	if guangyaToBool("0") || guangyaToBool(false) {
		t.Fatal("guangyaToBool 假值判定不符")
	}
	if guangyaToInt("42") != 42 || guangyaToInt(float64(7)) != 7 || guangyaToInt("x") != 0 {
		t.Fatal("guangyaToInt 判定不符")
	}
}
