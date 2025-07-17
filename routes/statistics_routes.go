package routes

import (
	"backend-sisteminformasi/controller"
	"backend-sisteminformasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func StatisticsRoutes(app fiber.Router) {
	app.Get("/statistics", middleware.AuthRequired(), middleware.AdminOnly(), controller.GetStatistics)
}
