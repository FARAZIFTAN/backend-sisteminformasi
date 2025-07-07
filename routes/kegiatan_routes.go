package routes

import (
	"backend-sisteminformasi/controller"
	"backend-sisteminformasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func KegiatanRoutes(app fiber.Router) {
	app.Get("/kegiatan", middleware.AuthRequired(), controller.GetKegiatan)
	app.Get("/kegiatan/:id", middleware.AuthRequired(), controller.GetKegiatanByID)
	app.Post("/kegiatan", middleware.AuthRequired(), middleware.AdminOnly(), controller.CreateKegiatan)
	app.Put("/kegiatan/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.UpdateKegiatan)
	app.Delete("/kegiatan/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.DeleteKegiatan)
}
