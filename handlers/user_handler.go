package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/middleware"
	"github.com/ctwj/urldb/plugin-system/triggers/plugins"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	clientIP, _ := c.Get("client_ip")
	utils.Info("Login - 尝试登录 - 用户名: %s, IP: %s", req.Username, clientIP)

	user, err := repoManager.UserRepository.FindByUsername(req.Username)
	if err != nil {
		utils.Warn("Login - 用户不存在或密码错误 - 用户名: %s, IP: %s", req.Username, clientIP)
		ErrorResponse(c, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		utils.Warn("Login - 账户已被禁用 - 用户名: %s, IP: %s", req.Username, clientIP)
		ErrorResponse(c, "账户已被禁用", http.StatusUnauthorized)
		return
	}

	if !middleware.CheckPassword(req.Password, user.Password) {
		utils.Warn("Login - 密码错误 - 用户名: %s, IP: %s", req.Username, clientIP)
		ErrorResponse(c, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	// 更新最后登录时间
	repoManager.UserRepository.UpdateLastLogin(user.ID)

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		utils.Error("Login - 生成令牌失败 - 用户名: %s, IP: %s, Error: %v", req.Username, clientIP, err)
		ErrorResponse(c, "生成令牌失败", http.StatusInternalServerError)
		return
	}

	utils.Info("Login - 登录成功 - 用户名: %s(ID:%d), IP: %s", req.Username, user.ID, clientIP)

	// 触发用户登录事件
	loginData := map[string]interface{}{
		"ip":         clientIP,
		"user_agent": c.GetHeader("User-Agent"),
		"login_time": time.Now(),
	}
	plugins.TriggerUserLogin(user, loginData)

	response := dto.LoginResponse{
		Token: token,
		User:  converter.ToUserResponse(user),
	}

	SuccessResponse(c, response)
}

// Register 用户注册
func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	clientIP, _ := c.Get("client_ip")
	utils.Info("Register - 尝试注册 - 用户名: %s, 邮箱: %s, IP: %s", req.Username, req.Email, clientIP)

	// 检查用户名是否已存在
	existingUser, _ := repoManager.UserRepository.FindByUsername(req.Username)
	if existingUser != nil {
		utils.Warn("Register - 用户名已存在 - 用户名: %s, IP: %s", req.Username, clientIP)
		ErrorResponse(c, "用户名已存在", http.StatusBadRequest)
		return
	}

	// 检查邮箱是否已存在
	existingEmail, _ := repoManager.UserRepository.FindByEmail(req.Email)
	if existingEmail != nil {
		utils.Warn("Register - 邮箱已存在 - 邮箱: %s, IP: %s", req.Email, clientIP)
		ErrorResponse(c, "邮箱已存在", http.StatusBadRequest)
		return
	}

	// 哈希密码
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		utils.Error("Register - 密码加密失败 - 用户名: %s, IP: %s, Error: %v", req.Username, clientIP, err)
		ErrorResponse(c, "密码加密失败", http.StatusInternalServerError)
		return
	}

	user := &entity.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     "user",
		IsActive: true,
	}

	err = repoManager.UserRepository.Create(user)
	if err != nil {
		utils.Error("Register - 创建用户失败 - 用户名: %s, IP: %s, Error: %v", req.Username, clientIP, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("Register - 注册成功 - 用户名: %s(ID:%d), 邮箱: %s, IP: %s", req.Username, user.ID, req.Email, clientIP)

	SuccessResponse(c, gin.H{
		"message": "注册成功",
		"user":    converter.ToUserResponse(user),
	})
}

// GetUsers 获取用户列表（管理员）
func GetUsers(c *gin.Context) {
	users, err := repoManager.UserRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToUserResponseList(users)
	SuccessResponse(c, responses)
}

// CreateUser 创建用户（管理员）
func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	adminUsername, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("CreateUser - 管理员创建用户 - 管理员: %s, 新用户名: %s, IP: %s", adminUsername, req.Username, clientIP)

	// 检查用户名是否已存在
	existingUser, _ := repoManager.UserRepository.FindByUsername(req.Username)
	if existingUser != nil {
		utils.Warn("CreateUser - 用户名已存在 - 管理员: %s, 用户名: %s, IP: %s", adminUsername, req.Username, clientIP)
		ErrorResponse(c, "用户名已存在", http.StatusBadRequest)
		return
	}

	// 检查邮箱是否已存在
	existingEmail, _ := repoManager.UserRepository.FindByEmail(req.Email)
	if existingEmail != nil {
		utils.Warn("CreateUser - 邮箱已存在 - 管理员: %s, 邮箱: %s, IP: %s", adminUsername, req.Email, clientIP)
		ErrorResponse(c, "邮箱已存在", http.StatusBadRequest)
		return
	}

	// 哈希密码
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		utils.Error("CreateUser - 密码加密失败 - 管理员: %s, 用户名: %s, IP: %s, Error: %v", adminUsername, req.Username, clientIP, err)
		ErrorResponse(c, "密码加密失败", http.StatusInternalServerError)
		return
	}

	user := &entity.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	err = repoManager.UserRepository.Create(user)
	if err != nil {
		utils.Error("CreateUser - 创建用户失败 - 管理员: %s, 用户名: %s, IP: %s, Error: %v", adminUsername, req.Username, clientIP, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("CreateUser - 用户创建成功 - 管理员: %s, 用户名: %s(ID:%d), 角色: %s, IP: %s", adminUsername, req.Username, user.ID, req.Role, clientIP)

	SuccessResponse(c, gin.H{
		"message": "用户创建成功",
		"user":    converter.ToUserResponse(user),
	})
}

// UpdateUser 更新用户（管理员）
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	adminUsername, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("UpdateUser - 管理员更新用户 - 管理员: %s, 目标用户ID: %d, IP: %s", adminUsername, id, clientIP)

	user, err := repoManager.UserRepository.FindByID(uint(id))
	if err != nil {
		utils.Warn("UpdateUser - 目标用户不存在 - 管理员: %s, 用户ID: %d, IP: %s", adminUsername, id, clientIP)
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

	// 记录变更前的信息
	oldInfo := fmt.Sprintf("用户名:%s,邮箱:%s,角色:%s,状态:%t", user.Username, user.Email, user.Role, user.IsActive)
	utils.Debug("UpdateUser - 更新前用户信息 - 管理员: %s, 用户ID: %d, 信息: %s", adminUsername, id, oldInfo)

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	user.IsActive = req.IsActive

	err = repoManager.UserRepository.Update(user)
	if err != nil {
		utils.Error("UpdateUser - 更新用户失败 - 管理员: %s, 用户ID: %d, IP: %s, Error: %v", adminUsername, id, clientIP, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 记录变更后信息
	newInfo := fmt.Sprintf("用户名:%s,邮箱:%s,角色:%s,状态:%t", user.Username, user.Email, user.Role, user.IsActive)
	utils.Info("UpdateUser - 用户更新成功 - 管理员: %s, 用户ID: %d, 更新前: %s, 更新后: %s, IP: %s", adminUsername, id, oldInfo, newInfo, clientIP)

	SuccessResponse(c, gin.H{"message": "用户更新成功"})
}

// ChangePassword 修改用户密码（管理员）
func ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	adminUsername, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("ChangePassword - 管理员修改用户密码 - 管理员: %s, 目标用户ID: %d, IP: %s", adminUsername, id, clientIP)

	user, err := repoManager.UserRepository.FindByID(uint(id))
	if err != nil {
		utils.Warn("ChangePassword - 目标用户不存在 - 管理员: %s, 用户ID: %d, IP: %s", adminUsername, id, clientIP)
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

	// 哈希新密码
	hashedPassword, err := middleware.HashPassword(req.NewPassword)
	if err != nil {
		utils.Error("ChangePassword - 密码加密失败 - 管理员: %s, 用户ID: %d, IP: %s, Error: %v", adminUsername, id, clientIP, err)
		ErrorResponse(c, "密码加密失败", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	err = repoManager.UserRepository.Update(user)
	if err != nil {
		utils.Error("ChangePassword - 更新密码失败 - 管理员: %s, 用户ID: %d, IP: %s, Error: %v", adminUsername, id, clientIP, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("ChangePassword - 密码修改成功 - 管理员: %s, 用户名: %s(ID:%d), IP: %s", adminUsername, user.Username, id, clientIP)

	SuccessResponse(c, gin.H{"message": "密码修改成功"})
}

// DeleteUser 删除用户（管理员）
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	adminUsername, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("DeleteUser - 管理员删除用户 - 管理员: %s, 目标用户ID: %d, IP: %s", adminUsername, id, clientIP)

	// 先获取用户信息用于日志记录
	user, err := repoManager.UserRepository.FindByID(uint(id))
	if err != nil {
		utils.Warn("DeleteUser - 目标用户不存在 - 管理员: %s, 用户ID: %d, IP: %s", adminUsername, id, clientIP)
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

	err = repoManager.UserRepository.Delete(uint(id))
	if err != nil {
		utils.Error("DeleteUser - 删除用户失败 - 管理员: %s, 用户ID: %d, IP: %s, Error: %v", adminUsername, id, clientIP, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("DeleteUser - 用户删除成功 - 管理员: %s, 用户名: %s(ID:%d), IP: %s", adminUsername, user.Username, id, clientIP)

	SuccessResponse(c, gin.H{"message": "用户删除成功"})
}

// GetProfile 获取用户资料
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, "未认证", http.StatusUnauthorized)
		return
	}

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("GetProfile - 用户获取个人资料 - 用户名: %s(ID:%d), IP: %s", username, userID, clientIP)

	user, err := repoManager.UserRepository.FindByID(userID.(uint))
	if err != nil {
		utils.Warn("GetProfile - 用户不存在 - 用户名: %s(ID:%d), IP: %s", username, userID, clientIP)
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

	response := converter.ToUserResponse(user)
	utils.Debug("GetProfile - 成功获取个人资料 - 用户名: %s(ID:%d), IP: %s", username, userID, clientIP)
	SuccessResponse(c, response)
}
