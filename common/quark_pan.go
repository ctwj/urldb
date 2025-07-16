package pan

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// QuarkPanService 夸克网盘服务
type QuarkPanService struct {
	*BasePanService
	configMutex sync.RWMutex // 保护配置的读写锁
}

// 单例相关变量
var (
	quarkInstance *QuarkPanService
	quarkOnce     sync.Once
)

// NewQuarkPanService 创建夸克网盘服务（单例模式）
func NewQuarkPanService(config *PanConfig) *QuarkPanService {
	quarkOnce.Do(func() {
		quarkInstance = &QuarkPanService{
			BasePanService: NewBasePanService(config),
		}

		// 设置夸克网盘的默认请求头
		quarkInstance.SetHeaders(map[string]string{
			"Accept":             "application/json, text/plain, */*",
			"Accept-Language":    "zh-CN,zh;q=0.9",
			"Content-Type":       "application/json;charset=UTF-8",
			"Sec-Ch-Ua":          `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": `"Windows"`,
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "same-site",
			"Referer":            "https://pan.quark.cn/",
			"Referrer-Policy":    "strict-origin-when-cross-origin",
			"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		})
	})

	// 更新配置
	quarkInstance.UpdateConfig(config)

	return quarkInstance
}

// GetQuarkInstance 获取夸克网盘服务单例实例
func GetQuarkInstance() *QuarkPanService {
	return NewQuarkPanService(nil)
}

// UpdateConfig 更新配置（线程安全）
func (q *QuarkPanService) UpdateConfig(config *PanConfig) {
	if config == nil {
		return
	}

	q.configMutex.Lock()
	defer q.configMutex.Unlock()

	q.config = config
}

// GetServiceType 获取服务类型
func (q *QuarkPanService) GetServiceType() ServiceType {
	return Quark
}

// Transfer 转存分享链接
func (q *QuarkPanService) Transfer(shareID string) (*TransferResult, error) {
	// 读取配置（线程安全）
	q.configMutex.RLock()
	config := q.config
	q.configMutex.RUnlock()

	log.Printf("开始处理夸克分享: %s", shareID)

	// 获取stoken
	var stoken string
	if config.Stoken == "" {
		stokenResult, err := q.getStoken(shareID)
		if err != nil {
			return ErrorResult(fmt.Sprintf("获取stoken失败: %v", err)), nil
		}

		stoken = strings.ReplaceAll(stokenResult.Stoken, " ", "+")
	} else {
		stoken = strings.ReplaceAll(config.Stoken, " ", "+")
	}

	// 获取分享详情
	shareResult, err := q.getShare(shareID, stoken)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取分享详情失败: %v", err)), nil
	}

	if config.IsType == 1 {
		// 遍历 资源目录结构
		for _, item := range shareResult.List {
			// 获取文件信息
			fileList, err := q.getDirFile(item.Fid)
			if err != nil {
				log.Printf("获取目录文件失败: %v", err)
				continue
			}

			// 处理文件列表
			if fileList != nil {
				log.Printf("目录 %s 包含 %d 个文件/文件夹", item.Fid, len(fileList))

				// 遍历所有文件，可以在这里添加具体的处理逻辑
				for _, file := range fileList {
					if fileName, ok := file["file_name"].(string); ok {
						if fileType, ok := file["file_type"].(float64); ok {
							fileTypeStr := "文件"
							if fileType == 1 {
								fileTypeStr = "目录"
							}
							log.Printf("  - %s (%s)", fileName, fileTypeStr)
						}
					}
				}
			}
		}

		// 直接返回资源信息
		return SuccessResult("检验成功", map[string]interface{}{
			"title":    shareResult.Share.Title,
			"shareUrl": config.URL,
		}), nil
	}

	// 提取文件信息
	fidList := make([]string, 0)
	fidTokenList := make([]string, 0)
	title := shareResult.Share.Title

	for _, item := range shareResult.List {
		fidList = append(fidList, item.Fid)
		fidTokenList = append(fidTokenList, item.ShareFidToken)
	}

	// 转存资源
	saveResult, err := q.getShareSave(shareID, stoken, fidList, fidTokenList)
	if err != nil {
		return ErrorResult(fmt.Sprintf("转存失败: %v", err)), nil
	}

	taskID := saveResult.TaskID

	// 等待转存完成
	myData, err := q.waitForTask(taskID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("等待转存完成失败: %v", err)), nil
	}

	// 删除广告文件（如果有配置）
	if err := q.deleteAdFiles(myData.SaveAs.SaveAsTopFids[0]); err != nil {
		log.Printf("删除广告文件失败: %v", err)
	}

	// 分享资源
	shareBtnResult, err := q.getShareBtn(myData.SaveAs.SaveAsTopFids, title)
	if err != nil {
		return ErrorResult(fmt.Sprintf("分享失败: %v", err)), nil
	}

	// 等待分享完成
	shareTaskResult, err := q.waitForTask(shareBtnResult.TaskID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("等待分享完成失败: %v", err)), nil
	}

	// 获取分享密码
	passwordResult, err := q.getSharePassword(shareTaskResult.ShareID)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取分享密码失败: %v", err)), nil
	}

	// 确定fid
	var fid string
	if len(myData.SaveAs.SaveAsTopFids) > 1 {
		fid = strings.Join(myData.SaveAs.SaveAsTopFids, ",")
	} else {
		fid = passwordResult.FirstFile.Fid
	}

	return SuccessResult("转存成功", map[string]interface{}{
		"shareUrl": passwordResult.ShareURL,
		"title":    passwordResult.ShareTitle,
		"fid":      fid,
		"code":     passwordResult.Code,
	}), nil
}

// GetFiles 获取文件列表
func (q *QuarkPanService) GetFiles(pdirFid string) (*TransferResult, error) {
	if pdirFid == "" {
		pdirFid = "0"
	}

	queryParams := map[string]string{
		"pr":              "ucpro",
		"fr":              "pc",
		"uc_param_str":    "",
		"pdir_fid":        pdirFid,
		"_page":           "1",
		"_size":           "50",
		"_fetch_total":    "1",
		"_fetch_sub_dirs": "0",
		"_sort":           "file_type:asc,updated_at:desc",
	}

	data, err := q.HTTPGet("https://drive-pc.quark.cn/1/clouddrive/file/sort", queryParams)
	if err != nil {
		return ErrorResult(fmt.Sprintf("获取文件列表失败: %v", err)), nil
	}

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			List []interface{} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return ErrorResult("解析响应失败"), nil
	}

	if response.Status != 200 {
		message := response.Message
		if message == "require login [guest]" {
			message = "夸克未登录，请检查cookie"
		}
		return ErrorResult(message), nil
	}

	return SuccessResult("获取成功", response.Data.List), nil
}

// DeleteFiles 删除文件
func (q *QuarkPanService) DeleteFiles(fileList []string) (*TransferResult, error) {
	if len(fileList) == 0 {
		return ErrorResult("文件列表为空"), nil
	}

	data := map[string]interface{}{
		"fid_list": fileList,
	}

	queryParams := map[string]string{
		"pr":           "ucpro",
		"fr":           "pc",
		"uc_param_str": "",
	}

	_, err := q.HTTPPost("https://drive-pc.quark.cn/1/clouddrive/file/delete", data, queryParams)
	if err != nil {
		return ErrorResult(fmt.Sprintf("删除文件失败: %v", err)), nil
	}

	return SuccessResult("删除成功", nil), nil
}

// getStoken 获取stoken
func (q *QuarkPanService) getStoken(shareID string) (*StokenResult, error) {
	data := map[string]interface{}{
		"passcode": "",
		"pwd_id":   shareID,
	}

	queryParams := map[string]string{
		"pr":           "ucpro",
		"fr":           "pc",
		"uc_param_str": "",
	}

	respData, err := q.HTTPPost("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/token", data, queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int          `json:"status"`
		Message string       `json:"message"`
		Data    StokenResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf(response.Message)
	}

	return &response.Data, nil
}

// getShare 获取分享详情
func (q *QuarkPanService) getShare(shareID, stoken string) (*ShareResult, error) {
	queryParams := map[string]string{
		"pr":            "ucpro",
		"fr":            "pc",
		"uc_param_str":  "",
		"pwd_id":        shareID,
		"stoken":        stoken,
		"pdir_fid":      "0",
		"force":         "0",
		"_page":         "1",
		"_size":         "100",
		"_fetch_banner": "1",
		"_fetch_share":  "1",
		"_fetch_total":  "1",
		"_sort":         "file_type:asc,updated_at:desc",
	}

	respData, err := q.HTTPGet("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/detail", queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    ShareResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf(response.Message)
	}

	return &response.Data, nil
}

// getShareSave 转存分享
func (q *QuarkPanService) getShareSave(shareID, stoken string, fidList, fidTokenList []string) (*SaveResult, error) {
	data := map[string]interface{}{
		"pwd_id":         shareID,
		"stoken":         stoken,
		"fid_list":       fidList,
		"fid_token_list": fidTokenList,
		"to_pdir_fid":    "0", // 默认存储到根目录
	}

	queryParams := map[string]string{
		"pr":           "ucpro",
		"fr":           "pc",
		"uc_param_str": "",
	}

	respData, err := q.HTTPPost("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/save", data, queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int        `json:"status"`
		Message string     `json:"message"`
		Data    SaveResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf(response.Message)
	}

	return &response.Data, nil
}

// getShareBtn 分享按钮
func (q *QuarkPanService) getShareBtn(fidList []string, title string) (*ShareBtnResult, error) {
	data := map[string]interface{}{
		"fid_list":     fidList,
		"title":        title,
		"expired_type": 1, // 永久分享
	}

	queryParams := map[string]string{
		"pr":           "ucpro",
		"fr":           "pc",
		"uc_param_str": "",
	}

	respData, err := q.HTTPPost("https://drive-pc.quark.cn/1/clouddrive/share/create", data, queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int            `json:"status"`
		Message string         `json:"message"`
		Data    ShareBtnResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf(response.Message)
	}

	return &response.Data, nil
}

// getShareTask 获取分享任务状态
func (q *QuarkPanService) getShareTask(taskID string, retryIndex int) (*TaskResult, error) {
	queryParams := map[string]string{
		"pr":           "ucpro",
		"fr":           "pc",
		"uc_param_str": "",
		"task_id":      taskID,
		"retry_index":  fmt.Sprintf("%d", retryIndex),
	}

	respData, err := q.HTTPGet("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/task", queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int        `json:"status"`
		Message string     `json:"message"`
		Data    TaskResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf(response.Message)
	}

	return &response.Data, nil
}

// getSharePassword 获取分享密码
func (q *QuarkPanService) getSharePassword(shareID string) (*PasswordResult, error) {
	queryParams := map[string]string{
		"pr":           "ucpro",
		"fr":           "pc",
		"uc_param_str": "",
		"share_id":     shareID,
	}

	respData, err := q.HTTPGet("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/password", queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Status  int            `json:"status"`
		Message string         `json:"message"`
		Data    PasswordResult `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf(response.Message)
	}

	return &response.Data, nil
}

// waitForTask 等待任务完成
func (q *QuarkPanService) waitForTask(taskID string) (*TaskResult, error) {
	maxRetries := 50
	retryDelay := 2 * time.Second

	for retryIndex := 0; retryIndex < maxRetries; retryIndex++ {
		result, err := q.getShareTask(taskID, retryIndex)
		if err != nil {
			if strings.Contains(err.Error(), "capacity limit[{0}]") {
				return nil, fmt.Errorf("容量不足")
			}
			return nil, err
		}

		if result.Status == 2 { // 任务完成
			return result, nil
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("任务超时")
}

// deleteAdFiles 删除广告文件
func (q *QuarkPanService) deleteAdFiles(pdirFid string) error {
	// 这里可以添加广告文件删除逻辑
	// 需要从配置中读取禁止的关键词列表
	return nil
}

// getDirFile 获取指定文件夹的文件列表
func (q *QuarkPanService) getDirFile(pdirFid string) ([]map[string]interface{}, error) {
	log.Printf("正在遍历父文件夹: %s", pdirFid)

	queryParams := map[string]string{
		"pr":              "ucpro",
		"fr":              "pc",
		"uc_param_str":    "",
		"pdir_fid":        pdirFid,
		"_page":           "1",
		"_size":           "50",
		"_fetch_total":    "1",
		"_fetch_sub_dirs": "0",
		"_sort":           "updated_at:desc",
	}

	respData, err := q.HTTPGet("https://drive-pc.quark.cn/1/clouddrive/file/sort", queryParams)
	if err != nil {
		log.Printf("获取目录文件失败: %v", err)
		return nil, err
	}

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			List []map[string]interface{} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		log.Printf("解析目录文件响应失败: %v", err)
		return nil, err
	}

	if response.Status != 200 {
		return nil, fmt.Errorf(response.Message)
	}

	// 递归处理子目录
	var allFiles []map[string]interface{}
	for _, item := range response.Data.List {
		// 添加当前文件/目录
		allFiles = append(allFiles, item)

		// 如果是目录，递归获取子目录内容
		if fileType, ok := item["file_type"].(float64); ok && fileType == 1 { // 1表示目录
			if fid, ok := item["fid"].(string); ok {
				subFiles, err := q.getDirFile(fid)
				if err != nil {
					log.Printf("获取子目录 %s 失败: %v", fid, err)
					continue
				}
				allFiles = append(allFiles, subFiles...)
			}
		}
	}

	return allFiles, nil
}

// 定义各种结果结构体
type StokenResult struct {
	Stoken string `json:"stoken"`
	Title  string `json:"title"`
}

type ShareResult struct {
	Share struct {
		Title string `json:"title"`
	} `json:"share"`
	List []struct {
		Fid           string `json:"fid"`
		ShareFidToken string `json:"share_fid_token"`
	} `json:"list"`
}

type SaveResult struct {
	TaskID string `json:"task_id"`
}

type ShareBtnResult struct {
	TaskID string `json:"task_id"`
}

type TaskResult struct {
	Status  int    `json:"status"`
	ShareID string `json:"share_id"`
	SaveAs  struct {
		SaveAsTopFids []string `json:"save_as_top_fids"`
	} `json:"save_as"`
}

type PasswordResult struct {
	ShareURL   string `json:"share_url"`
	ShareTitle string `json:"share_title"`
	Code       string `json:"code"`
	FirstFile  struct {
		Fid string `json:"fid"`
	} `json:"first_file"`
}

// GetUserInfo 获取用户信息
func (q *QuarkPanService) GetUserInfo(cookie string) (*UserInfo, error) {
	// 临时设置cookie
	originalCookie := q.GetHeader("Cookie")
	q.SetHeader("Cookie", cookie)
	defer q.SetHeader("Cookie", originalCookie) // 恢复原始cookie

	// 获取用户基本信息
	queryParams := map[string]string{
		"platform": "pc",
		"fr":       "pc",
	}

	data, err := q.HTTPGet("https://pan.quark.cn/account/info", queryParams)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	var response struct {
		Success bool   `json:"success"`
		Code    string `json:"code"`
		Data    struct {
			Nickname  string   `json:"nickname"`
			AvatarUri string   `json:"avatarUri"`
			Mobilekps string   `json:"mobilekps"`
			Config    struct{} `json:"config"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	if !response.Success || response.Code != "OK" {
		return nil, fmt.Errorf("获取用户信息失败: API返回错误")
	}

	// 获取用户详细信息（容量和会员信息）
	queryParams1 := map[string]string{
		"pr":              "ucpro",
		"fr":              "pc",
		"uc_param_str":    "",
		"fetch_subscribe": "true",
		"_ch":             "home",
		"fetch_identity":  "true",
	}
	data1, err := q.HTTPGet("https://drive-pc.quark.cn/1/clouddrive/member", queryParams1)
	if err != nil {
		return nil, fmt.Errorf("获取用户详细信息失败: %v", err)
	}

	var memberResponse struct {
		Status  int    `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TotalCapacity     int64  `json:"secret_total_capacity"`
			SecretUseCapacity int64  `json:"secret_use_capacity"`
			MemberType        string `json:"member_type"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data1, &memberResponse); err != nil {
		return nil, fmt.Errorf("解析用户详细信息失败: %v", err)
	}

	if memberResponse.Status != 200 || memberResponse.Code != 0 {
		return nil, fmt.Errorf("获取用户详细信息失败: %s", memberResponse.Message)
	}

	// 判断VIP状态
	vipStatus := memberResponse.Data.MemberType != "NORMAL"

	return &UserInfo{
		Username:    response.Data.Nickname,
		VIPStatus:   vipStatus,
		UsedSpace:   memberResponse.Data.SecretUseCapacity,
		TotalSpace:  memberResponse.Data.TotalCapacity,
		ServiceType: "quark",
	}, nil
}

// formatBytes 格式化字节数为可读格式
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
