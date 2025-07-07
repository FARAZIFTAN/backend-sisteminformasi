package routes

import (
	"backend-sisteminformasi/controller"
	"backend-sisteminformasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router) {
	app.Get("/users", middleware.AuthRequired(), middleware.AdminOnly(), controller.GetUsers)
	app.Get("/users/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.GetUser)
	app.Post("/users", middleware.AuthRequired(), middleware.AdminOnly(), controller.CreateUser)
	app.Put("/users/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.UpdateUser)
	app.Delete("/users/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.DeleteUser)
}
