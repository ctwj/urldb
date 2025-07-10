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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repoManager.UserRepository.FindByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账户已被禁用"})
		return
	}

	if !middleware.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 更新最后登录时间
	repoManager.UserRepository.UpdateLastLogin(user.ID)

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	response := dto.LoginResponse{
		Token: token,
		User:  converter.ToUserResponse(user),
	}

	c.JSON(http.StatusOK, response)
}

// Register 用户注册
func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := repoManager.UserRepository.FindByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	existingEmail, _ := repoManager.UserRepository.FindByEmail(req.Email)
	if existingEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在"})
		return
	}

	// 哈希密码
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user":    converter.ToUserResponse(user),
	})
}

// GetUsers 获取用户列表（管理员）
func GetUsers(c *gin.Context) {
	users, err := repoManager.UserRepository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := converter.ToUserResponseList(users)
	c.JSON(http.StatusOK, responses)
}

// CreateUser 创建用户（管理员）
func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := repoManager.UserRepository.FindByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	existingEmail, _ := repoManager.UserRepository.FindByEmail(req.Email)
	if existingEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在"})
		return
	}

	// 哈希密码
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户创建成功",
		"user":    converter.ToUserResponse(user),
	})
}

// UpdateUser 更新用户（管理员）
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repoManager.UserRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户更新成功"})
}

// DeleteUser 删除用户（管理员）
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = repoManager.UserRepository.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// GetProfile 获取当前用户信息
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	user, err := repoManager.UserRepository.FindByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	response := converter.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}
