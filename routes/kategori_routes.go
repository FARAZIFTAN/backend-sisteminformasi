package routes

import (
	"backend-sisteminformasi/controller"
	"backend-sisteminformasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func KategoriRoutes(app fiber.Router) {
	app.Get("/kategori", controller.GetKategori) // Public GET all kategori
	app.Get("/kategori/:id", middleware.AuthRequired(), controller.GetKategoriByID)
	app.Post("/kategori", middleware.AuthRequired(), middleware.AdminOnly(), controller.CreateKategori)
	app.Put("/kategori/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.UpdateKategori)
	app.Delete("/kategori/:id", middleware.AuthRequired(), middleware.AdminOnly(), controller.DeleteKategori)
}
