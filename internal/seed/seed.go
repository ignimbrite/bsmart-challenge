package seed

import (
	"log"

	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

func Run(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Category{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("seed: skipping, data already present")
		return nil
	}

	categories := []models.Category{
		{Name: "Electronics", Description: "Devices and gadgets"},
		{Name: "Office", Description: "Office supplies and equipment"},
		{Name: "Groceries", Description: "Food and household items"},
	}

	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	products := []models.Product{
		{
			Name:        "Laptop",
			Description: "Lightweight laptop",
			Price:       1200.00,
			Stock:       10,
			Categories:  []models.Category{categories[0]},
		},
		{
			Name:        "Mechanical Keyboard",
			Description: "RGB mechanical keyboard",
			Price:       150.00,
			Stock:       30,
			Categories:  []models.Category{categories[0], categories[1]},
		},
		{
			Name:        "Coffee Beans",
			Description: "500g specialty coffee",
			Price:       18.50,
			Stock:       50,
			Categories:  []models.Category{categories[2]},
		},
	}

	if err := db.Create(&products).Error; err != nil {
		return err
	}

	log.Println("seed: sample data inserted")
	return nil
}
