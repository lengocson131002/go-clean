package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/internal/model"
	"github.com/lengocson131002/go-clean/internal/usecase"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/presenter/http/v1/middleware"
	http "github.com/lengocson131002/go-clean/presenter/http/v1/model"
)

type UserController struct {
	Log     logger.LoggerInterface
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger logger.LoggerInterface) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

// GetCrag GetById swagger documentation
// @Summary Get a crag by ID
// @Description Get a crag by ID
// @Tags Crag
// @Accept json
// @Produce json
// @Param id path string true "Crag ID"
// @Router /api/v1/auth/register [get]
func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warn("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to register user : %+v", err)
		return err
	}

	resp := http.SuccessResponse[model.UserResponse](response)
	return resp.JSON(ctx)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warn("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to login user : %+v", err)
		return err
	}

	resp := http.SuccessResponse[model.UserResponse](response)
	return resp.JSON(ctx)
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to get current user")
		return err
	}

	resp := http.SuccessResponse[model.UserResponse](response)
	return resp.JSON(ctx)
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to logout user")
		return err
	}

	resp := http.SuccessResponse[bool](&response)
	return resp.JSON(ctx)
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warn("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = auth.ID
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to update user")
		return err
	}

	resp := http.SuccessResponse[model.UserResponse](response)
	return resp.JSON(ctx)
}
