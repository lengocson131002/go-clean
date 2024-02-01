package usecase

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/common"
	"github.com/lengocson131002/go-clean/pkg/logger"
	mapper "github.com/lengocson131002/go-clean/pkg/util"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	Log            logger.Logger
	Validator      validation.Validator
	UserRepository domain.UserRepository
}

func NewUserUseCase(
	logger logger.Logger,
	validate validation.Validator,
	userRepository domain.UserRepository) *userUseCase {
	return &userUseCase{
		Log:            logger,
		Validator:      validate,
		UserRepository: userRepository,
	}
}

func (c *userUseCase) Verify(ctx context.Context, request *domain.VerifyUserRequest) (*domain.Auth, error) {
	err := c.Validator.Validate(request)
	if err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, common.ErrBadRequest
	}

	user, err := c.UserRepository.FindByToken(ctx, request.Token)
	if err != nil {
		c.Log.Warn("Failed find user by token : %+v", err)
		return nil, common.ErrBadRequest
	}

	return &domain.Auth{ID: user.ID}, nil
}

func (c *userUseCase) Create(ctx context.Context, request *domain.RegisterUserRequest) (*domain.UserResponse, error) {
	// begin traction

	err := c.Validator.Validate(request)
	if err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, common.ErrBadRequest
	}

	total, err := c.UserRepository.CountById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed count user from database : %+v", err)
		return nil, common.ErrInternalServer
	}

	if total > 0 {
		c.Log.Warn("User already exists")
		return nil, domain.ErrorAccountExisted
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warn("Failed to generate bcrype hash : %+v", err)
		return nil, common.ErrInternalServer
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
		return nil, common.ErrInternalServer
	}

	res := &domain.UserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}

func (c *userUseCase) Login(ctx context.Context, request *domain.LoginUserRequest) (*domain.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body  : %+v", err)
		return nil, common.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return nil, common.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warn("Failed to compare user password with bcrype hash : %+v", err)
		return nil, common.ErrUnauthorized
	}

	token := uuid.New().String()
	user.Token = &token

	if err := c.UserRepository.UpdateUser(ctx, user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return nil, common.ErrInternalServer
	}

	res := &domain.UserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}

func (c *userUseCase) Current(ctx context.Context, request *domain.GetUserRequest) (*domain.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, common.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return nil, common.ErrNotFound
	}

	res := &domain.UserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}

func (c *userUseCase) Logout(ctx context.Context, request *domain.LogoutUserRequest) (bool, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return false, common.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return false, common.ErrNotFound
	}

	user.Token = nil

	if err := c.UserRepository.UpdateUser(ctx, user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return false, common.ErrInternalServer
	}

	return true, nil
}

func (c *userUseCase) Update(ctx context.Context, request *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, common.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return nil, common.ErrNotFound
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warn("Failed to generate bcrype hash : %+v", err)
			return nil, common.ErrInternalServer
		}
		user.Password = string(password)
	}

	if err := c.UserRepository.UpdateUser(ctx, user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return nil, common.ErrInternalServer
	}

	res := &domain.UserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}
