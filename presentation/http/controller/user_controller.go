package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/http"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/presentation/http/middleware"
)

type UserController struct {
	Log     logger.Logger
	UseCase domain.UserUseCase
}

func NewUserController(useCase domain.UserUseCase, logger logger.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

// Register user, return status code
// @Summary Register user
// @Tags Users
// @Accepts json
// @Produces json
// @Param request body domain.RegisterUserRequest true "RegisterUserRequest request"
// @Success 200 {object} http.DataResponse[domain.UserResponse]
// @Router /api/v1/users/register [post]
func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(domain.RegisterUserRequest)
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

	resp := http.SuccessResponse[domain.UserResponse](response)
	return ctx.Status(resp.Status).JSON(resp)
}

// @Summary Login
// @Description Login User using ID and password
// @Tags Users
// @Accepts json
// @Produces json
// @Param request body domain.LoginUserRequest true "LoginUserRequest request"
// @Success 200 {object} http.DataResponse[domain.UserResponse]
// @Router /api/v1/users/login [get]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(domain.LoginUserRequest)
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

	resp := http.SuccessResponse[domain.UserResponse](response)
	return ctx.Status(resp.Status).JSON(resp)
}

// @Summary Current user
// @Description Get current user
// @Tags Users
// @Accepts json
// @Produces json
// @Param token header string true "Token string"
// @Success 200 {object} http.DataResponse[domain.UserResponse]
// @Router /api/v1/users/me [get]
func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &domain.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to get current user")
		return err
	}

	resp := http.SuccessResponse[domain.UserResponse](response)
	return ctx.Status(resp.Status).JSON(resp)
}

// @Summary Logout
// @Description Log out user
// @Tags Users
// @Accepts json
// @Produces json
// @Param token header string true "Token string"
// @Success 200 {object} http.DataResponse[bool]
// @Router /api/v1/users/me [delete]
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &domain.LogoutUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to logout user")
		return err
	}

	resp := http.SuccessResponse[bool](&response)
	return ctx.Status(resp.Status).JSON(resp)
}

// @Summary Update current user
// @Description Update current user
// @Tags Users
// @Accepts json
// @Produces json
// @Param token header string true "ID string"
// @Param request body domain.UpdateUserRequest true "UpdateUserRequest request"
// @Success 200 {object} http.DataResponse[domain.UserResponse]
// @Router /api/v1/users [put]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(domain.UpdateUserRequest)
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

	resp := http.SuccessResponse[domain.UserResponse](response)
	return ctx.Status(resp.Status).JSON(resp)
}
