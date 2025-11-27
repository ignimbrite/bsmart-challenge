package server

import (
	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

func (s *Server) recordHistory(productID uint, price float64, stock int) error {
	entry := models.ProductHistory{
		ProductID: productID,
		Price:     price,
		Stock:     stock,
	}
	return s.db.Create(&entry).Error
}
