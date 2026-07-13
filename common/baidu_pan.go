package pan

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

const baiduPanBaseURL = "https://pan.baidu.com"

// baiduTransferDir 百度转存目标子目录。
// 百度对根目录 / 的文件删除有风控（errno=132 强制短信二次验证），
// 改用子目录规避。⚠️ 该目录必须由用户在百度网盘网页版手动创建，
// 百度 share/transfer API 不会自动创建目标目录，缺失时转存会失败。
const baiduTransferDir = "/urldb"

// BaiduPanService 百度网盘服务
type BaiduPanService struct {
	*BasePanService
}

// NewBaiduPanService 创建百度网盘服务
func NewBaiduPanService(config *PanConfig) *BaiduPanService {
	service := &BaiduPanService{
		BasePanService: NewBasePanService(config),
	}

	// 设置百度网盘的默认请求头（注意：不要设置 Content-Type，
	// 让 HTTPPost/HTTPPostForm 按需设置 application/json 或 application/x-www-form-urlencoded）
	service.SetHeaders(map[string]string{
		"Host":                      "pan.baidu.com",
		"Connection":                "keep-alive",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept-Encoding":           "gzip, deflate, br", // 触发 BasePanService 手动 gzip 解压
		"Referer":                   "https://pan.baidu.com",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "same-site",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})

	// 如果配置中带 Cookie，清理换行/制表符后设置到请求头
	// （参考 moss baidu_utils.go:115-118；用户从 DevTools 复制时常带入 \n\r\t 导致请求头格式错误）
	if config != nil && config.Cookie != "" {
		service.SetHeader("Cookie", SanitizeCookie(config.Cookie))
	}

	return service
}

// SanitizeCookie 清理 Cookie 字符串：移除换行/回车/制表符，trim 首尾空白。
// 用户从浏览器 DevHeaders 直接复制的 Cookie 经常夹带这些不可见字符，
// 导致 HTTP 请求头格式错误、百度返回 401 或 set-cookie 异常。
// 参考实现：moss baidu_utils.go:115-118。
func SanitizeCookie(raw string) string {
	s := strings.ReplaceAll(raw, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	return strings.TrimSpace(s)
}

// GetServiceType 获取服务类型
func (b *BaiduPanService) GetServiceType() ServiceType {
	return BaiduPan
}

// updateCookieValue 在 cookie 字符串中更新/新增一个键值对，返回新的 cookie 字符串。
// 纯函数，无副作用；用于 verifyPassCode 后回写 BDCLND。
func updateCookieValue(cookie, key, value string) string {
	m := map[string]string{}
	for _, pair := range strings.Split(cookie, ";") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	m[key] = value
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, "; ")
}

// parseBaiduErrno 解析百度响应 JSON，提取 errno 与整个 map。
func parseBaiduErrno(data []byte) (int, map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return 0, nil, fmt.Errorf("解析响应失败: %w", err)
	}
	errno := 0
	if v, ok := m["errno"].(float64); ok {
		errno = int(v)
	}
	return errno, m, nil
}

// getBdstoken 获取 bdstoken
func (b *BaiduPanService) getBdstoken() (string, error) {
	queryParams := map[string]string{
		"clienttype": "0",
		"app_id":     "250528",
		"web":        "1",
		"fields":     `["bdstoken","token","uk","isdocuser","servertime"]`,
	}
	data, err := b.HTTPGet(baiduPanBaseURL+"/api/gettemplatevariable", queryParams)
	if err != nil {
		return "", fmt.Errorf("获取 bdstoken 失败: %v", err)
	}

	errno, m, err := parseBaiduErrno(data)
	if err != nil {
		return "", err
	}
	if errno != 0 {
		return "", fmt.Errorf("获取 bdstoken 失败: %s", ErrnoMessage(errno))
	}

	result, ok := m["result"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("解析 bdstoken 失败")
	}
	bdstoken, ok := result["bdstoken"].(string)
	if !ok || bdstoken == "" {
		return "", fmt.Errorf("解析 bdstoken 失败")
	}
	return bdstoken, nil
}

// verifyPassCode 验证提取码，返回 randsk（用于回写 Cookie BDCLND）
func (b *BaiduPanService) verifyPassCode(surl, pwd, bdstoken string) (string, error) {
	queryParams := map[string]string{
		"surl":       surl,
		"bdstoken":   bdstoken,
		"t":          strconv.FormatInt(time.Now().UnixMilli(), 10),
		"channel":    "chunlei",
		"web":        "1",
		"clienttype": "0",
	}
	rawBody := fmt.Sprintf("pwd=%s&vcode=&vcode_str=", pwd)

	data, err := b.HTTPPostForm(baiduPanBaseURL+"/share/verify", rawBody, queryParams)
	if err != nil {
		return "", fmt.Errorf("验证提取码失败: %v", err)
	}

	errno, m, err := parseBaiduErrno(data)
	if err != nil {
		return "", err
	}
	if errno != 0 {
		return "", fmt.Errorf("验证提取码失败: %s", ErrnoMessage(errno))
	}

	randsk, ok := m["randsk"].(string)
	if !ok {
		return "", fmt.Errorf("解析 randsk 失败")
	}
	return randsk, nil
}

// getSharedPaths 解析百度分享页，返回文件信息列表
func (b *BaiduPanService) getSharedPaths(shareURL string) ([]baiduShareFile, error) {
	surl := ExtractSurl(shareURL)
	if surl == "" {
		return nil, fmt.Errorf("无效的分享链接")
	}
	body, err := b.HTTPGet(shareURL, nil)
	if err != nil {
		return nil, fmt.Errorf("获取分享页失败: %v", err)
	}
	return ParseSharePageHTML(string(body))
}

// transferFile 转存文件到指定目录，返回转存后的 to_fs_id 列表（仅用于确认成功/计数）
func (b *BaiduPanService) transferFile(shareID, uk int64, bdstoken, remotedir string, fsIDs []int64) ([]int64, error) {
	queryParams := map[string]string{
		"shareid":    strconv.FormatInt(shareID, 10),
		"from":       strconv.FormatInt(uk, 10),
		"bdstoken":   bdstoken,
		"channel":    "chunlei",
		"web":        "1",
		"clienttype": "0",
	}

	// 构建 fsidlist JSON 数组 [1,2,3]
	fsidBuilder := strings.Builder{}
	fsidBuilder.WriteByte('[')
	for i, id := range fsIDs {
		if i > 0 {
			fsidBuilder.WriteByte(',')
		}
		fsidBuilder.WriteString(strconv.FormatInt(id, 10))
	}
	fsidBuilder.WriteByte(']')
	fsidlistJSON := fsidBuilder.String()

	rawBody := fmt.Sprintf("fsidlist=%s&path=%s", fsidlistJSON, url.QueryEscape(remotedir))

	data, err := b.HTTPPostForm(baiduPanBaseURL+"/share/transfer", rawBody, queryParams)
	if err != nil {
		return nil, fmt.Errorf("转存失败: %v", err)
	}

	errno, m, err := parseBaiduErrno(data)
	if err != nil {
		return nil, err
	}
	if errno != 0 {
		return nil, fmt.Errorf("转存失败: %s", ErrnoMessage(errno))
	}

	var toFsIDs []int64
	if extra, ok := m["extra"].(map[string]any); ok {
		if list, ok := extra["list"].([]any); ok {
			for _, item := range list {
				if im, ok := item.(map[string]any); ok {
					if v, ok := im["to_fs_id"].(float64); ok {
						toFsIDs = append(toFsIDs, int64(v))
					}
				}
			}
		}
	}
	if len(toFsIDs) == 0 {
		return nil, fmt.Errorf("未能获取转存后的文件信息")
	}
	return toFsIDs, nil
}

// createShare 创建分享，返回分享链接（带提取码）
func (b *BaiduPanService) createShare(fsIDs []int64, period, pwd, bdstoken string) (string, error) {
	queryParams := map[string]string{
		"channel":    "chunlei",
		"bdstoken":   bdstoken,
		"clienttype": "0",
		"app_id":     "250528",
		"web":        "1",
		"dp-logid":   strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	// 构建 fid_list JSON 数组
	fidBuilder := strings.Builder{}
	fidBuilder.WriteByte('[')
	for i, id := range fsIDs {
		if i > 0 {
			fidBuilder.WriteByte(',')
		}
		fidBuilder.WriteString(strconv.FormatInt(id, 10))
	}
	fidBuilder.WriteByte(']')
	fidListJSON := fidBuilder.String()

	rawBody := fmt.Sprintf("period=%s&pwd=%s&eflag_disable=true&channel_list=[]&schannel=4&fid_list=%s", period, pwd, fidListJSON)

	data, err := b.HTTPPostForm(baiduPanBaseURL+"/share/pset", rawBody, queryParams)
	if err != nil {
		return "", fmt.Errorf("创建分享失败: %v", err)
	}

	errno, m, err := parseBaiduErrno(data)
	if err != nil {
		return "", err
	}
	if errno != 0 {
		return "", fmt.Errorf("创建分享失败: %s", ErrnoMessage(errno))
	}

	link, ok := m["link"].(string)
	if !ok || link == "" {
		return "", fmt.Errorf("解析分享链接失败")
	}
	shareURL := link
	if pwd != "" {
		shareURL += "?pwd=" + pwd
	}
	return shareURL, nil
}

// Share 对系统已存文件按 fid 重新生成百度分享链接（实现 Sharer，FR-015）。
// fid 即百度 fs_id（int64，存为字符串）；period=1 永久分享。失败时由调用方回退到「判原始→转存」。
func (b *BaiduPanService) Share(fid string) (*TransferResult, error) {
	if fid == "" {
		return &TransferResult{Success: false, Message: "fid 为空"}, nil
	}
	fsID, err := strconv.ParseInt(fid, 10, 64)
	if err != nil {
		return &TransferResult{Success: false, Message: fmt.Sprintf("fid 非法: %v", err)}, nil
	}
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return &TransferResult{Success: false, Message: fmt.Sprintf("获取 bdstoken 失败: %v", err)}, nil
	}
	shareURL, err := b.createShare([]int64{fsID}, "1", "", bdstoken)
	if err != nil {
		return &TransferResult{Success: false, Message: fmt.Sprintf("创建分享失败: %v", err)}, nil
	}
	if shareURL == "" {
		return &TransferResult{Success: false, Message: "未获取到分享链接"}, nil
	}
	return &TransferResult{Success: true, ShareURL: shareURL, Fid: fid}, nil
}

// deleteByPaths 按路径批量删除文件（百度删除 API 基于路径，非 fs_id）
//
// 关键：参数与 filelist 格式必须与浏览器抓包一致，否则触发 errno=132 短信二次验证：
//   - filelist 用字符串数组 ["path1","path2"]，不是对象数组 [{"path":"..."}]
//   - query 必须带 async=2 / onnest=fail / newVerify=1 / dp-logid，不能带 channel=chunlei
//   - async=2 模式下返回 taskid，需要轮询 /share/taskquery 确认实际删除结果
func (b *BaiduPanService) deleteByPaths(paths []string) error {
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return err
	}

	// filelist 是字符串数组（不是对象数组），与浏览器一致
	// 用 json.Marshal 避免文件名特殊字符破坏 JSON
	filelistJSON, err := json.Marshal(paths)
	if err != nil {
		return fmt.Errorf("构建删除列表失败: %v", err)
	}

	// dp-logid：浏览器随机生成，用纳秒时间戳模拟（非关键，但缺失可能被风控识别）
	dpLogid := strconv.FormatInt(time.Now().UnixNano(), 10)

	queryParams := map[string]string{
		"async":      "2",
		"onnest":     "fail",
		"opera":      "delete",
		"bdstoken":   bdstoken,
		"newVerify":  "1",
		"clienttype": "0",
		"app_id":     "250528",
		"web":        "1",
		"dp-logid":   dpLogid,
	}
	rawBody := "filelist=" + url.QueryEscape(string(filelistJSON))

	// 诊断日志：记录请求体与目标路径，便于排查 errno
	log.Printf("[baidu] deleteByPaths 请求 paths=%v body=%s", paths, rawBody)

	data, err := b.HTTPPostForm(baiduPanBaseURL+"/api/filemanager", rawBody, queryParams)
	if err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}

	// 诊断日志：记录原始响应前 500 字符
	respSnippet := string(data)
	if len(respSnippet) > 500 {
		respSnippet = respSnippet[:500]
	}
	log.Printf("[baidu] deleteByPaths 响应: %s", respSnippet)

	errno, m, err := parseBaiduErrno(data)
	if err != nil {
		return err
	}
	// 防御：errno 字段缺失（异常空响应，如被中间层拦截返回 {}）不应误判为删除成功
	if _, hasErrno := m["errno"]; !hasErrno {
		return fmt.Errorf("删除响应异常，缺少 errno 字段")
	}
	if errno != 0 {
		// 错误消息保持 ErrnoMessage 原文，便于 cleanup_service.isFileNotExist 匹配。
		// 注意 errno=132 是百度风控要求短信二次验证（响应含 verify_scene/authwidget），
		// 不是"文件不存在"—— historically 这里误判过。
		return fmt.Errorf("%s", ErrnoMessage(errno))
	}

	// async=2 模式：返回 taskid 表示异步任务已受理，需要轮询 /share/taskquery 确认实际结果。
	// 浏览器流程：filemanager 返回 taskid → 轮询 taskquery 直到 state=success → 刷新文件列表。
	if taskID, ok := m["taskid"]; ok && taskID != nil {
		if err := b.pollDeleteTask(taskID); err != nil {
			return err
		}
	}
	return nil
}

// pollDeleteTask 轮询异步删除任务结果（async=2 模式）。
//
// 重要：filemanager 返回 errno=0 + taskid 已表示删除请求被百度受理并执行，
// taskquery 仅用于显式确认，其失败不推翻受理结果（百度 API 怪异行为：
// 删除已完成后查询可能返回非 0 errno，实测文件已被删除但 taskquery 报 errno=4）。
//
// 响应字段语义（基于浏览器抓包）：
//   - errno：外层 API 调用结果，0 = 调用成功
//   - status：任务状态字符串，"success" = 任务完成
//   - task_errno：任务实际执行结果，0 = 任务成功（文件已删）
//
// 策略：只要看到 errno=0 + status=success + task_errno=0 就显式成功；
// 其他情况（包括 taskquery 自己返回非 0 errno）继续轮询；
// 超时后视为成功（因为 filemanager 已受理，乐观处理避免假阴性）。
func (b *BaiduPanService) pollDeleteTask(taskID any) error {
	taskIDStr := fmt.Sprintf("%v", taskID)
	queryParams := map[string]string{
		"taskid":     taskIDStr,
		"clienttype": "0",
		"app_id":     "250528",
		"web":        "1",
	}

	const maxAttempts = 10
	const interval = 500 * time.Millisecond

	for i := 0; i < maxAttempts; i++ {
		data, err := b.HTTPGet(baiduPanBaseURL+"/share/taskquery", queryParams)
		if err != nil {
			// 网络错误不直接失败，继续重试
			log.Printf("[baidu] taskquery 网络错误 (attempt=%d): %v", i+1, err)
			time.Sleep(interval)
			continue
		}

		errno, m, perr := parseBaiduErrno(data)
		if perr != nil {
			log.Printf("[baidu] taskquery 解析失败 (attempt=%d): %v", i+1, perr)
			time.Sleep(interval)
			continue
		}

		// 只有 errno=0 时才进一步看 status/task_errno
		if errno == 0 {
			status, _ := m["status"].(string)
			if status == "success" {
				taskErrno := toInt64(m["task_errno"])
				if taskErrno != 0 {
					// 任务执行层失败（文件不存在等）。由于 filemanager 已受理，
					// 这里仍乐观视为成功——文件可能本来就不存在，cleanup 目标已达成。
					log.Printf("[baidu] taskquery task_errno=%d (taskid=%s)，乐观视为成功", taskErrno, taskIDStr)
					return nil
				}
				log.Printf("[baidu] taskquery 任务完成 (taskid=%s)", taskIDStr)
				return nil
			}
			// status 非 success = 任务仍在进行，继续等
		}
		// errno != 0：任务查询层错误（百度怪异行为，删除已完成后查询可能返回非 0），
		// 继续重试确认，超时再乐观处理
		time.Sleep(interval)
	}

	// 超时：filemanager 已受理，乐观视为成功（避免假阴性）。
	// 下一轮 cleanup 不会扫到该资源（因为外层 cleanup 会清空 fid/save_url）。
	log.Printf("[baidu] taskquery 超时未确认 (taskid=%s)，乐观视为成功（filemanager 已受理）", taskIDStr)
	return nil
}

// Transfer 转存分享链接
func (b *BaiduPanService) Transfer(shareID string) (*TransferResult, error) {
	config := b.GetConfig()
	if config == nil {
		return ErrorResult("百度转存配置为空"), nil
	}

	// 1. bdstoken
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("Cookie 失效或网络错误: %v", err)), nil
	}

	// 2. surl
	surl := ExtractSurl(config.URL)
	if surl == "" {
		return ErrorResult("无效的百度分享链接"), nil
	}

	// 3. 提取码验证（若有）
	// 上层调用方（transferToCloud）目前不会把 URL 里的 ?pwd=xxx 拆到 Code 字段，
	// 这里兜底：若 Code 为空，尝试从 URL 的 query 中解析 pwd。
	code := config.Code
	if code == "" {
		if u, perr := url.Parse(config.URL); perr == nil {
			code = u.Query().Get("pwd")
		}
	}
	if code != "" {
		randsk, err := b.verifyPassCode(surl, code, bdstoken)
		if err != nil {
			return ErrorResult(fmt.Sprintf("提取码错误或链接失效: %v", err)), nil
		}
		// 回写 BDCLND 到 Cookie 头
		newCookie := updateCookieValue(b.GetHeader("Cookie"), "BDCLND", randsk)
		b.SetHeader("Cookie", newCookie)
		if config.Cookie != "" {
			config.Cookie = newCookie
		}
	}

	// 4. 解析分享页
	files, err := b.getSharedPaths(config.URL)
	if err != nil || len(files) == 0 {
		msg := "解析分享链接失败"
		if err != nil {
			msg = err.Error()
		}
		return ErrorResult(msg), nil
	}

	title := files[0].ServerName

	// 5a. IsType==1：scheduler 只校验+取标题，不真转存
	if config.IsType == 1 {
		return SuccessResult("校验成功", map[string]interface{}{
			"title":    title,
			"shareUrl": config.URL,
		}), nil
	}

	// 5b. IsType==0：真转存 + 重新分享
	fsIDs := make([]int64, 0, len(files))
	for _, f := range files {
		fsIDs = append(fsIDs, f.FsID)
	}

	toFsIDs, err := b.transferFile(files[0].ShareID, files[0].UK, bdstoken, baiduTransferDir, fsIDs)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}

	// 有效期：统一永久（period=0）
	period := "0"
	shareURL, err := b.createShare(toFsIDs, period, code, bdstoken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存成功但创建分享失败: %v", err)), nil
	}

	// fid 存完整路径（百度删除基于路径）。
	// 关键：必须带上 baiduTransferDir 前缀，否则 cleanup 删除时路径错位；
	// 且根目录删除会触发 errno=132 风控，所以一定要落在子目录里。
	// 多文件用逗号连接。
	pathParts := make([]string, 0, len(files))
	for _, f := range files {
		pathParts = append(pathParts, baiduTransferDir+"/"+f.ServerName)
	}
	fid := strings.Join(pathParts, ",")

	return SuccessResult("转存成功", map[string]interface{}{
		"shareUrl": shareURL,
		"title":    title,
		"fid":      fid,
		"code":     code,
	}), nil
}

// DeleteFiles 删除文件（fileList 元素可能是逗号连接的多路径）
func (b *BaiduPanService) DeleteFiles(fileList []string) (*TransferResult, error) {
	if len(fileList) == 0 {
		return ErrorResult("文件列表为空"), nil
	}
	// fileList 元素可能是逗号连接的多路径（与 Transfer 存储一致）
	var paths []string
	for _, item := range fileList {
		for _, p := range strings.Split(item, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				paths = append(paths, p)
			}
		}
	}
	if len(paths) == 0 {
		return ErrorResult("文件列表为空"), nil
	}
	if err := b.deleteByPaths(paths); err != nil {
		// 错误消息保持 ErrnoMessage 原文，便于 cleanup_service.isFileNotExist 匹配
		return ErrorResult(err.Error()), nil
	}
	return SuccessResult("删除成功", nil), nil
}

// GetFiles 获取文件列表
func (b *BaiduPanService) GetFiles(pdirFid string) (*TransferResult, error) {
	if pdirFid == "" {
		pdirFid = "/"
	}

	bdstoken, err := b.getBdstoken()
	if err != nil {
		return ErrorResult(fmt.Sprintf("Cookie 失效或网络错误: %v", err)), nil
	}

	queryParams := map[string]string{
		"order":     "time",
		"desc":      "1",
		"showempty": "0",
		"web":       "1",
		"page":      "1",
		"num":       "1000",
		"dir":       pdirFid,
		"bdstoken":  bdstoken,
	}

	data, err := b.HTTPGet(baiduPanBaseURL+"/api/list", queryParams)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取文件列表失败: %v", err)), nil
	}

	errno, m, err := parseBaiduErrno(data)
	if err != nil {
		return ErrorResult("解析响应失败"), nil
	}
	if errno != 0 {
		return ErrorResult(ErrnoMessage(errno)), nil
	}

	list, _ := m["list"].([]any)
	return SuccessResult("获取成功", list), nil
}

// GetUserInfo 获取用户信息
//
// 百度网盘网页版的用户信息分散在两个 API：
//   1. /api/gettemplatevariable —— 仅返回 username / uk / vip_type（不含容量）
//   2. /api/quota               —— 返回 total / used（字节数，int64）
//
// 历史代码试图用单个 API 拿全部字段（fields 里塞 total_capacity/used_capacity），
// 但百度会静默忽略不支持的 field，导致容量恒为 0。这里改成两次请求组合。
func (b *BaiduPanService) GetUserInfo(cookie *string) (*UserInfo, error) {
	// 设置Cookie（清理换行/制表符，防御历史脏数据或 DevTools 复制残留）
	if cookie != nil && *cookie != "" {
		b.SetHeader("Cookie", SanitizeCookie(*cookie))
	}

	// Step 1: 用户基本信息（username + vip_type）
	userParams := map[string]string{
		"clienttype": "0",
		"app_id":     "250528",
		"web":        "1",
		"fields":     `["username","uk","vip_type"]`,
	}
	userResp, err := b.HTTPGet("https://pan.baidu.com/api/gettemplatevariable", userParams)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 用 map[string]any 解析，避免百度 API 字段类型漂移
	// （uk 实际是 number，但历史代码曾定义为 string 导致反序列化失败）
	var userRaw map[string]any
	if err := b.ParseJSONResponse(userResp, &userRaw); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}
	if errno := toInt64(userRaw["errno"]); errno != 0 {
		return nil, fmt.Errorf("API返回错误: %s", ErrnoMessage(int(errno)))
	}
	userResult, ok := userRaw["result"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("用户信息响应缺少 result 字段")
	}
	username, _ := userResult["username"].(string)
	vipType := toInt64(userResult["vip_type"])

	// Step 2: 容量信息（total / used，单位字节）
	quotaParams := map[string]string{
		"clienttype": "0",
		"app_id":     "250528",
		"web":        "1",
	}
	quotaResp, err := b.HTTPGet("https://pan.baidu.com/api/quota", quotaParams)
	if err != nil {
		// 容量查询失败不阻塞账号创建/刷新，仅记日志返回 0
		log.Printf("[baidu] 获取容量失败（不阻塞）: %v", err)
	} else {
		var quotaRaw map[string]any
		if err := b.ParseJSONResponse(quotaResp, &quotaRaw); err == nil {
			if errno := toInt64(quotaRaw["errno"]); errno == 0 {
				// 写回 userResult 以便下面统一取值
				userResult["total"] = quotaRaw["total"]
				userResult["used"] = quotaRaw["used"]
			} else {
				log.Printf("[baidu] quota API errno=%d: %s", errno, ErrnoMessage(int(errno)))
			}
		}
	}

	totalCapacity := toInt64(userResult["total"])
	usedCapacity := toInt64(userResult["used"])

	return &UserInfo{
		Username:    username,
		VIPStatus:   vipType > 0,
		UsedSpace:   usedCapacity,
		TotalSpace:  totalCapacity,
		ServiceType: "baidu",
	}, nil
}

// toInt64 把 any 安全转成 int64（兼容 float64 / int / json.Number / 数字字符串）。
// 百度 API 同一字段在不同账号/版本下可能返回 number 或 string-number，统一兜底。
func toInt64(v any) int64 {
	switch n := v.(type) {
	case nil:
		return 0
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	case string:
		i, _ := strconv.ParseInt(n, 10, 64)
		return i
	}
	return 0
}

// GetUserInfoByEntity 根据 entity.Cks 获取用户信息（待实现）
func (b *BaiduPanService) GetUserInfoByEntity(cks entity.Cks) (*UserInfo, error) {
	return nil, nil
}

func (u *BaiduPanService) SetCKSRepository(cksRepo repo.CksRepository, entity entity.Cks) {
}

func (x *BaiduPanService) UpdateConfig(config *PanConfig) {
	if config == nil {
		return
	}
	x.config = config
	if config.Cookie != "" {
		x.SetHeader("Cookie", config.Cookie)
	}
}
