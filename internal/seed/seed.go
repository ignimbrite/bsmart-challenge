package seed

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

func Run(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Category{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("seed: categories already present, skipping categories/products")
	} else {
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

		log.Println("seed: sample categories and products inserted")
	}

	if err := seedAdmin(db); err != nil {
		return err
	}

	return nil
}

func seedAdmin(db *gorm.DB) error {
	const email = "admin@bsmart.test"
	const password = "admin123"

	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:        email,
		PasswordHash: string(hash),
		Role:         "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	log.Printf("seed: admin user created email=%s password=%s", email, password)
	return nil
}
