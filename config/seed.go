package config

import (
	"context"
	"log"
	"os"
	"time"

	"backend-sisteminformasi/model"
	"backend-sisteminformasi/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedAdminUser membuat user admin default jika belum ada di database
func SeedAdminUser(db *mongo.Database) {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@ukm.com"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}
	adminNama := os.Getenv("ADMIN_NAMA")
	if adminNama == "" {
		adminNama = "Admin UKM"
	}
	adminUKM := os.Getenv("ADMIN_UKM")
	if adminUKM == "" {
		adminUKM = "Semua UKM"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cek apakah admin sudah ada
	var existing model.User
	err := db.Collection("users").FindOne(ctx, bson.M{"email": adminEmail}).Decode(&existing)
	if err == nil {
		return // Admin sudah ada
	}

	hash, err := utils.HashPassword(adminPassword)
	if err != nil {
		log.Println("Gagal hash password admin:", err)
		return
	}

	admin := model.User{
		Nama:     adminNama,
		Email:    adminEmail,
		Password: hash,
		Role:     "admin",
		UKM:      adminUKM,
	}
	_, err = db.Collection("users").InsertOne(ctx, admin)
	if err != nil {
		log.Println("Gagal insert admin user:", err)
	} else {
		log.Println("Admin user berhasil dibuat (", adminEmail, ")")
	}
}
