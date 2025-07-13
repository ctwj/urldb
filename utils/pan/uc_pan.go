package pan

// UCService UC网盘服务
type UCService struct {
	*BasePanService
}

// NewUCService 创建UC网盘服务
func NewUCService(config *PanConfig) *UCService {
	service := &UCService{
		BasePanService: NewBasePanService(config),
	}

	// 设置UC网盘的默认请求头
	service.SetHeaders(map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Content-Type":    "application/json",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})

	return service
}

// GetServiceType 获取服务类型
func (u *UCService) GetServiceType() ServiceType {
	return UC
}

// Transfer 转存分享链接
func (u *UCService) Transfer(shareID string) (*TransferResult, error) {
	// TODO: 实现UC网盘转存逻辑
	return ErrorResult("UC网盘转存功能暂未实现"), nil
}

// GetFiles 获取文件列表
func (u *UCService) GetFiles(pdirFid string) (*TransferResult, error) {
	// TODO: 实现UC网盘文件列表获取
	return ErrorResult("UC网盘文件列表功能暂未实现"), nil
}

// DeleteFiles 删除文件
func (u *UCService) DeleteFiles(fileList []string) (*TransferResult, error) {
	// TODO: 实现UC网盘文件删除
	return ErrorResult("UC网盘文件删除功能暂未实现"), nil
}
