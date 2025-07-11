package handlers

import (
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"
	"res_db/middleware"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := repoManager.UserRepository.FindByUsername(req.Username)
	if err != nil {
		ErrorResponse(c, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		ErrorResponse(c, "账户已被禁用", http.StatusUnauthorized)
		return
	}

	if !middleware.CheckPassword(req.Password, user.Password) {
		ErrorResponse(c, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	// 更新最后登录时间
	repoManager.UserRepository.UpdateLastLogin(user.ID)

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		ErrorResponse(c, "生成令牌失败", http.StatusInternalServerError)
		return
	}

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

	// 检查用户名是否已存在
	existingUser, _ := repoManager.UserRepository.FindByUsername(req.Username)
	if existingUser != nil {
		ErrorResponse(c, "用户名已存在", http.StatusBadRequest)
		return
	}

	// 检查邮箱是否已存在
	existingEmail, _ := repoManager.UserRepository.FindByEmail(req.Email)
	if existingEmail != nil {
		ErrorResponse(c, "邮箱已存在", http.StatusBadRequest)
		return
	}

	// 哈希密码
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
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
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

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

	// 检查用户名是否已存在
	existingUser, _ := repoManager.UserRepository.FindByUsername(req.Username)
	if existingUser != nil {
		ErrorResponse(c, "用户名已存在", http.StatusBadRequest)
		return
	}

	// 检查邮箱是否已存在
	existingEmail, _ := repoManager.UserRepository.FindByEmail(req.Email)
	if existingEmail != nil {
		ErrorResponse(c, "邮箱已存在", http.StatusBadRequest)
		return
	}

	// 哈希密码
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
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
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

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

	user, err := repoManager.UserRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

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
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "用户更新成功"})
}

// DeleteUser 删除用户（管理员）
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.UserRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "用户删除成功"})
}

// GetProfile 获取用户资料
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, "未认证", http.StatusUnauthorized)
		return
	}

	user, err := repoManager.UserRepository.FindByID(userID.(uint))
	if err != nil {
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

	response := converter.ToUserResponse(user)
	SuccessResponse(c, response)
}
