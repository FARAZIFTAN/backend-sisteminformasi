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

var kehadiranValidate = validator.New()

// @Security BearerAuth
// GetKehadiran godoc
// @Summary Get all kehadiran
// @Tags Kehadiran
// @Produce json
// @Success 200 {array} model.Kehadiran
// @Failure 500 {object} map[string]interface{}
// @Router /kehadiran [get]
func GetKehadiran(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Aggregation pipeline untuk JOIN dengan users dan kegiatan
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"user_object_id": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$user_id"}, "string"}},
						"then": bson.M{"$toObjectId": "$user_id"},
						"else": "$user_id",
					},
				},
				"kegiatan_object_id": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$kegiatan_id"}, "string"}},
						"then": bson.M{"$toObjectId": "$kegiatan_id"},
						"else": "$kegiatan_id",
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_object_id",
				"foreignField": "_id",
				"as":           "user_data",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "kegiatan",
				"localField":   "kegiatan_object_id",
				"foreignField": "_id",
				"as":           "kegiatan_data",
			},
		},
		{
			"$project": bson.M{
				"_id":           1,
				"user_id":       1,
				"kegiatan_id":   1,
				"status":        1,
				"waktu_cek":     1,
				"user_nama":     bson.M{"$arrayElemAt": []interface{}{"$user_data.nama", 0}},
				"kegiatan_nama": bson.M{"$arrayElemAt": []interface{}{"$kegiatan_data.judul", 0}},
				"ukm":           bson.M{"$arrayElemAt": []interface{}{"$kegiatan_data.kategori", 0}},
			},
		},
	}

	cursor, err := config.DB.Collection("kehadiran").Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch kehadiran"})
	}

	var kehadirans []bson.M
	if err := cursor.All(ctx, &kehadirans); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode kehadiran"})
	}

	return c.JSON(kehadirans)
}

// @Security BearerAuth
// GetKehadiranByID godoc
// @Summary Get kehadiran by ID
// @Tags Kehadiran
// @Produce json
// @Param id path string true "Kehadiran ID"
// @Success 200 {object} model.Kehadiran
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kehadiran/{id} [get]
func GetKehadiranByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var kehadiran model.Kehadiran
	err = config.DB.Collection("kehadiran").FindOne(ctx, bson.M{"_id": objID}).Decode(&kehadiran)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Kehadiran not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch kehadiran"})
	}
	return c.JSON(kehadiran)
}

// @Security BearerAuth
// CreateKehadiran godoc
// @Summary Create kehadiran
// @Tags Kehadiran
// @Accept json
// @Produce json
// @Param kehadiran body model.Kehadiran true "Kehadiran Data"
// @Success 201 {object} model.Kehadiran
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kehadiran [post]
func CreateKehadiran(c *fiber.Ctx) error {
	var kehadiran model.Kehadiran
	if err := c.BodyParser(&kehadiran); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := kehadiranValidate.Struct(kehadiran); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := config.DB.Collection("kehadiran").InsertOne(ctx, kehadiran)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create kehadiran"})
	}
	kehadiran.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return c.Status(201).JSON(kehadiran)
}

// @Security BearerAuth
// UpdateKehadiran godoc
// @Summary Update kehadiran
// @Tags Kehadiran
// @Accept json
// @Produce json
// @Param id path string true "Kehadiran ID"
// @Param kehadiran body model.Kehadiran true "Kehadiran Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kehadiran/{id} [put]
func UpdateKehadiran(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var kehadiran model.Kehadiran
	if err := c.BodyParser(&kehadiran); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := kehadiranValidate.Struct(kehadiran); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{
		"user_id":     kehadiran.UserID,
		"kegiatan_id": kehadiran.KegiatanID,
		"status":      kehadiran.Status,
		"waktu_cek":   kehadiran.WaktuCek,
	}
	_, err = config.DB.Collection("kehadiran").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update kehadiran"})
	}
	return c.JSON(fiber.Map{"message": "Kehadiran updated"})
}

// @Security BearerAuth
// DeleteKehadiran godoc
// @Summary Delete kehadiran
// @Tags Kehadiran
// @Produce json
// @Param id path string true "Kehadiran ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kehadiran/{id} [delete]
func DeleteKehadiran(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = config.DB.Collection("kehadiran").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete kehadiran"})
	}
	return c.JSON(fiber.Map{"message": "Kehadiran deleted"})
}
