package pan

import "encoding/json"

// XunleiAccountCredentials 迅雷账号凭据结构
type XunleiAccountCredentials struct {
	Username     string `json:"username"`     // 手机号（不包含+86前缀）
	Password     string `json:"password"`     // 密码
	RefreshToken string `json:"refresh_token"` // 当前有效的refresh_token
}

// ParseCredentialsFromCk 从ck字段解析账号凭据
func ParseCredentialsFromCk(ck string) (*XunleiAccountCredentials, error) {
	var credentials XunleiAccountCredentials
	if err := json.Unmarshal([]byte(ck), &credentials); err != nil {
		return nil, err
	}
	return &credentials, nil
}

// IsAccountCredentials 检查ck是否包含账号密码信息
func IsAccountCredentials(ck string) bool {
	var credentials map[string]interface{}
	if err := json.Unmarshal([]byte(ck), &credentials); err != nil {
		return false
	}
	_, hasUsername := credentials["username"]
	_, hasPassword := credentials["password"]
	return hasUsername && hasPassword
}

// ToJsonString 转换为JSON字符串
func (c *XunleiAccountCredentials) ToJsonString() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}