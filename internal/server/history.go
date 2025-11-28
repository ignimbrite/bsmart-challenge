package server

import (
	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

func (s *Server) recordHistory(db *gorm.DB, productID uint, price float64, stock int) error {
	entry := models.ProductHistory{
		ProductID: productID,
		Price:     price,
		Stock:     stock,
	}
	return db.Create(&entry).Error
}
