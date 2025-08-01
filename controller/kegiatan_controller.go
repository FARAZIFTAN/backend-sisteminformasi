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

var kegiatanValidate = validator.New()

// Tambahkan struct response
type KegiatanResponse struct {
	ID              string `json:"id"`
	Judul           string `json:"judul"`
	Deskripsi       string `json:"deskripsi"`
	Tanggal         string `json:"tanggal"`
	Lokasi          string `json:"lokasi"`
	Kategori        string `json:"kategori"`
	MaxParticipants int    `json:"maxParticipants"`
	DokumentasiURL  string `json:"dokumentasi_url"`
	CreatedBy       string `json:"created_by"`
}

func safeObjectIDHex(id interface{}) string {
	if oid, ok := id.(primitive.ObjectID); ok {
		return oid.Hex()
	}
	if s, ok := id.(string); ok {
		return s
	}
	return ""
}

// GetKegiatan godoc
// @Summary Get all kegiatan
// @Tags Kegiatan
// @Produce json
// @Success 200 {array} model.Kegiatan
// @Failure 500 {object} map[string]interface{}
// @Router /kegiatan [get]
// @Security BearerAuth
func GetKegiatan(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := config.DB.Collection("kegiatan").Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch kegiatan"})
	}
	var rawKegiatans []bson.M
	if err := cursor.All(ctx, &rawKegiatans); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode kegiatan"})
	}
	var responses []KegiatanResponse
	for _, k := range rawKegiatans {
		idStr := ""
		if oid, ok := k["_id"].(primitive.ObjectID); ok {
			idStr = oid.Hex()
		} else if s, ok := k["_id"].(string); ok {
			idStr = s
		}
		responses = append(responses, KegiatanResponse{
			ID:              idStr,
			Judul:           getStringFromMap(k, "judul"),
			Deskripsi:       getStringFromMap(k, "deskripsi"),
			Tanggal:         getStringFromMap(k, "tanggal"),
			Lokasi:          getStringFromMap(k, "lokasi"),
			Kategori:        getStringFromMap(k, "kategori"),
			MaxParticipants: getIntFromMap(k, "maxParticipants"),
			DokumentasiURL:  getStringFromMap(k, "dokumentasi_url"),
			CreatedBy:       getStringFromMap(k, "created_by"),
		})
	}
	return c.JSON(responses)
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getIntFromMap(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int:
			return val
		case int32:
			return int(val)
		case int64:
			return int(val)
		case float64:
			return int(val)
		}
	}
	return 0
}

// GetKegiatanByID godoc
// @Summary Get kegiatan by ID
// @Tags Kegiatan
// @Produce json
// @Param id path string true "Kegiatan ID"
// @Success 200 {object} model.Kegiatan
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kegiatan/{id} [get]
// @Security BearerAuth
func GetKegiatanByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var kegiatan model.Kegiatan
	err = config.DB.Collection("kegiatan").FindOne(ctx, bson.M{"_id": objID}).Decode(&kegiatan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Kegiatan not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch kegiatan"})
	}
	response := KegiatanResponse{
		ID:              safeObjectIDHex(kegiatan.ID),
		Judul:           kegiatan.Judul,
		Deskripsi:       kegiatan.Deskripsi,
		Tanggal:         kegiatan.Tanggal,
		Lokasi:          kegiatan.Lokasi,
		Kategori:        kegiatan.Kategori,
		MaxParticipants: kegiatan.MaxParticipants,
		DokumentasiURL:  kegiatan.DokumentasiURL,
		CreatedBy:       kegiatan.CreatedBy,
	}
	return c.JSON(response)
}

// CreateKegiatan godoc
// @Summary Create kegiatan
// @Tags Kegiatan
// @Accept json
// @Produce json
// @Param kegiatan body model.Kegiatan true "Kegiatan Data"
// @Success 201 {object} model.Kegiatan
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kegiatan [post]
// @Security BearerAuth
func CreateKegiatan(c *fiber.Ctx) error {
	var kegiatan model.Kegiatan
	if err := c.BodyParser(&kegiatan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := kegiatanValidate.Struct(kegiatan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := config.DB.Collection("kegiatan").InsertOne(ctx, kegiatan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create kegiatan"})
	}
	kegiatan.ID = res.InsertedID.(primitive.ObjectID)
	response := KegiatanResponse{
		ID:              kegiatan.ID.Hex(),
		Judul:           kegiatan.Judul,
		Deskripsi:       kegiatan.Deskripsi,
		Tanggal:         kegiatan.Tanggal,
		Lokasi:          kegiatan.Lokasi,
		Kategori:        kegiatan.Kategori,
		MaxParticipants: kegiatan.MaxParticipants,
		DokumentasiURL:  kegiatan.DokumentasiURL,
		CreatedBy:       kegiatan.CreatedBy,
	}
	return c.Status(201).JSON(response)
}

// UpdateKegiatan godoc
// @Summary Update kegiatan
// @Tags Kegiatan
// @Accept json
// @Produce json
// @Param id path string true "Kegiatan ID"
// @Param kegiatan body model.Kegiatan true "Kegiatan Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kegiatan/{id} [put]
// @Security BearerAuth
func UpdateKegiatan(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var kegiatan model.Kegiatan
	if err := c.BodyParser(&kegiatan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := kegiatanValidate.Struct(kegiatan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{
		"judul":           kegiatan.Judul,
		"deskripsi":       kegiatan.Deskripsi,
		"tanggal":         kegiatan.Tanggal,
		"lokasi":          kegiatan.Lokasi,
		"kategori":        kegiatan.Kategori,
		"maxParticipants": kegiatan.MaxParticipants,
		"dokumentasi_url": kegiatan.DokumentasiURL,
		"created_by":      kegiatan.CreatedBy,
	}
	_, err = config.DB.Collection("kegiatan").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update kegiatan"})
	}
	return c.JSON(fiber.Map{"message": "Kegiatan updated"})
}

// DeleteKegiatan godoc
// @Summary Delete kegiatan
// @Tags Kegiatan
// @Produce json
// @Param id path string true "Kegiatan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /kegiatan/{id} [delete]
// @Security BearerAuth
func DeleteKegiatan(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = config.DB.Collection("kegiatan").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete kegiatan"})
	}
	return c.JSON(fiber.Map{"message": "Kegiatan deleted"})
}
