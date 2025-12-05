package pan

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// BaiduPanService 百度网盘服务
type BaiduPanService struct {
	*BasePanService
	configMutex sync.RWMutex // 保护配置的读写锁
	bdstoken    string       // 百度网盘操作令牌
	entity      entity.Cks   // 实体信息
	cksRepo     repo.CksRepository
}

// NewBaiduPanService 创建百度网盘服务
func NewBaiduPanService(config *PanConfig) *BaiduPanService {
	service := &BaiduPanService{
		BasePanService: NewBasePanService(config),
	}

	// 设置百度网盘的默认请求头
	service.SetHeaders(map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Content-Type":    "application/x-www-form-urlencoded",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Referer":         "https://pan.baidu.com/",
		"Origin":          "https://pan.baidu.com",
	})

	service.UpdateConfig(config)
	return service
}

// GetServiceType 获取服务类型
func (b *BaiduPanService) GetServiceType() ServiceType {
	return BaiduPan
}

// Transfer 转存分享链接
func (b *BaiduPanService) Transfer(shareID string) (*TransferResult, error) {
	b.configMutex.RLock()
	config := b.config
	b.configMutex.RUnlock()

	log.Printf("开始处理百度网盘分享: %s", shareID)

	// 解析分享链接
	urlParts, err := url.Parse(config.URL)
	if err != nil {
		return ErrorResult(fmt.Sprintf("解析分享链接失败: %v", err)), nil
	}

	linkURL := urlParts.Scheme + "://" + urlParts.Host + urlParts.Path
	passCode := extractCode(config.URL)

	// 验证提取码（如果有）
	if passCode != "" {
		randsk, err := b.verifyPassCode(linkURL, passCode)
		if err != nil {
			return ErrorResult(fmt.Sprintf("验证提取码失败: %v", err)), nil
		}
		if randsk == "0" {
			return ErrorResult("提取码验证失败"), nil
		}
		// 更新cookie
		b.updateBdclnd(randsk)
	}

	// 获取转存参数
	transferParams, err := b.getTransferParams(linkURL)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取转存参数失败: %v", err)), nil
	}

	shareID, userID, fsIds, fileNames := transferParams[0], transferParams[1], transferParams[2], transferParams[3]

	// 检验模式
	if config.IsType == 1 {
		fileNameList := strings.Split(fileNames, ",")
		title := ""
		if len(fileNameList) > 0 {
			title = fileNameList[0]
		}
		return SuccessResult("检验成功", map[string]interface{}{
			"title":     title,
			"shareUrl":  config.URL,
			"stoken":    "",
		}), nil
	}

	// 确定存储路径
	folderName := b.getStorageFolder()

	// 检查并创建目录
	if err := b.ensureDirectoryExists(folderName); err != nil {
		return ErrorResult(fmt.Sprintf("创建目录失败: %v", err)), nil
	}

	// 执行转存
	if err := b.transferFile([]string{shareID, userID, fsIds}, "/"+folderName); err != nil {
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}

	// 获取转存后的文件列表
	dirList, err := b.getDirList("/" + folderName)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取文件列表失败: %v", err)), nil
	}

	// 处理广告过滤和文件整理
	filePaths, fsIdList, err := b.processTransferredFiles(dirList, strings.Split(fileNames, ","), folderName)
	if err != nil {
		return ErrorResult(fmt.Sprintf("处理转存文件失败: %v", err)), nil
	}

	if len(filePaths) == 0 {
		return ErrorResult("资源内容为空或所有转存的文件都包含广告内容，已全部删除"), nil
	}

	// 创建分享
	password := "6666"
	shareLink, err := b.createShare(strings.Join(fsIdList, ","), 0, password)
	if err != nil {
		return ErrorResult(fmt.Sprintf("创建分享失败: %v", err)), nil
	}

	if shareLink == "0" {
		return ErrorResult("创建分享失败"), nil
	}

	if password != "" {
		shareLink = shareLink + "?pwd=" + password
	}

	fileNameList := strings.Split(fileNames, ",")
	title := ""
	if len(fileNameList) > 0 {
		title = fileNameList[0]
	}

	return SuccessResult("文件转存成功", map[string]interface{}{
		"title":    title,
		"shareUrl": shareLink,
		"fid":      filePaths,
		"code":     password,
	}), nil
}

// GetFiles 获取文件列表
func (b *BaiduPanService) GetFiles(pdirFid string) (*TransferResult, error) {
	if pdirFid == "0" || pdirFid == "" {
		pdirFid = "/"
	}

	log.Printf("开始获取百度网盘文件列表，目录ID: %s", pdirFid)

	dirList, err := b.getDirList(pdirFid)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取文件列表失败: %v", err)), nil
	}

	return SuccessResult("获取成功", dirList), nil
}

// DeleteFiles 删除文件
func (b *BaiduPanService) DeleteFiles(fileList []string) (*TransferResult, error) {
	if len(fileList) == 0 {
		return ErrorResult("文件列表为空"), nil
	}

	log.Printf("开始删除百度网盘文件，文件数量: %d", len(fileList))

	result, err := b.batchDeleteFiles(fileList)
	if err != nil {
		return ErrorResult(fmt.Sprintf("删除文件失败: %v", err)), nil
	}

	// 检查删除结果
	if errno, ok := result["errno"].(float64); ok && int(errno) != 0 {
		return ErrorResult(fmt.Sprintf("删除文件失败: %v", result)), nil
	}

	return SuccessResult("删除成功", nil), nil
}

// GetUserInfo 获取用户信息 - 与PHP实现完全一致
func (b *BaiduPanService) GetUserInfo(cookie *string) (*UserInfo, error) {
	// 设置Cookie
	b.SetHeader("Cookie", *cookie)

	// 首先调用gettemplatevariable获取基础token信息，与PHP实现一致
	userInfoURL := "https://pan.baidu.com/api/gettemplatevariable"

	utils.Info("调用百度网盘用户信息API - URL: %s, Cookie长度: %d", userInfoURL, len(*cookie))
	resp, err := b.HTTPGet(userInfoURL, map[string]string{
		"fields":     `["bdstoken","token","uk","isdocuser","servertime"]`,
		"clienttype": "0",
		"app_id":     "38824127",
		"web":        "1",
	})
	if err != nil {
		utils.Info("百度网盘API调用失败: %v", err)
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	utils.Info("百度网盘API调用成功，响应长度: %d", len(resp))

	// 解析响应 - 根据PHP实现调整结构
	var result struct {
		Errno int `json:"errno"`
		Data  struct {
			Bdstoken   string `json:"bdstoken"`
			Token      string `json:"token"`
			UK         string `json:"uk"`
			IsDocUser  int    `json:"isdocuser"`
			ServerTime int64  `json:"servertime"`
		} `json:"data"`
	}

	utils.Info("API响应内容: %s", string(resp))

	if err := b.ParseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	if result.Errno != 0 {
		return nil, fmt.Errorf("API返回错误: %d", result.Errno)
	}

	utils.Info("成功获取用户信息token: %+v", result)

	// 现在使用获取的bdstoken来获取详细的用户信息
	// 这里需要调用另一个API来获取容量等详细信息
	detailUserInfo, err := b.getUserInfoDetail(cookie, result.Data.Bdstoken, result.Data.UK)
	if err != nil {
		utils.Error("获取详细用户信息失败: %v", err)
		// 如果获取详细信息失败，返回基本信息
		return &UserInfo{
			Username:    "", // PHP版本在这个阶段也没有username
			VIPStatus:   false, // 默认为false，后续可以更新
			UsedSpace:   0,
			TotalSpace:  0,
			ServiceType: "baidu",
		}, nil
	}

	return detailUserInfo, nil
}

// getUserInfoDetail 获取详细用户信息
func (b *BaiduPanService) getUserInfoDetail(cookie *string, bdstoken, uk string) (*UserInfo, error) {
	// 设置Cookie
	b.SetHeader("Cookie", *cookie)

	// 调用百度网盘用户信息API
	userInfoURL := "https://pan.baidu.com/rest/2.0/xpan/nas"
	params := map[string]string{
		"method":     "uinfo",
		"clienttype": "0",
		"app_id":     "250528", // 注意这个app_id与gettemplatevariable不同
	}

	utils.Info("调用百度网盘详细用户信息API - URL: %s", userInfoURL)
	resp, err := b.HTTPGet(userInfoURL, params)
	if err != nil {
		utils.Info("百度网盘详细用户信息API调用失败: %v", err)
		return nil, fmt.Errorf("获取详细用户信息失败: %v", err)
	}

	// 解析响应
	var result struct {
		Errno int `json:"errno"`
		Username string `json:"username"`
		UK      string `json:"uk"`
		VipType int    `json:"vip_type"`
		VipEndtime int64 `json:"vip_endtime"`
		TotalCapacity int64 `json:"total_capacity"`
		UsedCapacity  int64 `json:"used_capacity"`
	}

	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return nil, fmt.Errorf("解析详细用户信息失败: %v", err)
	}

	if result.Errno != 0 {
		return nil, fmt.Errorf("详细用户信息API返回错误: %d", result.Errno)
	}

	// 转换VIP状态
	vipStatus := result.VipType > 0

	utils.Info("成功获取详细用户信息: username=%s, uk=%s, vip=%d", result.Username, result.UK, result.VipType)

	return &UserInfo{
		Username:    result.Username,
		VIPStatus:   vipStatus,
		UsedSpace:   result.UsedCapacity,
		TotalSpace:  result.TotalCapacity,
		ServiceType: "baidu",
	}, nil
}

// buildQuery 构建查询字符串
func (b *BaiduPanService) buildQuery(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	return values.Encode()
}

// SetCKSRepository 设置 CksRepository 和 entity
func (b *BaiduPanService) SetCKSRepository(cksRepo repo.CksRepository, entity entity.Cks) {
	b.cksRepo = cksRepo
	b.entity = entity
}

// UpdateConfig 更新配置（线程安全）
func (b *BaiduPanService) UpdateConfig(config *PanConfig) {
	if config == nil {
		return
	}
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	b.config = config
	if config.Cookie != "" {
		b.SetHeader("Cookie", config.Cookie)
	}
}

// getBdstoken 获取百度网盘操作令牌
func (b *BaiduPanService) getBdstoken() (string, error) {
	if b.bdstoken != "" {
		return b.bdstoken, nil
	}

	// 调用百度网盘API获取bdstoken
	url := "https://pan.baidu.com/api/gettemplatevariable"
	data := map[string]interface{}{
		"fields": "['bdstoken']",
	}

	respData, err := b.HTTPPost(url, data, nil)
	if err != nil {
		return "", fmt.Errorf("获取bdstoken失败: %v", err)
	}

	var result struct {
		Errno int    `json:"errno"`
		Data  struct {
			Bdstoken string `json:"bdstoken"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return "", fmt.Errorf("解析bdstoken响应失败: %v", err)
	}

	if result.Errno != 0 {
		return "", fmt.Errorf("获取bdstoken失败: %d", result.Errno)
	}

	b.bdstoken = result.Data.Bdstoken
	return b.bdstoken, nil
}

// verifyPassCode 验证提取码
func (b *BaiduPanService) verifyPassCode(shareURL, passCode string) (string, error) {
	data := map[string]interface{}{
		"surl":       extractSurl(shareURL),
		"pwd":        passCode,
		"pwdid":      "",
		"channel":    "chunlei",
		"web":        "1",
		"app_id":     "250528",
		"clienttype": "0",
	}

	respData, err := b.HTTPPost("https://pan.baidu.com/share/verify", data, nil)
	if err != nil {
		return "0", fmt.Errorf("验证提取码失败: %v", err)
	}

	var result struct {
		Errno int `json:"errno"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return "0", fmt.Errorf("解析验证码响应失败: %v", err)
	}

	if result.Errno != 0 {
		return fmt.Sprintf("%d", result.Errno), nil
	}

	// 返回randsk
	var response map[string]interface{}
	if err := json.Unmarshal(respData, &response); err != nil {
		return "0", fmt.Errorf("解析randsk失败: %v", err)
	}

	if randsk, ok := response["randsk"].(string); ok {
		return randsk, nil
	}

	return "0", fmt.Errorf("未找到randsk")
}

// getTransferParams 获取转存参数
func (b *BaiduPanService) getTransferParams(shareURL string) ([]string, error) {
	surl := extractSurl(shareURL)
	url := fmt.Sprintf("https://pan.baidu.com/s/tbnet?&surl=%s&root=1", surl)

	respData, err := b.HTTPGet(url, map[string]string{
		"channel":     "chunlei",
		"web":         "1",
		"app_id":      "250528",
		"clienttype":  "0",
		"shareid":     "0",
		"uk":          "0",
		"fid":         "0",
		"page":        "1",
		"num":         "100",
		"by":          "filename",
		"order":       "asc",
		"fn":          "1",
		"fsid":        "0",
		"ct":          "0",
		"cs":          "0",
		"interation":  "1",
		"_":           fmt.Sprintf("%d", time.Now().UnixNano()/1000000),
	})
	if err != nil {
		return nil, fmt.Errorf("获取转存参数失败: %v", err)
	}

	var result struct {
		Errno int `json:"errno"`
		List  []struct {
			ShareID string `json:"shareid"`
			UserID  string `json:"uk"`
			FsID    string `json:"fs_id"`
			FileName string `json:"server_filename"`
			IsDir   int    `json:"isdir"`
		} `json:"list"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("解析转存参数响应失败: %v", err)
	}

	if result.Errno != 0 {
		return nil, fmt.Errorf("获取转存参数失败: %d", result.Errno)
	}

	if len(result.List) == 0 {
		return nil, fmt.Errorf("分享列表为空")
	}

	// 提取参数
	shareID := result.List[0].ShareID
	userID := result.List[0].UserID
	fsIds := make([]string, 0)
	fileNames := make([]string, 0)
	isDirs := make([]string, 0)

	for _, item := range result.List {
		fsIds = append(fsIds, item.FsID)
		fileNames = append(fileNames, item.FileName)
		isDirs = append(isDirs, strconv.Itoa(item.IsDir))
	}

	return []string{shareID, userID, strings.Join(fsIds, ","), strings.Join(fileNames, ","), strings.Join(isDirs, ",")}, nil
}

// createDir 创建目录
func (b *BaiduPanService) createDir(dirPath string) error {
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"a":          "commit",
		"channel":    "chunlei",
		"web":        "1",
		"app_id":     "250528",
		"clienttype": "0",
		"bdstoken":   bdstoken,
		"path":       dirPath,
		"isdir":      "1",
		"rtype":      "1",
	}

	respData, err := b.HTTPPost("https://pan.baidu.com/api/create", data, nil)
	if err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	var result struct {
		Errno int `json:"errno"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return fmt.Errorf("解析创建目录响应失败: %v", err)
	}

	if result.Errno != 0 {
		return fmt.Errorf("创建目录失败: %d", result.Errno)
	}

	return nil
}

// transferFile 转存文件
func (b *BaiduPanService) transferFile(params []string, targetPath string) error {
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"shareid":    params[0], // shareID
		"from":       params[1], // userID
		"bdstoken":   bdstoken,
		"channel":    "chunlei",
		"web":        "1",
		"app_id":     "250528",
		"clienttype": "0",
		"fsidlist":   fmt.Sprintf("[%s]", params[2]), // fsIds
		"to":         targetPath,
	}

	respData, err := b.HTTPPost("https://pan.baidu.com/api/transfer", data, nil)
	if err != nil {
		return fmt.Errorf("转存文件失败: %v", err)
	}

	var result struct {
		Errno int `json:"errno"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return fmt.Errorf("解析转存文件响应失败: %v", err)
	}

	if result.Errno != 0 {
		return fmt.Errorf("转存文件失败: %d", result.Errno)
	}

	return nil
}

// createShare 创建分享
func (b *BaiduPanService) createShare(fsIds string, expiry int, password string) (string, error) {
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"channel":    "chunlei",
		"web":        "1",
		"app_id":     "250528",
		"clienttype": "0",
		"bdstoken":   bdstoken,
		"schannel":   "0",
		"fid_list":   fmt.Sprintf("[%s]", fsIds),
		"pwd":        password,
	}

	if expiry > 0 {
		data["period"] = strconv.Itoa(expiry)
	} else {
		data["period"] = "0" // 永久
	}

	respData, err := b.HTTPPost("https://pan.baidu.com/share/set", data, nil)
	if err != nil {
		return "0", fmt.Errorf("创建分享失败: %v", err)
	}

	var result struct {
		Errno int `json:"errno"`
		Data  struct {
			Link string `json:"link"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return "0", fmt.Errorf("解析创建分享响应失败: %v", err)
	}

	if result.Errno != 0 {
		return fmt.Sprintf("%d", result.Errno), nil
	}

	return result.Data.Link, nil
}

// batchDeleteFiles 批量删除文件
func (b *BaiduPanService) batchDeleteFiles(filePaths []string) (map[string]interface{}, error) {
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return nil, err
	}

	// 构建filelist参数
	filelist := ""
	for i, path := range filePaths {
		if i > 0 {
			filelist += ","
		}
		filelist += fmt.Sprintf(`"%s"`, path)
	}

	data := map[string]interface{}{
		"channel":    "chunlei",
		"web":        "1",
		"app_id":     "250528",
		"clienttype": "0",
		"bdstoken":   bdstoken,
		"async":      "1",
		"filelist":   fmt.Sprintf("[%s]", filelist),
		"ondup":      "overwrite",
	}

	respData, err := b.HTTPPost("https://pan.baidu.com/api/filemanager", data, nil)
	if err != nil {
		return nil, fmt.Errorf("批量删除文件失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("解析批量删除文件响应失败: %v", err)
	}

	return result, nil
}

// containsAdKeywords 检查文件名是否包含广告关键词
func (b *BaiduPanService) containsAdKeywords(filename string) bool {
	// 获取广告关键词配置
	adKeywordsStr, err := b.getSystemConfigValue(entity.ConfigKeyAdKeywords)
	if err != nil {
		log.Printf("获取广告关键词配置失败: %v", err)
		return false
	}

	if adKeywordsStr == "" {
		return false
	}

	// 按逗号分割关键词
	adKeywords := b.splitKeywords(adKeywordsStr)
	return b.checkKeywordsInFilename(filename, adKeywords)
}

// getSystemConfigValue 获取系统配置值
func (b *BaiduPanService) getSystemConfigValue(key string) (string, error) {
	systemConfigRepo := repo.NewSystemConfigRepository(db.DB)
	return systemConfigRepo.GetConfigValue(key)
}

// splitKeywords 按逗号分割关键词
func (b *BaiduPanService) splitKeywords(keywordsStr string) []string {
	if keywordsStr == "" {
		return []string{}
	}

	re := regexp.MustCompile(`[,，]`)
	parts := re.Split(keywordsStr, -1)

	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// checkKeywordsInFilename 检查文件名是否包含指定关键词
func (b *BaiduPanService) checkKeywordsInFilename(filename string, keywords []string) bool {
	lowercaseFilename := strings.ToLower(filename)

	for _, keyword := range keywords {
		if strings.Contains(lowercaseFilename, strings.ToLower(keyword)) {
			log.Printf("文件 %s 包含广告关键词: %s", filename, keyword)
			return true
		}
	}

	return false
}

// extractSurl 从分享链接中提取surl
func extractSurl(shareURL string) string {
	re := regexp.MustCompile(`/s/([a-zA-Z0-9_-]+)`)
	matches := re.FindStringSubmatch(shareURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// updateBdclnd 更新cookie中的bdclnd
func (b *BaiduPanService) updateBdclnd(randsk string) {
	// 这里需要更新cookie，实现方式依赖于BasePanService的cookie管理
	// 暂时简化处理
	log.Printf("更新bdclnd: %s", randsk)
}

// getStorageFolder 获取存储文件夹名称
func (b *BaiduPanService) getStorageFolder() string {
	b.configMutex.RLock()
	config := b.config
	b.configMutex.RUnlock()

	// 根据配置返回不同的文件夹名称
	if config.ExpiredType == 2 {
		// 临时资源路径
		return "心悦临时转存文件"
	}
	return "心悦转存文件"
}

// ensureDirectoryExists 确保目录存在，不存在则创建
func (b *BaiduPanService) ensureDirectoryExists(folderName string) error {
	// 检查目录是否存在
	_, err := b.getDirList("/" + folderName)
	if err == nil {
		return nil // 目录已存在
	}

	// 创建目录
	return b.createDir("/" + folderName)
}

// getDirList 获取目录文件列表
func (b *BaiduPanService) getDirList(dirPath string) ([]map[string]interface{}, error) {
	bdstoken, err := b.getBdstoken()
	if err != nil {
		return nil, err
	}

	queryParams := map[string]string{
		"channel":    "chunlei",
		"web":        "1",
		"app_id":     "250528",
		"clienttype": "0",
		"bdstoken":   bdstoken,
		"dir":        dirPath,
		"num":        "100",
		"order":      "name",
		"desc":       "0",
	}

	respData, err := b.HTTPGet("https://pan.baidu.com/api/list", queryParams)
	if err != nil {
		return nil, fmt.Errorf("获取目录列表失败: %v", err)
	}

	var result struct {
		Errno int `json:"errno"`
		List  []map[string]interface{} `json:"list"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("解析目录列表响应失败: %v", err)
	}

	if result.Errno != 0 {
		return nil, fmt.Errorf("获取目录列表失败: %d", result.Errno)
	}

	return result.List, nil
}

// processTransferredFiles 处理转存后的文件，包括广告过滤
func (b *BaiduPanService) processTransferredFiles(dirList []map[string]interface{}, fileNames []string, folderName string) ([]string, []string, error) {
	targetFiles := []map[string]interface{}{}
	fsIdList := []string{}
	filePaths := []string{}
	adFilePaths := []string{}
	allFilesAreAds := true

	// 找到转存的文件
	for _, file := range dirList {
		if fileName, ok := file["server_filename"].(string); ok {
			// 检查是否是目标文件
			for _, targetName := range fileNames {
				if fileName == targetName {
					targetFiles = append(targetFiles, file)
					if fsID, ok := file["fs_id"].(float64); ok {
						fsIdList = append(fsIdList, strconv.FormatInt(int64(fsID), 10))
					}
					filePath := "/" + folderName + "/" + fileName
					filePaths = append(filePaths, filePath)

					// 检查是否包含广告
					if b.containsAdKeywords(fileName) {
						adFilePaths = append(adFilePaths, filePath)
					} else {
						allFilesAreAds = false
					}

					// 如果是目录，递归检查
					if isDir, ok := file["isdir"].(float64); ok && int(isDir) == 1 {
						subFiles, err := b.getDirList(filePath)
						if err == nil {
							for _, subFile := range subFiles {
								if subFileName, ok := subFile["server_filename"].(string); ok {
									if b.containsAdKeywords(subFileName) {
										adFilePaths = append(adFilePaths, filePath+"/"+subFileName)
									} else {
										allFilesAreAds = false
									}
								}
							}
						}
					}
					break
				}
			}
		}
	}

	if len(targetFiles) == 0 {
		return nil, nil, fmt.Errorf("找不到转存的文件")
	}

	// 如果所有文件都是广告，删除所有文件
	if allFilesAreAds && len(targetFiles) > 0 {
		_, err := b.batchDeleteFiles(filePaths)
		if err == nil {
			return nil, nil, fmt.Errorf("资源内容为空或所有转存的文件都包含广告内容，已全部删除")
		}
	}

	// 删除广告文件
	if len(adFilePaths) > 0 {
		_, err := b.batchDeleteFiles(adFilePaths)
		if err == nil {
			// 更新文件列表
			newFilePaths := []string{}
			newFsIdList := []string{}
			for i, path := range filePaths {
				isAd := false
				for _, adPath := range adFilePaths {
					if path == adPath {
						isAd = true
						break
					}
				}
				if !isAd {
					newFilePaths = append(newFilePaths, path)
					if i < len(fsIdList) {
						newFsIdList = append(newFsIdList, fsIdList[i])
					}
				}
			}
			filePaths = newFilePaths
			fsIdList = newFsIdList
		}
	}

	if len(filePaths) == 0 {
		return nil, nil, fmt.Errorf("删除广告后没有剩余文件")
	}

	return filePaths, fsIdList, nil
}

// GetUserInfoByEntity 根据 entity.Cks 获取用户信息（待实现）
func (b *BaiduPanService) GetUserInfoByEntity(cks entity.Cks) (*UserInfo, error) {
	return nil, nil
}
