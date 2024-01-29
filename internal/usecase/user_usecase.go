package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/lengocson131002/go-clean/internal/domain"
	"github.com/lengocson131002/go-clean/internal/domain/response"
	repo "github.com/lengocson131002/go-clean/internal/interfaces"
	"github.com/lengocson131002/go-clean/internal/mapper"
	"github.com/lengocson131002/go-clean/internal/model"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	Log            logger.LoggerInterface
	Validator      validation.Validator
	UserRepository repo.UserRepositoryInterface
}

func NewUserUseCase(
	logger logger.LoggerInterface,
	validate validation.Validator,
	userRepository repo.UserRepositoryInterface) *UserUseCase {
	return &UserUseCase{
		Log:            logger,
		Validator:      validate,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	err := c.Validator.Validate(request)
	if err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindByToken(request.Token)
	if err != nil {
		c.Log.Warn("Failed find user by token : %+v", err)
		return nil, response.ErrBadRequest
	}

	return &model.Auth{ID: user.ID}, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	// begin traction

	err := c.Validator.Validate(request)
	if err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, response.ErrBadRequest
	}

	total, err := c.UserRepository.CountById(request.ID)
	if err != nil {
		c.Log.Warn("Failed count user from database : %+v", err)
		return nil, response.ErrInternalServer
	}

	if total > 0 {
		c.Log.Warn("User already exists : %+v", err)
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

	if err := c.UserRepository.CreateUser(user); err != nil {
		c.Log.Warn("Failed create user to database : %+v", err)
		return nil, response.ErrInternalServer
	}

	return mapper.UserToResponse(user), nil
	//
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body  : %+v", err)
		return nil, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return nil, response.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warn("Failed to compare user password with bcrype hash : %+v", err)
		return nil, response.ErrUnauthorized
	}

	user.Token = uuid.New().String()
	if err := c.UserRepository.UpdateUser(user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return nil, response.ErrInternalServer
	}

	return mapper.UserToTokenResponse(user), nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return nil, response.ErrNotFound
	}

	return mapper.UserToResponse(user), nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return false, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(request.ID)
	if err != nil {
		c.Log.Warn("Failed find user by id : %+v", err)
		return false, response.ErrNotFound
	}

	user.Token = ""

	if err := c.UserRepository.UpdateUser(user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return false, response.ErrInternalServer
	}

	return true, nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, response.ErrBadRequest
	}

	user, err := c.UserRepository.FindUserById(request.ID)
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

	if err := c.UserRepository.UpdateUser(user); err != nil {
		c.Log.Warn("Failed save user : %+v", err)
		return nil, response.ErrInternalServer
	}

	return mapper.UserToResponse(user), nil
}
