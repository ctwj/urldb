package pan

import (
	"fmt"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

// BaiduPanService 百度网盘服务
type BaiduPanService struct {
	*BasePanService
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
		"Content-Type":    "application/json",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})

	return service
}

// GetServiceType 获取服务类型
func (b *BaiduPanService) GetServiceType() ServiceType {
	return BaiduPan
}

// Transfer 转存分享链接
func (b *BaiduPanService) Transfer(shareID string) (*TransferResult, error) {
	// TODO: 实现百度网盘转存逻辑
	return ErrorResult("百度网盘转存功能暂未实现"), nil
}

// GetFiles 获取文件列表
func (b *BaiduPanService) GetFiles(pdirFid string) (*TransferResult, error) {
	// TODO: 实现百度网盘文件列表获取
	return ErrorResult("百度网盘文件列表功能暂未实现"), nil
}

// DeleteFiles 删除文件
func (b *BaiduPanService) DeleteFiles(fileList []string) (*TransferResult, error) {
	// TODO: 实现百度网盘文件删除
	return ErrorResult("百度网盘文件删除功能暂未实现"), nil
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
