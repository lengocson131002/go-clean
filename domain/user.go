package domain

import (
	"context"

	"github.com/lengocson131002/go-clean/pkg/database"
)

type User struct {
	ID        string
	Password  string
	Name      string
	Token     *string
	CreatedAt int64
	UpdatedAt int64
}

// Data interfaces
type UserRepository interface {
	FindByToken(ctx context.Context, token string) (*User, error)
	CountById(ctx context.Context, id string) (int64, error)
	CreateUser(ctx context.Context, user *User) error
	FindUserById(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	database.EnableTransactor
}

// Usecase interfaces
type UserUseCase interface {
	Verify(ctx context.Context, request *VerifyUserRequest) (*Auth, error)
	Create(ctx context.Context, request *RegisterUserRequest) (*UserResponse, error)
	Login(ctx context.Context, request *LoginUserRequest) (*UserResponse, error)
	Current(ctx context.Context, request *GetUserRequest) (*UserResponse, error)
	Logout(ctx context.Context, request *LogoutUserRequest) (bool, error)
	Update(ctx context.Context, request *UpdateUserRequest) (*UserResponse, error)
}
