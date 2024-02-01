package usecase

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lengocson131002/go-clean/internal/domain"
	"github.com/lengocson131002/go-clean/internal/domain/response"
	"github.com/lengocson131002/go-clean/internal/interfaces"
	"github.com/lengocson131002/go-clean/internal/model"
	"github.com/lengocson131002/go-clean/pkg/logger"
	mapper "github.com/lengocson131002/go-clean/pkg/util"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error)
	Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error)
	Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error)
	Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error)
}

type userUseCase struct {
	Log            logger.Logger
	Validator      validation.Validator
	UserRepository interfaces.UserRepositoryInterface
}

func NewUserUseCase(
	logger logger.Logger,
	validate validation.Validator,
	userRepository interfaces.UserRepositoryInterface) *userUseCase {
	return &userUseCase{
		Log:            logger,
		Validator:      validate,
		UserRepository: userRepository,
	}
}

func (c *userUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	err := c.Validator.Validate(request)
	if err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindByToken(ctx, request.Token)
	if err != nil {
		c.Log.Warn("Failed find user by token : %+v", err)
		return nil, response.ErrBadRequest
	}

	return &model.Auth{ID: user.ID}, nil
}

func (c *userUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	// begin traction

	err := c.Validator.Validate(request)
	if err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, response.ErrBadRequest
	}

	total, err := c.UserRepository.CountById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed count user from database : %+v", err)
		return nil, response.ErrInternalServer
	}

	if total > 0 {
		c.Log.Warn("User already exists")
		return nil, response.ErrorAccountExisted
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warn("Failed to generate bcrype hash : %+v", err)
		return nil, response.ErrInternalServer
	}

	user := &domain.User{
		ID:       request.ID,
		Password: string(password),
		Name:     request.Name,
	}

	err = c.UserRepository.WithinTransactionOptions(ctx, func(ctx context.Context) error {
		return c.UserRepository.CreateUser(ctx, user)
	}, &sql.TxOptions{
		Isolation: sql.IsolationLevel(2),
		ReadOnly:  false,
	})

	if err != nil {
		c.Log.Warn("Failed create user to database : %+v", err)
		return nil, response.ErrInternalServer
	}

	res := &model.UserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}

func (c *userUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body  : %+v", err)
		return nil, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return nil, response.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warn("Failed to compare user password with bcrype hash : %+v", err)
		return nil, response.ErrUnauthorized
	}

	token := uuid.New().String()
	user.Token = &token

	if err := c.UserRepository.UpdateUser(ctx, user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return nil, response.ErrInternalServer
	}

	res := &model.UserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}

func (c *userUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return nil, response.ErrNotFound
	}

	res := &model.UserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}

func (c *userUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return false, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return false, response.ErrNotFound
	}

	user.Token = nil

	if err := c.UserRepository.UpdateUser(ctx, user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return false, response.ErrInternalServer
	}

	return true, nil
}

func (c *userUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return nil, response.ErrNotFound
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warn("Failed to generate bcrype hash : %+v", err)
			return nil, response.ErrInternalServer
		}
		user.Password = string(password)
	}

	if err := c.UserRepository.UpdateUser(ctx, user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return nil, response.ErrInternalServer
	}

	res := &model.UserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}
