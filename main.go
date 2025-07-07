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
	"github.com/joho/godotenv"
)

// @title Sistem Informasi UKM Kampus API
// @version 1.0
// @description API untuk tugas besar backend Fiber + MongoDB
// @host localhost:3000
// @BasePath /
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

	app.Use(middleware.Logger())
	app.Use(middleware.Cors())

	// Swagger endpoint
	app.Get("/swagger/*", fiberswagger.WrapHandler)

	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
