package repo

import "github.com/lengocson131002/go-clean/internal/domain"

type UserRepositoryInterface interface {
	FindByToken(token string) (*domain.User, error)

	CountById(id string) (int64, error)

	CreateUser(user *domain.User) error

	FindUserById(id string) (*domain.User, error)

	UpdateUser(user *domain.User) error
}
