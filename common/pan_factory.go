package pan

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// ServiceType 定义网盘服务类型
type ServiceType int

const (
	Quark ServiceType = iota
	Alipan
	BaiduPan
	UC
	NotFound
	Xunlei
	Tianyi
	Pan123
	Pan115
)

// String 返回服务类型的字符串表示
func (s ServiceType) String() string {
	switch s {
	case Quark:
		return "quark"
	case Alipan:
		return "alipan"
	case BaiduPan:
		return "baidu"
	case UC:
		return "uc"
	case Xunlei:
		return "xunlei"
	case Tianyi:
		return "tianyi"
	case Pan123:
		return "123pan"
	case Pan115:
		return "115"
	default:
		return "unknown"
	}
}

// PanConfig 网盘配置
type PanConfig struct {
	URL         string `json:"url"`
	Code        string `json:"code"`
	IsType      int    `json:"isType"`      // 0: 转存并分享后的资源信息, 1: 直接获取资源信息
	ExpiredType int    `json:"expiredType"` // 1: 分享永久, 2: 临时
	AdFid       string `json:"adFid"`       // 夸克专用 - 分享时带上这个文件的fid
	Stoken      string `json:"stoken"`
	Cookie      string `json:"cookie"`
}

// TransferResult 转存结果
type TransferResult struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
	ShareURL string      `json:"shareUrl,omitempty"`
	Title    string      `json:"title,omitempty"`
	Fid      string      `json:"fid,omitempty"`
}

// UserInfo 用户信息结构体
type UserInfo struct {
	Username    string `json:"username"`    // 用户名
	VIPStatus   bool   `json:"vipStatus"`   // VIP状态
	UsedSpace   int64  `json:"usedSpace"`   // 已使用空间
	TotalSpace  int64  `json:"totalSpace"`  // 总空间
	ServiceType string `json:"serviceType"` // 服务类型
}

// PanService 网盘服务接口
type PanService interface {
	// Transfer 转存分享链接
	Transfer(shareID string) (*TransferResult, error)

	// GetFiles 获取文件列表
	GetFiles(pdirFid string) (*TransferResult, error)

	// DeleteFiles 删除文件
	DeleteFiles(fileList []string) (*TransferResult, error)

	// GetUserInfo 获取用户信息
	GetUserInfo(cookie string) (*UserInfo, error)

	// GetServiceType 获取服务类型
	GetServiceType() ServiceType
}

// PanFactory 网盘工厂
type PanFactory struct{}

// 单例相关变量
var (
	instance *PanFactory
	once     sync.Once
)

// NewPanFactory 创建网盘工厂实例（单例模式）
func NewPanFactory() *PanFactory {
	once.Do(func() {
		instance = &PanFactory{}
	})
	return instance
}

// GetInstance 获取单例实例（推荐使用）
func GetInstance() *PanFactory {
	return NewPanFactory()
}

// CreatePanService 根据URL创建对应的网盘服务
func (f *PanFactory) CreatePanService(url string, config *PanConfig) (PanService, error) {
	serviceType := ExtractServiceType(url)

	switch serviceType {
	case Quark:
		return NewQuarkPanService(config), nil
	case Alipan:
		return NewAlipanService(config), nil
	case BaiduPan:
		return NewBaiduPanService(config), nil
	case UC:
		return NewUCService(config), nil
	default:
		return nil, fmt.Errorf("不支持的服务类型: %s", url)
	}
}

// CreatePanServiceByType 根据服务类型创建对应的网盘服务
func (f *PanFactory) CreatePanServiceByType(serviceType ServiceType, config *PanConfig) (PanService, error) {
	switch serviceType {
	case Quark:
		return NewQuarkPanService(config), nil
	case Alipan:
		return NewAlipanService(config), nil
	case BaiduPan:
		return NewBaiduPanService(config), nil
	case UC:
		return NewUCService(config), nil
	// case Xunlei:
	// 	return NewXunleiService(config), nil
	// case Tianyi:
	// 	return NewTianyiService(config), nil
	default:
		return nil, fmt.Errorf("不支持的服务类型: %d", serviceType)
	}
}

// GetQuarkService 获取夸克网盘服务单例
func (f *PanFactory) GetQuarkService(config *PanConfig) PanService {
	service := NewQuarkPanService(config)
	return service
}

// GetAlipanService 获取阿里云盘服务单例
func (f *PanFactory) GetAlipanService(config *PanConfig) PanService {
	service := NewAlipanService(config)
	return service
}

// GetBaiduService 获取百度网盘服务单例
func (f *PanFactory) GetBaiduService(config *PanConfig) PanService {
	service := NewBaiduPanService(config)
	return service
}

// GetUCService 获取UC网盘服务单例
func (f *PanFactory) GetUCService(config *PanConfig) PanService {
	service := NewUCService(config)
	return service
}

// ExtractServiceType 从URL中提取服务类型
func ExtractServiceType(url string) ServiceType {
	url = strings.ToLower(url)

	// "https://www.123pan.com/s/i4uaTd-WHn0", // 公开分享
	// "https://www.123912.com/s/U8f2Td-ZeOX",
	// "https://www.123684.coms/u9izjv-k3uWv",
	// "https://www.123pan.com/s/A6cA-AKH11", // 外链不存在

	patterns := map[string]ServiceType{
		"pan.quark.cn":        Quark,
		"www.alipan.com":      Alipan,
		"www.aliyundrive.com": Alipan,
		"pan.baidu.com":       BaiduPan,
		"drive.uc.cn":         UC,
		"fast.uc.cn":          UC,
		"pan.xunlei.com":      Xunlei,
		"cloud.189.cn":        Tianyi,
		"www.123pan.com":      Pan123,
		"www.123912.com":      Pan123,
		"www.123684.com":      Pan123,
		"115cdn.com":          Pan115,
		"anxia.com":           Pan115,
		"115.com/":            Pan115,
	}

	for pattern, serviceType := range patterns {
		if strings.Contains(url, pattern) {
			return serviceType
		}
	}

	return NotFound
}

// ExtractShareId 从URL中提取分享ID
func ExtractShareId(url string) (string, ServiceType) {
	// 处理entry参数
	if strings.Contains(url, "?entry=") {
		url = strings.Split(url, "?entry=")[0]
	}

	// 提取分享ID
	shareID := ""
	substring := strings.Index(url, "/s/")
	if substring == -1 {
		substring = strings.Index(url, "/t/") // 天翼云 是 t
		shareID = url[substring+3:]
	}
	if substring == -1 {
		substring = strings.Index(url, "/web/share?code=") // 天翼云 带密码
		shareID = url[substring+11:]
	}
	if substring == -1 {
		substring = strings.Index(url, "/p/") // 天翼云 是 p
		shareID = url[substring+3:]
	}
	if substring == -1 {
		return "", NotFound
	}

	// 去除可能的锚点
	if hashIndex := strings.Index(shareID, "#"); hashIndex != -1 {
		shareID = shareID[:hashIndex]
	}

	serviceType := ExtractServiceType(url)
	return shareID, serviceType
}

// SuccessResult 创建成功结果
func SuccessResult(message string, data interface{}) *TransferResult {
	return &TransferResult{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResult 创建错误结果
func ErrorResult(message string) *TransferResult {
	return &TransferResult{
		Success: false,
		Message: message,
	}
}

// ParseCapacityString 解析容量字符串为字节数
func ParseCapacityString(capacityStr string) int64 {
	if capacityStr == "" {
		return 0
	}

	// 移除空格并转换为小写
	capacityStr = strings.TrimSpace(strings.ToLower(capacityStr))

	var multiplier int64 = 1
	if strings.Contains(capacityStr, "gb") {
		multiplier = 1024 * 1024 * 1024
		capacityStr = strings.Replace(capacityStr, "gb", "", -1)
	} else if strings.Contains(capacityStr, "mb") {
		multiplier = 1024 * 1024
		capacityStr = strings.Replace(capacityStr, "mb", "", -1)
	} else if strings.Contains(capacityStr, "kb") {
		multiplier = 1024
		capacityStr = strings.Replace(capacityStr, "kb", "", -1)
	} else if strings.Contains(capacityStr, "b") {
		capacityStr = strings.Replace(capacityStr, "b", "", -1)
	}

	// 解析数字
	capacityStr = strings.TrimSpace(capacityStr)
	if capacityStr == "" {
		return 0
	}

	// 尝试解析浮点数
	if strings.Contains(capacityStr, ".") {
		if val, err := strconv.ParseFloat(capacityStr, 64); err == nil {
			return int64(val * float64(multiplier))
		}
	} else {
		if val, err := strconv.ParseInt(capacityStr, 10, 64); err == nil {
			return val * multiplier
		}
	}

	return 0
}
