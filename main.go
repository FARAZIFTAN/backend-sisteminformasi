package main

import (
	"log"
	"os"

	"backend-sisteminformasi/config"
	_ "backend-sisteminformasi/docs"
	"backend-sisteminformasi/middleware"
	"backend-sisteminformasi/routes"

	fiberswagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// @title Sistem Informasi UKM Kampus API
// @version 1.0
// @description API untuk tugas besar backend Fiber + MongoDB
// @host backend-sisteminformasi-production.up.railway.app
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDB()

	// Seed admin user jika belum ada
	config.SeedAdminUser(config.DB)

	app := fiber.New()

	// Logger middleware
	app.Use(middleware.Logger())

	// CORS middleware (hanya satu, konfigurasi sudah benar)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Swagger endpoint
	app.Get("/swagger/*", fiberswagger.WrapHandler)

	// Setup all routes
	routes.SetupRoutes(app)

	// Health check endpoint for Railway/Render
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
