package pan

import (
	"fmt"
	"testing"
)

func TestQuarkServiceSingleton(t *testing.T) {
	// 测试夸克服务单例模式
	config1 := &PanConfig{
		URL:         "https://pan.quark.cn/s/123456",
		IsType:      0,
		ExpiredType: 1,
	}

	config2 := &PanConfig{
		URL:         "https://pan.quark.cn/s/789012",
		IsType:      1,
		ExpiredType: 2,
	}

	service1 := NewQuarkPanService(config1)
	service2 := NewQuarkPanService(config2)

	// 应该返回相同的实例
	if service1 != service2 {
		t.Error("QuarkPanService 应该返回相同的单例实例")
	}

	// 验证服务类型
	if service1.GetServiceType() != Quark {
		t.Error("服务类型应该是 Quark")
	}
}

func TestAlipanServiceSingleton(t *testing.T) {
	// 测试阿里云盘服务单例模式
	config1 := &PanConfig{
		URL:         "https://www.alipan.com/s/123456",
		IsType:      0,
		ExpiredType: 1,
	}

	config2 := &PanConfig{
		URL:         "https://www.alipan.com/s/789012",
		IsType:      1,
		ExpiredType: 2,
	}

	service1 := NewAlipanService(config1)
	service2 := NewAlipanService(config2)

	// 应该返回相同的实例
	if service1 != service2 {
		t.Error("AlipanService 应该返回相同的单例实例")
	}

	// 验证服务类型
	if service1.GetServiceType() != Alipan {
		t.Error("服务类型应该是 Alipan")
	}
}

func TestServiceConcurrency(t *testing.T) {
	// 测试并发情况下的服务单例
	done := make(chan bool, 10)
	quarkServices := make([]PanService, 10)
	alipanServices := make([]PanService, 10)

	// 并发创建夸克服务
	for i := 0; i < 10; i++ {
		go func(index int) {
			config := &PanConfig{
				URL:         fmt.Sprintf("https://pan.quark.cn/s/%d", index),
				IsType:      0,
				ExpiredType: 1,
			}
			service := NewQuarkPanService(config)
			quarkServices[index] = service
			done <- true
		}(i)
	}

	// 并发创建阿里云盘服务
	for i := 0; i < 10; i++ {
		go func(index int) {
			config := &PanConfig{
				URL:         fmt.Sprintf("https://www.alipan.com/s/%d", index),
				IsType:      0,
				ExpiredType: 1,
			}
			service := NewAlipanService(config)
			alipanServices[index] = service
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 20; i++ {
		<-done
	}

	// 验证夸克服务单例
	firstQuarkService := quarkServices[0]
	for i := 1; i < 10; i++ {
		if quarkServices[i] != firstQuarkService {
			t.Errorf("夸克服务并发测试失败：实例 %d 与第一个实例不同", i)
		}
	}

	// 验证阿里云盘服务单例
	firstAlipanService := alipanServices[0]
	for i := 1; i < 10; i++ {
		if alipanServices[i] != firstAlipanService {
			t.Errorf("阿里云盘服务并发测试失败：实例 %d 与第一个实例不同", i)
		}
	}
}

func BenchmarkQuarkServiceCreation(b *testing.B) {
	// 测试夸克服务创建性能
	config := &PanConfig{
		URL:         "https://pan.quark.cn/s/123456",
		IsType:      0,
		ExpiredType: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewQuarkPanService(config)
	}
}

func BenchmarkAlipanServiceCreation(b *testing.B) {
	// 测试阿里云盘服务创建性能
	config := &PanConfig{
		URL:         "https://www.alipan.com/s/123456",
		IsType:      0,
		ExpiredType: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewAlipanService(config)
	}
}

func BenchmarkFactoryGetQuarkService(b *testing.B) {
	// 测试工厂获取夸克服务性能
	factory := GetInstance()
	config := &PanConfig{
		URL:         "https://pan.quark.cn/s/123456",
		IsType:      0,
		ExpiredType: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = factory.GetQuarkService(config)
	}
}

func BenchmarkFactoryGetAlipanService(b *testing.B) {
	// 测试工厂获取阿里云盘服务性能
	factory := GetInstance()
	config := &PanConfig{
		URL:         "https://www.alipan.com/s/123456",
		IsType:      0,
		ExpiredType: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = factory.GetAlipanService(config)
	}
}
