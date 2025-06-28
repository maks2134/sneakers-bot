package utils

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"os"
	"snakers-bot/internal/usecases"
)

func SeedProducts(db *gorm.DB) {
	var count int64
	db.Model(&usecases.Product{}).Count(&count)

	if count > 0 {
		log.Println("Database already seeded. Skipping.")
		return
	}

	log.Println("Seeding database with initial products...")

	file, err := os.ReadFile("data/products.json")
	if err != nil {
		log.Fatalf("Failed to read seed file: %v", err)
	}

	var products []usecases.Product
	if err := json.Unmarshal(file, &products); err != nil {
		log.Fatalf("Failed to unmarshal seed data: %v", err)
	}

	if err := db.Create(&products).Error; err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	log.Printf("Successfully seeded %d products.", len(products))
}
