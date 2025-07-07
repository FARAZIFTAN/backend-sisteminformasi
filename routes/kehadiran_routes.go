package routes

import (
	"backend-sisteminformasi/controller"
	"backend-sisteminformasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func KehadiranRoutes(app fiber.Router) {
	app.Get("/kehadiran", middleware.AuthRequired(), controller.GetKehadiran)
	app.Get("/kehadiran/:id", middleware.AuthRequired(), controller.GetKehadiranByID)
	app.Post("/kehadiran", middleware.AuthRequired(), controller.CreateKehadiran)
	app.Put("/kehadiran/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.UpdateKehadiran)
	app.Delete("/kehadiran/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.DeleteKehadiran)
}
