package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	AuthRoutes(app)
	UserRoutes(app)
	KegiatanRoutes(app)
	KehadiranRoutes(app)
	KategoriRoutes(app)
	StatisticsRoutes(app)
}
