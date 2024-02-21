package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/common"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/usecase/data"
)

type verifyUserHandler struct {
	Log            logger.Logger
	Validator      validation.Validator
	UserRepository data.UserRepository
}

func NewVerifyUserHandler(log logger.Logger, validator validation.Validator, userRepository data.UserRepository) domain.VerifyUserHandler {
	return &verifyUserHandler{
		Log:            log,
		Validator:      validator,
		UserRepository: userRepository,
	}
}

func (c *verifyUserHandler) Handle(ctx context.Context, request *domain.VerifyUserRequest) (*domain.VerifyUserResponse, error) {
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

	return &domain.VerifyUserResponse{ID: user.ID}, nil
}
