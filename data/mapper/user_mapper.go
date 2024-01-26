package mapper

import (
	"github.com/lengocson131002/go-clean/data/entity"
	"github.com/lengocson131002/go-clean/internal/domain"
)

func ToUserEntity(user *domain.User) *entity.UserEntity {
	return &entity.UserEntity{
		ID:        user.ID,
		Password:  user.Password,
		Name:      user.Name,
		Token:     user.Token,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserDomain(user *entity.UserEntity) *domain.User {
	return &domain.User{
		ID:        user.ID,
		Password:  user.Password,
		Name:      user.Name,
		Token:     user.Token,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
