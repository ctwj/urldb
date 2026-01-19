package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key")

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		clientIP := c.ClientIP()

	
		if authHeader == "" {
			// utils.Warn("AuthMiddleware - 未提供认证令牌 - IP: %s, Path: %s", clientIP, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Warn("AuthMiddleware - 无效的认证格式 - IP: %s, Header: %s", clientIP, authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		claims, err := parseToken(tokenString)
		if err != nil {
			utils.Warn("AuthMiddleware - 令牌解析失败 - IP: %s, Error: %v", clientIP, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		// utils.Info("AuthMiddleware - 认证成功 - 用户: %s(ID:%d), 角色: %s, IP: %s",
		// 	claims.Username, claims.UserID, claims.Role, clientIP)

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("client_ip", clientIP)

		c.Next()
	}
}

// AdminMiddleware 管理员中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		username, _ := c.Get("username")
		clientIP, _ := c.Get("client_ip")

		if !exists {
			utils.Warn("AdminMiddleware - 未认证访问管理员接口 - IP: %s, Path: %s", clientIP, c.Request.URL.Path)
			c.Abort()
			return
		}

		if role != "admin" {
			utils.Warn("AdminMiddleware - 非管理员用户尝试访问管理员接口 - 用户: %s, 角色: %s, IP: %s, Path: %s",
				username, role, clientIP, c.Request.URL.Path)
			c.Abort()
			return
		}

				c.Next()
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(user *entity.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(utils.GetCurrentTime().Add(30 * 24 * time.Hour)), // 30天有效期
			IssuedAt:  jwt.NewNumericDate(utils.GetCurrentTime()),
			NotBefore: jwt.NewNumericDate(utils.GetCurrentTime()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// parseToken 解析JWT令牌
func parseToken(tokenString string) (*Claims, error) {
	// utils.Info("parseToken - 开始解析令牌")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		// utils.Error("parseToken - JWT解析失败: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// utils.Info("parseToken - 令牌解析成功，用户ID: %d", claims.UserID)
		return claims, nil
	}

	// utils.Error("parseToken - 令牌无效或签名错误")
	return nil, jwt.ErrSignatureInvalid
}

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 检查密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
