package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/pipeline"
)

type AuthMiddleware struct {
	log logger.Logger
}

func NewAuthMiddleware(log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		log: log,
	}
}

func (m *AuthMiddleware) Handle(ctx *fiber.Ctx) error {
	request := &domain.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
	m.log.Debugf(ctx.Context(), "Authorization : %s", request.Token)

	auth, err := pipeline.Send[*domain.VerifyUserRequest, *domain.VerifyUserResponse](ctx.Context(), request)
	if err != nil {
		m.log.Warnf(ctx.Context(), "Failed find user by token : %+v", err)
		return fiber.ErrUnauthorized
	}

	m.log.Debugf(ctx.Context(), "User : %+v", auth.ID)
	ctx.Locals("auth", auth)
	return ctx.Next()
}

func GetUser(ctx *fiber.Ctx) *domain.VerifyUserResponse {
	return ctx.Locals("auth").(*domain.VerifyUserResponse)
}
