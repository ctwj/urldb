package utils

import (
	"log"
)

// PanExample 网盘工厂模式使用示例
func PanExample() {
	// 创建工厂实例
	factory := NewPanFactory()

	// 示例URL
	urls := []string{
		"https://pan.quark.cn/s/123456789",
		"https://www.alipan.com/s/abcdef123",
		"https://pan.baidu.com/s/xyz789",
		"https://drive.uc.cn/s/uc123456",
	}

	for _, url := range urls {
		log.Printf("处理URL: %s", url)

		// 创建配置
		config := &PanConfig{
			URL:         url,
			Code:        "",
			IsType:      0, // 转存并分享后的资源信息
			ExpiredType: 1, // 永久分享
			AdFid:       "",
			Stoken:      "",
		}

		// 使用工厂创建对应的网盘服务
		panService, err := factory.CreatePanService(url, config)
		if err != nil {
			log.Printf("创建网盘服务失败: %v", err)
			continue
		}

		// 获取服务类型
		serviceType := panService.GetServiceType()
		log.Printf("检测到服务类型: %s", serviceType.String())

		// 提取分享ID
		shareID, extractedType := ExtractShareId(url)
		if extractedType == NotFound {
			log.Printf("不支持的链接格式: %s", url)
			continue
		}

		log.Printf("分享ID: %s", shareID)

		// 根据服务类型进行不同处理
		switch serviceType {
		case Quark:
			log.Println("使用夸克网盘服务")
			// 这里可以调用 panService.Transfer(shareID)

		case Alipan:
			log.Println("使用阿里云盘服务")
			// 这里可以调用 panService.Transfer(shareID)

		case BaiduPan:
			log.Println("使用百度网盘服务")
			// 这里可以调用 panService.Transfer(shareID)

		case UC:
			log.Println("使用UC网盘服务")
			// 这里可以调用 panService.Transfer(shareID)

		default:
			log.Printf("暂不支持的服务类型: %s", serviceType.String())
		}

		log.Println("---")
	}
}

// TransferExample 转存示例
func TransferExample(url string) (*TransferResult, error) {
	factory := NewPanFactory()

	config := &PanConfig{
		URL:         url,
		Code:        "",
		IsType:      0,
		ExpiredType: 1,
		AdFid:       "",
		Stoken:      "",
	}

	// 创建服务
	panService, err := factory.CreatePanService(url, config)
	if err != nil {
		return ErrorResult("创建网盘服务失败"), err
	}

	// 提取分享ID
	shareID, serviceType := ExtractShareId(url)
	if serviceType == NotFound {
		return ErrorResult("不支持的链接格式"), nil
	}

	// 执行转存
	result, err := panService.Transfer(shareID)
	if err != nil {
		return ErrorResult("转存失败"), err
	}

	return result, nil
}
