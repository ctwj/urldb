package pan

import (
	"fmt"

	"github.com/ctwj/urldb/db/entity"
)

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

// GetUserInfo 获取用户信息
func (u *UCService) GetUserInfo(cookie string) (*UserInfo, error) {
	// 设置Cookie
	u.SetHeader("Cookie", cookie)

	// 调用UC网盘用户信息API
	userInfoURL := "https://drive.uc.cn/api/user/info"

	resp, err := u.HTTPGet(userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 解析响应
	var result struct {
		Code int `json:"code"`
		Data struct {
			Username   string `json:"username"`
			Nickname   string `json:"nickname"`
			VipStatus  int    `json:"vip_status"`
			TotalSpace int64  `json:"total_space"`
			UsedSpace  int64  `json:"used_space"`
		} `json:"data"`
	}

	if err := u.ParseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API返回错误: %d", result.Code)
	}

	// 转换VIP状态
	vipStatus := result.Data.VipStatus > 0

	// 使用nickname或username
	username := result.Data.Nickname
	if username == "" {
		username = result.Data.Username
	}

	return &UserInfo{
		Username:    username,
		VIPStatus:   vipStatus,
		UsedSpace:   result.Data.UsedSpace,
		TotalSpace:  result.Data.TotalSpace,
		ServiceType: "uc",
	}, nil
}

// GetUserInfoByEntity 根据 entity.Cks 获取用户信息（待实现）
func (u *UCService) GetUserInfoByEntity(cks entity.Cks) (*UserInfo, error) {
	return nil, nil
}
