package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/errors"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/usecase/data"
)

type logoutUserHandler struct {
	Log            logger.Logger
	Validator      validation.Validator
	UserRepository data.UserRepository
}

func NewLogoutUserHandler(
	logger logger.Logger,
	validate validation.Validator,
	userRepository data.UserRepository) domain.LogoutUserHandler {
	return &logoutUserHandler{
		Log:            logger,
		Validator:      validate,
		UserRepository: userRepository,
	}
}

func (c *logoutUserHandler) Handle(ctx context.Context, request *domain.LogoutUserRequest) (bool, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warnf(ctx, "Invalid request body : %+v", err)
		return false, errors.DomainValidationError
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warnf(ctx, "Failed find user by id : %+v", err)
		return false, domain.ErrorAccountNotFound
	}

	user.Token = nil

	if err := c.UserRepository.UpdateUser(ctx, user); err != nil {
		c.Log.Warnf(ctx, "Failed save user : %+v", err)
		return false, errors.InternalServerError
	}

	return true, nil
}
