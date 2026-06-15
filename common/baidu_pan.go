package pan

import (
	"encoding/json"
	"fmt"
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

	// 如果配置中带 Cookie，设置到请求头
	if config != nil && config.Cookie != "" {
		service.SetHeader("Cookie", config.Cookie)
	}

	return service
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
	filelistBuilder := strings.Builder{}
	filelistBuilder.WriteByte('[')
	for i, p := range paths {
		if i > 0 {
			filelistBuilder.WriteByte(',')
		}
		filelistBuilder.WriteString(`{"path":"`)
		filelistBuilder.WriteString(p)
		filelistBuilder.WriteString(`"}`)
	}
	filelistBuilder.WriteByte(']')
	filelistJSON := filelistBuilder.String()

	queryParams := map[string]string{
		"opera":       "delete",
		"bdstoken":    bdstoken,
		"web":         "1",
		"clienttype":  "0",
		"channel":     "chunlei",
		"app_id":      "250528",
	}
	rawBody := "filelist=" + url.QueryEscape(filelistJSON)

	data, err := b.HTTPPostForm(baiduPanBaseURL+"/api/filemanager", rawBody, queryParams)
	if err != nil {
		return fmt.Errorf("%s", "删除失败: "+err.Error())
	}

	errno, _, err := parseBaiduErrno(data)
	if err != nil {
		return err
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
		"order":      "time",
		"desc":       "1",
		"showempty":  "0",
		"web":        "1",
		"page":       "1",
		"num":        "1000",
		"dir":        pdirFid,
		"bdstoken":   bdstoken,
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
func (b *BaiduPanService) GetUserInfo(cookie *string) (*UserInfo, error) {
	// 设置Cookie
	b.SetHeader("Cookie", *cookie)

	// 调用百度网盘用户信息API
	userInfoURL := "https://pan.baidu.com/api/gettemplatevariable"
	data := map[string]interface{}{
		"fields": "['username','uk','vip_type','vip_endtime','total_capacity','used_capacity']",
	}

	resp, err := b.HTTPPost(userInfoURL, data, nil)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 解析响应
	var result struct {
		Errno int `json:"errno"`
		Data  struct {
			Username      string `json:"username"`
			Uk            string `json:"uk"`
			VipType       int    `json:"vip_type"`
			VipEndtime    string `json:"vip_endtime"`
			TotalCapacity string `json:"total_capacity"`
			UsedCapacity  string `json:"used_capacity"`
		} `json:"data"`
	}

	if err := b.ParseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	if result.Errno != 0 {
		return nil, fmt.Errorf("API返回错误: %d", result.Errno)
	}

	// 转换VIP状态
	vipStatus := result.Data.VipType > 0

	// 解析容量字符串
	totalCapacityBytes := ParseCapacityString(result.Data.TotalCapacity)
	usedCapacityBytes := ParseCapacityString(result.Data.UsedCapacity)

	return &UserInfo{
		Username:    result.Data.Username,
		VIPStatus:   vipStatus,
		UsedSpace:   usedCapacityBytes,
		TotalSpace:  totalCapacityBytes,
		ServiceType: "baidu",
	}, nil
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
