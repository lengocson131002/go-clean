package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/common"
	"github.com/lengocson131002/go-clean/pkg/logger"
	mapper "github.com/lengocson131002/go-clean/pkg/util"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/usecase/data"
	"golang.org/x/crypto/bcrypt"
)

type loginUserHandler struct {
	Log            logger.Logger
	Validator      validation.Validator
	UserRepository data.UserRepository
}

func NewLoginUserHandler(
	logger logger.Logger,
	validate validation.Validator,
	userRepository data.UserRepository) domain.LoginUserHandler {
	return &loginUserHandler{
		Log:            logger,
		Validator:      validate,
		UserRepository: userRepository,
	}
}

func (c *loginUserHandler) Handle(ctx context.Context, request *domain.LoginUserRequest) (*domain.LoginUserResponse, error) {
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

	res := &domain.LoginUserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}
