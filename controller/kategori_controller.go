package controller

import (
	"context"
	"time"

	"backend-sisteminformasi/config"
	"backend-sisteminformasi/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var kategoriValidate = validator.New()

// @Security BearerAuth
// GetKategori godoc
// @Summary Get all kategori
// @Tags Kategori
// @Produce json
// @Success 200 {array} model.Kategori
// @Failure 500 {object} map[string]interface{}
// @Router /kategori [get]
func GetKategori(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := config.DB.Collection("kategori").Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil kategori"})
	}
	var kategoris []model.Kategori
	if err := cursor.All(ctx, &kategoris); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode kategori"})
	}
	return c.JSON(kategoris)
}

// @Security BearerAuth
// GetKategoriByID godoc
// @Summary Get kategori by ID
// @Tags Kategori
// @Produce json
// @Param id path string true "Kategori ID"
// @Success 200 {object} model.Kategori
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kategori/{id} [get]
func GetKategoriByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var kategori model.Kategori
	err = config.DB.Collection("kategori").FindOne(ctx, bson.M{"_id": objID}).Decode(&kategori)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Kategori not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch kategori"})
	}
	return c.JSON(kategori)
}

// @Security BearerAuth
// CreateKategori godoc
// @Summary Create kategori
// @Tags Kategori
// @Accept json
// @Produce json
// @Param kategori body model.Kategori true "Kategori Data"
// @Success 201 {object} model.Kategori
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kategori [post]
func CreateKategori(c *fiber.Ctx) error {
	var kategori model.Kategori
	if err := c.BodyParser(&kategori); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := kategoriValidate.Struct(kategori); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := config.DB.Collection("kategori").InsertOne(ctx, kategori)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create kategori"})
	}
	kategori.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return c.Status(201).JSON(kategori)
}

// @Security BearerAuth
// UpdateKategori godoc
// @Summary Update kategori
// @Tags Kategori
// @Accept json
// @Produce json
// @Param id path string true "Kategori ID"
// @Param kategori body model.Kategori true "Kategori Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kategori/{id} [put]
func UpdateKategori(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var kategori model.Kategori
	if err := c.BodyParser(&kategori); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := kategoriValidate.Struct(kategori); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{
		"nama_kategori": kategori.NamaKategori,
	}
	_, err = config.DB.Collection("kategori").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update kategori"})
	}
	return c.JSON(fiber.Map{"message": "Kategori updated"})
}

// @Security BearerAuth
// DeleteKategori godoc
// @Summary Delete kategori
// @Tags Kategori
// @Produce json
// @Param id path string true "Kategori ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kategori/{id} [delete]
func DeleteKategori(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = config.DB.Collection("kategori").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete kategori"})
	}
	return c.JSON(fiber.Map{"message": "Kategori deleted"})
}
