package models

import "time"

type Product struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `gorm:"size:255;not null;index:idx_products_name,sort:asc"`
	Description string     `gorm:"type:text"`
	Price       float64    `gorm:"type:numeric(12,2);not null"`
	Stock       int        `gorm:"not null;default:0;index"`
	Categories  []Category `gorm:"many2many:product_categories;constraint:OnDelete:CASCADE"`
	History     []ProductHistory `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
}

type Category struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null;uniqueIndex"`
	Description string    `gorm:"type:text"`
	Products    []Product `gorm:"many2many:product_categories;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
}

type ProductCategory struct {
	ProductID  uint `gorm:"primaryKey"`
	CategoryID uint `gorm:"primaryKey"`
}

type ProductHistory struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID uint      `gorm:"not null;index"`
	Price     float64   `gorm:"type:numeric(12,2);not null"`
	Stock     int       `gorm:"not null"`
	ChangedAt time.Time `gorm:"autoCreateTime"`
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"size:50;not null;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func AutoMigrate(db GormMigrator) error {
	return db.AutoMigrate(&Category{}, &Product{}, &ProductCategory{}, &ProductHistory{}, &User{})
}

type GormMigrator interface {
	AutoMigrate(dst ...interface{}) error
}
