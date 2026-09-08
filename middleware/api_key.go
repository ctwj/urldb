package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuth 开放接口 API Key 认证（016-api-access-application，FR-009/FR-016/FR-017）
// 校验顺序：缺失/不匹配 401 → 停用 403 → 过期 403；通过后注入 api_credential_id / api_user_id
func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			key = c.Query("api_key")
		}
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "缺少 API Key",
				"code":    401,
			})
			c.Abort()
			return
		}

		sum := sha256.Sum256([]byte(key))
		keyHash := hex.EncodeToString(sum[:])

		cred, err := repoManager.ApiCredentialRepository.FindByKeyHash(keyHash)
		if err != nil {
			utils.Error("ApiKeyAuth - 凭证查询失败 - error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "系统内部错误",
				"code":    500,
			})
			c.Abort()
			return
		}
		if cred == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "凭证无效",
				"code":    401,
			})
			c.Abort()
			return
		}

		now := utils.GetCurrentTime()
		if cred.Status != entity.ApiCredentialStatusActive {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "权限已停用，请联系管理员或重新申请",
				"code":    403,
			})
			c.Abort()
			return
		}
		if cred.IsExpired(now) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "凭证已过期，请重新申请",
				"code":    403,
			})
			c.Abort()
			return
		}

		// last_used_at 节流更新（距上次 ≥60s 才写库，避免写放大）
		if cred.LastUsedAt == nil || now.Sub(*cred.LastUsedAt) >= time.Minute {
			cred.LastUsedAt = &now
			_ = repoManager.ApiCredentialRepository.Update(cred)
		}

		c.Set("api_credential_id", cred.ID)
		c.Set("api_user_id", cred.UserID)
		c.Next()
	}
}

// ApiRateLimit 开放接口频率限制（016-api-access-application，US4/FR-012/FR-013）
// 每用户三窗口（分钟/小时/天）独立判定，超任一档即 429；0=该档不限
func ApiRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("api_user_id")
		if userID == 0 {
			c.Next()
			return
		}

		allowed, dimension, limit := services.AllowApiRequest(userID)
		if !allowed {
			label := map[string]string{
				"minute": "分钟",
				"hour":   "小时",
				"day":    "天",
			}[dimension]
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "请求过于频繁：每" + label + "上限 " + strconv.Itoa(limit) + " 次",
				"code":    429,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
