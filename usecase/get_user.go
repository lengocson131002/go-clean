package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/errors"
	"github.com/lengocson131002/go-clean/pkg/logger"
	mapper "github.com/lengocson131002/go-clean/pkg/util"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/usecase/data"
)

type getUserHandler struct {
	Log            logger.Logger
	Validator      validation.Validator
	UserRepository data.UserRepository
}

func NewGetUserHandler(
	logger logger.Logger,
	validate validation.Validator,
	userRepository data.UserRepository) domain.GetUserHandler {
	return &getUserHandler{
		Log:            logger,
		Validator:      validate,
		UserRepository: userRepository,
	}
}

func (c *getUserHandler) Handle(ctx context.Context, request *domain.GetUserRequest) (*domain.GetUserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warn(ctx, "Invalid request body : %+v", err)
		return nil, errors.DomainValidationError
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warnf(ctx, "Failed find user by id : %+v", err)
		return nil, domain.ErrorAccountNotFound
	}

	res := &domain.GetUserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}
