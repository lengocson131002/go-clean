package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/internal/model"
	"github.com/lengocson131002/go-clean/internal/usecase"
	"github.com/lengocson131002/go-clean/pkg/logger"
)

type AuthMiddleware struct {
	userUseCase *usecase.UserUseCase
	log         logger.Logger
}

func NewAuthMiddleware(userUserCase *usecase.UserUseCase, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		userUseCase: userUserCase,
		log:         log,
	}
}

func (m *AuthMiddleware) Handle(ctx *fiber.Ctx) error {
	request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
	m.log.Debug("Authorization : %s", request.Token)

	auth, err := m.userUseCase.Verify(ctx.UserContext(), request)
	if err != nil {
		m.userUseCase.Log.Warn("Failed find user by token : %+v", err)
		return fiber.ErrUnauthorized
	}

	m.log.Debug("User : %+v", auth.ID)
	ctx.Locals("auth", auth)
	return ctx.Next()
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
