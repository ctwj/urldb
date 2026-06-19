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

// deleteByPaths 按路径批量删除文件（百度删除 API 基于路径，非 fs_id）
func (b *BaiduPanService) deleteByPaths(paths []string) error {
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return err
	}

	// 构建 filelist JSON: [{"path":"/a"},{"path":"/b"}]
	// 用 json.Marshal 而非手拼，避免文件名特殊字符（如 "）破坏 JSON。
	type pathItem struct {
		Path string `json:"path"`
	}
	items := make([]pathItem, 0, len(paths))
	for _, p := range paths {
		items = append(items, pathItem{Path: p})
	}
	filelistJSON, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("构建删除列表失败: %v", err)
	}

	queryParams := map[string]string{
		"opera":      "delete",
		"bdstoken":   bdstoken,
		"web":        "1",
		"clienttype": "0",
		"channel":    "chunlei",
		"app_id":     "250528",
	}
	rawBody := "filelist=" + url.QueryEscape(string(filelistJSON))

	data, err := b.HTTPPostForm(baiduPanBaseURL+"/api/filemanager", rawBody, queryParams)
	if err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}

	errno, m, err := parseBaiduErrno(data)
	if err != nil {
		return err
	}
	// 防御：errno 字段缺失（异常空响应，如被中间层拦截返回 {}）不应误判为删除成功
	if _, hasErrno := m["errno"]; !hasErrno {
		return fmt.Errorf("删除响应异常，缺少 errno 字段")
	}
	if errno != 0 {
		// 错误消息保持 ErrnoMessage 原文，便于 cleanup_service.isFileNotExist 匹配
		return fmt.Errorf("%s", ErrnoMessage(errno))
	}
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
	if config.Code != "" {
		randsk, err := b.verifyPassCode(surl, config.Code, bdstoken)
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

	toFsIDs, err := b.transferFile(files[0].ShareID, files[0].UK, bdstoken, "/", fsIDs)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}

	// 有效期：统一永久（period=0）
	period := "0"
	shareURL, err := b.createShare(toFsIDs, period, config.Code, bdstoken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存成功但创建分享失败: %v", err)), nil
	}

	// fid 存路径（百度删除基于路径）。多文件用逗号连接
	pathParts := make([]string, 0, len(files))
	for _, f := range files {
		pathParts = append(pathParts, "/"+f.ServerName)
	}
	fid := strings.Join(pathParts, ",")

	return SuccessResult("转存成功", map[string]interface{}{
		"shareUrl": shareURL,
		"title":    title,
		"fid":      fid,
		"code":     config.Code,
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
