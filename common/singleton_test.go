package pan

import (
	"testing"
)

func TestSingletonPattern(t *testing.T) {
	// 测试多次调用NewPanFactory()返回相同实例
	factory1 := NewPanFactory()
	factory2 := NewPanFactory()
	factory3 := GetInstance()

	if factory1 != factory2 {
		t.Error("NewPanFactory() 应该返回相同的实例")
	}

	if factory1 != factory3 {
		t.Error("GetInstance() 应该返回相同的实例")
	}

	if factory2 != factory3 {
		t.Error("所有获取方法应该返回相同的实例")
	}

	// 验证实例不为nil
	if factory1 == nil {
		t.Error("工厂实例不应该为nil")
	}
}

func TestSingletonConcurrency(t *testing.T) {
	// 测试并发情况下的单例模式
	done := make(chan bool, 10)
	instances := make([]*PanFactory, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			instances[index] = GetInstance()
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有实例都是相同的
	firstInstance := instances[0]
	for i := 1; i < 10; i++ {
		if instances[i] != firstInstance {
			t.Errorf("并发测试失败：实例 %d 与第一个实例不同", i)
		}
	}
}

func BenchmarkSingletonCreation(b *testing.B) {
	// 测试单例创建的性能
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetInstance()
	}
}

func BenchmarkNewPanFactory(b *testing.B) {
	// 测试NewPanFactory的性能
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewPanFactory()
	}
}
