package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/errors"
	"github.com/lengocson131002/go-clean/pkg/logger"
	mapper "github.com/lengocson131002/go-clean/pkg/util"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/usecase/data"
	"golang.org/x/crypto/bcrypt"
)

type updateUserHandler struct {
	Log            logger.Logger
	Validator      validation.Validator
	UserRepository data.UserRepository
}

func NewUpdateUserHandler(
	logger logger.Logger,
	validate validation.Validator,
	userRepository data.UserRepository) domain.UpdateUserHandler {
	return &updateUserHandler{
		Log:            logger,
		Validator:      validate,
		UserRepository: userRepository,
	}
}
func (c *updateUserHandler) Handle(ctx context.Context, request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	if err := c.Validator.Validate(request); err != nil {
		c.Log.Warnf(ctx, "Invalid request body : %+v", err)
		return nil, errors.DomainValidationError
	}

	user, err := c.UserRepository.FindUserById(ctx, request.ID)
	if err != nil {
		c.Log.Warnf(ctx, "Failed find user by id : %+v", err)
		return nil, errors.DomainValidationError
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf(ctx, "Failed to generate bcrype hash : %+v", err)
			return nil, domain.ErrorAccountNotFound
		}
		user.Password = string(password)
	}

	if err := c.UserRepository.UpdateUser(ctx, user); err != nil {
		c.Log.Warnf(ctx, "Failed save user : %+v", err)
		return nil, errors.InternalServerError
	}

	res := &domain.UpdateUserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}
