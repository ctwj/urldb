package pan

import (
	"testing"
)

// TestQuarkGetUserInfo 测试夸克网盘GetUserInfo功能
func TestQuarkGetUserInfo(t *testing.T) {
	// 创建夸克网盘服务实例
	service := NewQuarkPanService(&PanConfig{})

	// 测试无效cookie
	_, err := service.GetUserInfo("invalid_cookie")
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	} else {
		t.Logf("正确返回错误: %v", err)
	}
}

// TestQuarkGetUserInfoWithValidCookie 测试有效cookie的情况
func TestQuarkGetUserInfoWithValidCookie(t *testing.T) {
	// 创建夸克网盘服务实例
	service := NewQuarkPanService(&PanConfig{})

	// 使用测试cookie（需要替换为有效的cookie）
	testCookie := "your_test_cookie_here"
	if testCookie == "your_test_cookie_here" {
		t.Skip("跳过测试，需要提供有效的cookie")
	}

	userInfo, err := service.GetUserInfo(testCookie)
	if err != nil {
		t.Logf("获取用户信息失败: %v", err)
		return
	}

	// 验证返回的用户信息
	if userInfo == nil {
		t.Fatal("用户信息为空")
	}

	t.Logf("用户名: %s", userInfo.Username)
	t.Logf("VIP状态: %t", userInfo.VIPStatus)
	t.Logf("容量信息: %s", userInfo.Capacity)
	t.Logf("服务类型: %s", userInfo.ServiceType)

	// 基本验证
	if userInfo.ServiceType != "quark" {
		t.Errorf("服务类型不匹配，期望: quark, 实际: %s", userInfo.ServiceType)
	}
}
