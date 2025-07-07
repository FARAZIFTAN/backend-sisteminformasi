package controller

import (
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	return c.SendString("Get all users")
}

func GetUser(c *fiber.Ctx) error {
	return c.SendString("Get user by ID")
}

func CreateUser(c *fiber.Ctx) error {
	return c.SendString("Create user")
}

func UpdateUser(c *fiber.Ctx) error {
	return c.SendString("Update user")
}

func DeleteUser(c *fiber.Ctx) error {
	return c.SendString("Delete user")
}
