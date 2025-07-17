package converter

import (
	"github.com/ctwj/panResManage/db/dto"
	"github.com/ctwj/panResManage/db/entity"
)

// ToUserResponse 将User实体转换为UserResponse
func ToUserResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserResponseList 将User实体列表转换为UserResponse列表
func ToUserResponseList(users []entity.User) []dto.UserResponse {
	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(&user)
	}
	return responses
}
