package middleware

import (
	"errors"
	"strings"
	"testing"
)

// stub generator：可断言是否被调用
func stubGenerator(returnValue string, err error) (func() (string, error), *int) {
	calls := 0
	return func() (string, error) {
		calls++
		return returnValue, err
	}, &calls
}

func TestResolveSecretEnvPriority(t *testing.T) {
	gen, calls := stubGenerator("random", nil)

	secret, source, needPersist, err := resolveSecret("env-secret", "stored-secret", gen)

	if err != nil {
		t.Fatalf("期望无错误，得到: %v", err)
	}
	if secret != "env-secret" {
		t.Errorf("期望密钥来自环境变量 env-secret，得到: %s", secret)
	}
	if source != SecretSourceEnv {
		t.Errorf("期望来源为 env，得到: %s", source)
	}
	if needPersist {
		t.Error("环境变量密钥不应持久化")
	}
	if *calls != 0 {
		t.Errorf("不应调用生成器，实际调用了 %d 次", *calls)
	}
}

func TestResolveSecretFromStored(t *testing.T) {
	gen, calls := stubGenerator("random", nil)

	secret, source, needPersist, err := resolveSecret("", "stored-secret", gen)

	if err != nil {
		t.Fatalf("期望无错误，得到: %v", err)
	}
	if secret != "stored-secret" {
		t.Errorf("期望密钥来自存储值 stored-secret，得到: %s", secret)
	}
	if source != SecretSourceStored {
		t.Errorf("期望来源为 stored，得到: %s", source)
	}
	if needPersist {
		t.Error("已有存储密钥不应再次持久化")
	}
	if *calls != 0 {
		t.Errorf("不应调用生成器，实际调用了 %d 次", *calls)
	}
}

func TestResolveSecretGeneratesWhenMissing(t *testing.T) {
	gen, calls := stubGenerator("random-secret", nil)

	secret, source, needPersist, err := resolveSecret("", "", gen)

	if err != nil {
		t.Fatalf("期望无错误，得到: %v", err)
	}
	if secret != "random-secret" {
		t.Errorf("期望密钥为生成值 random-secret，得到: %s", secret)
	}
	if source != SecretSourceGenerated {
		t.Errorf("期望来源为 generated，得到: %s", source)
	}
	if !needPersist {
		t.Error("生成的密钥需要持久化")
	}
	if *calls != 1 {
		t.Errorf("应调用生成器 1 次，实际调用了 %d 次", *calls)
	}
}

func TestResolveSecretGeneratorError(t *testing.T) {
	gen, _ := stubGenerator("", errors.New("rand failed"))

	_, _, _, err := resolveSecret("", "", gen)

	if err == nil {
		t.Fatal("生成器失败时应返回错误")
	}
	if !strings.Contains(err.Error(), "rand failed") {
		t.Errorf("错误信息应包含生成器错误，得到: %v", err)
	}
}

func TestGenerateRandomSecret(t *testing.T) {
	first, err := generateRandomSecret()
	if err != nil {
		t.Fatalf("期望无错误，得到: %v", err)
	}
	if len(first) != 64 {
		t.Errorf("32字节密钥hex编码应为64字符，得到 %d 字符", len(first))
	}

	second, err := generateRandomSecret()
	if err != nil {
		t.Fatalf("期望无错误，得到: %v", err)
	}
	if first == second {
		t.Error("两次生成的密钥不应相同")
	}
}
