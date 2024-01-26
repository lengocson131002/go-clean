package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/internal/model"
	"github.com/lengocson131002/go-clean/internal/usecase"
	"github.com/lengocson131002/go-clean/pkg/logger"
)

func NewAuth(userUserCase *usecase.UserUseCase, log logger.LoggerInterface) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUserCase.Log.Debug("Authorization : %s", request.Token)

		auth, err := userUserCase.Verify(ctx.UserContext(), request)
		if err != nil {
			userUserCase.Log.Warn("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		userUserCase.Log.Debug("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
