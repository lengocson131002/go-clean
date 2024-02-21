package domain

import (
	"context"
)

type User struct {
	ID        string
	Password  string
	Name      string
	Token     *string
	CreatedAt int64
	UpdatedAt int64
}

type CreateUserRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Name     string `json:"name" validate:"required,max=100"`
}

type CreateUserResponse struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Token     *string `json:"token,omitempty"`
	CreatedAt int64   `json:"created_at,omitempty"`
	UpdatedAt int64   `json:"updated_at,omitempty"`
}

type LoginUserRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginUserResponse struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Token     *string `json:"token,omitempty"`
	CreatedAt int64   `json:"created_at,omitempty"`
	UpdatedAt int64   `json:"updated_at,omitempty"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserResponse struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Token     *string `json:"token,omitempty"`
	CreatedAt int64   `json:"created_at,omitempty"`
	UpdatedAt int64   `json:"updated_at,omitempty"`
}

type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type UpdateUserRequest struct {
	ID       string `json:"-" validate:"required,max=100"`
	Password string `json:"password,omitempty" validate:"max=100"`
	Name     string `json:"name,omitempty" validate:"max=100"`
}

type UpdateUserResponse struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Token     *string `json:"token,omitempty"`
	CreatedAt int64   `json:"created_at,omitempty"`
	UpdatedAt int64   `json:"updated_at,omitempty"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

type VerifyUserResponse struct {
	ID string
}

type CreateUserHandler interface {
	Handle(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error)
}

type GetUserHandler interface {
	Handle(ctx context.Context, request *GetUserRequest) (*GetUserResponse, error)
}

type LoginUserHandler interface {
	Handle(ctx context.Context, request *LoginUserRequest) (*LoginUserResponse, error)
}

type LogoutUserHandler interface {
	Handle(ctx context.Context, request *LogoutUserRequest) (bool, error)
}

type UpdateUserHandler interface {
	Handle(ctx context.Context, request *UpdateUserRequest) (*UpdateUserResponse, error)
}

type VerifyUserHandler interface {
	Handle(ctx context.Context, request *VerifyUserRequest) (*VerifyUserResponse, error)
}
