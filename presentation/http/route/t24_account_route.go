package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/presentation/http/controller"
)

func RegisterT24Route(root *fiber.Router, t24Con *controller.T24AccountController) {
	t24Acc := (*root).Group("/t24/accounts")

	t24Acc.Post("/open", t24Con.OpenAccount)
	t24Acc.Post("/open-rpc", t24Con.OpenAccountRpc)
}
