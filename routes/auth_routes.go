package routes

import (
	"backend-sisteminformasi/controller"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app fiber.Router) {
	app.Post("/register", controller.Register)
	app.Post("/login", controller.Login)
}
