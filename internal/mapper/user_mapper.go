package mapper

import (
	"github.com/lengocson131002/go-clean/internal/domain"
	"github.com/lengocson131002/go-clean/internal/model"
)

func UserToResponse(user *domain.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToTokenResponse(user *domain.User) *model.UserResponse {
	return &model.UserResponse{
		Token: user.Token,
	}
}
