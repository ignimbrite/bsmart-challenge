package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

func (s *Server) search(c *gin.Context) {
	var query SearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "invalid query params")
		return
	}

	switch query.Type {
	case "product":
		s.searchProducts(c, query)
	case "category":
		s.searchCategories(c, query)
	default:
		respondError(c, http.StatusBadRequest, "unsupported search type")
	}
}

func (s *Server) searchProducts(c *gin.Context, query SearchQuery) {
	page, pageSize, _ := parsePagination(query.PaginationQuery)
	order := sanitizeSort(query.Sort, productSortOptions, "created_at desc")

	db := s.db.Model(&models.Product{}).Preload("Categories")

	if query.Query != "" {
		like := "%" + query.Query + "%"
		db = db.Where("products.name ILIKE ? OR products.description ILIKE ?", like, like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to count products")
		return
	}

	var products []models.Product
	if err := db.Order(order).Limit(pageSize).Offset((page - 1) * pageSize).Find(&products).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to search products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      products,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

func (s *Server) searchCategories(c *gin.Context, query SearchQuery) {
	page, pageSize, _ := parsePagination(query.PaginationQuery)
	order := sanitizeSort(query.Sort, categorySortOptions, "created_at desc")

	db := s.db.Model(&models.Category{})
	if query.Query != "" {
		like := "%" + query.Query + "%"
		db = db.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to count categories")
		return
	}

	var categories []models.Category
	if err := db.Order(order).Limit(pageSize).Offset((page - 1) * pageSize).Find(&categories).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to search categories")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      categories,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}
