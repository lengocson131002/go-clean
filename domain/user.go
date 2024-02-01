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

// Models
type Auth struct {
	// Login user id
	ID string
}

type UserResponse struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Token     *string `json:"token,omitempty"`
	CreatedAt int64   `json:"created_at,omitempty"`
	UpdatedAt int64   `json:"updated_at,omitempty"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

type RegisterUserRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Name     string `json:"name" validate:"required,max=100"`
}

type UpdateUserRequest struct {
	ID       string `json:"-" validate:"required,max=100"`
	Password string `json:"password,omitempty" validate:"max=100"`
	Name     string `json:"name,omitempty" validate:"max=100"`
}

type LoginUserRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
