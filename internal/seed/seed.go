package seed

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

func Run(db *gorm.DB) error {
	seed := time.Now().UnixNano()
	gofakeit.Seed(seed)
	rand.Seed(seed)

	var count int64
	if err := db.Model(&models.Category{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("seed: categories already present, skipping categories/products")
	} else {
		categories := buildCategories(10)
		if err := db.Create(&categories).Error; err != nil {
			return err
		}

		products := buildProducts(100, categories)
		if err := db.Create(&products).Error; err != nil {
			return err
		}

		log.Printf("seed: inserted %d categories and %d products", len(categories), len(products))
	}

	if err := seedAdmin(db); err != nil {
		return err
	}

	if err := seedClient(db); err != nil {
		return err
	}

	return nil
}

func buildCategories(n int) []models.Category {
	names := make(map[string]struct{})
	var categories []models.Category

	for len(categories) < n {
		name := randomCategory()
		if _, exists := names[name]; exists {
			continue
		}
		names[name] = struct{}{}
		categories = append(categories, models.Category{
			Name:        name,
			Description: fakeDescription(6, 8),
		})
	}

	return categories
}

func buildProducts(n int, categories []models.Category) []models.Product {
	var products []models.Product
	names := make(map[string]struct{})

	for len(products) < n {
		name := fakeProductName()
		if _, exists := names[name]; exists {
			continue
		}
		names[name] = struct{}{}

		product := models.Product{
			Name:        name,
			Description: fakeDescription(6, 12),
			Price:       gofakeit.Price(5, 750),
			Stock:       gofakeit.Number(0, 500),
		}

		for _, idx := range randomCategoryIndexes(len(categories)) {
			product.Categories = append(product.Categories, categories[idx])
		}

		products = append(products, product)
	}

	return products
}

func randomCategoryIndexes(max int) []int {
	if max == 0 {
		return nil
	}
	count := rand.Intn(3) + 1 // 1-3 categorías por producto
	if count > max {
		count = max
	}

	indexes := make([]int, 0, count)
	seen := make(map[int]struct{})
	for len(indexes) < count {
		idx := rand.Intn(max)
		if _, ok := seen[idx]; ok {
			continue
		}
		seen[idx] = struct{}{}
		indexes = append(indexes, idx)
	}
	return indexes
}

func fakeProductName() string {
	nouns := []string{"Headphones", "Coffee Maker", "Lamp", "Keyboard", "Monitor", "Battery", "Backpack", "Router", "Camera", "Speaker", "Chair", "Desk", "Pen", "Notebook", "Drone", "Watch", "Microphone"}
	adjectives := []string{"Eco", "Premium", "Classic", "Modern", "Compact", "Smart", "Portable", "Deluxe", "Fast", "Silent"}

	noun := gofakeit.RandomString(nouns)
	adj := gofakeit.RandomString(adjectives)
	return noun + " " + adj
}

func fakeDescription(minWords, maxWords int) string {
	words := []string{
		"quality", "warranty", "design", "ergonomic", "lightweight", "durable", "practical", "versatile",
		"daily", "office", "home", "work", "comfort", "performance", "fast", "quiet",
		"connectivity", "wireless", "battery", "rechargeable", "materials", "premium", "long-lasting",
	}
	count := minWords
	if maxWords > minWords {
		count = minWords + rand.Intn(maxWords-minWords+1)
	}
	if count > len(words) {
		count = len(words)
	}
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
	selected := words[:count]
	s := strings.Join(selected, " ")
	return strings.Title(s) + "."
}

func randomCategory() string {
	categories := []string{
		"Electronics",
		"Home",
		"Office",
		"Sports",
		"Fashion",
		"Food",
		"Garden",
		"Automotive",
		"Beauty",
		"Pets",
		"Toys",
		"Health",
	}
	return gofakeit.RandomString(categories)
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

func seedClient(db *gorm.DB) error {
	const email = "client@bsmart.test"
	const password = "client123"

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
		Role:         "client",
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	log.Printf("seed: client user created email=%s password=%s", email, password)
	return nil
}
