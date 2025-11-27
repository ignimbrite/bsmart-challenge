package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

var productSortOptions = map[string]string{
	"price_asc":  "price asc",
	"price_desc": "price desc",
	"name_asc":   "name asc",
	"name_desc":  "name desc",
	"newest":     "created_at desc",
	"oldest":     "created_at asc",
}

func (s *Server) listProducts(c *gin.Context) {
	var query ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "invalid query params")
		return
	}

	page, pageSize, _ := parsePagination(query.PaginationQuery)
	order := sanitizeSort(query.Sort, productSortOptions, "created_at desc")

	db := s.db.Model(&models.Product{}).Preload("Categories")

	if query.CategoryID > 0 {
		db = db.Joins("JOIN product_categories pc ON pc.product_id = products.id").Where("pc.category_id = ?", query.CategoryID)
	}

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
		respondError(c, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      products,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

func (s *Server) getProduct(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var product models.Product
	if err := s.db.Preload("Categories").First(&product, id).Error; err != nil {
		if errorsIs(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "product not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to fetch product")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (s *Server) createProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if len(req.CategoryIDs) > 0 {
		var categories []models.Category
		if err := s.db.Where("id IN ?", req.CategoryIDs).Find(&categories).Error; err != nil {
			respondError(c, http.StatusBadRequest, "invalid categories")
			return
		}
		if len(categories) != len(req.CategoryIDs) {
			respondError(c, http.StatusBadRequest, "some categories not found")
			return
		}
		product.Categories = categories
	}

	if err := s.db.Create(&product).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to create product")
		return
	}

	if err := s.recordHistory(product.ID, product.Price, product.Stock); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to record history")
		return
	}

	s.wsHub.Broadcast(NewWSMessage("product.created", product))

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

func (s *Server) updateProduct(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	var product models.Product
	if err := s.db.Preload("Categories").First(&product, id).Error; err != nil {
		if errorsIs(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "product not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to fetch product")
		return
	}

	originalPrice := product.Price
	originalStock := product.Stock

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if req.CategoryIDs != nil {
		var categories []models.Category
		if len(req.CategoryIDs) > 0 {
			if err := s.db.Where("id IN ?", req.CategoryIDs).Find(&categories).Error; err != nil {
				respondError(c, http.StatusBadRequest, "invalid categories")
				return
			}
			if len(categories) != len(req.CategoryIDs) {
				respondError(c, http.StatusBadRequest, "some categories not found")
				return
			}
		}
		if err := s.db.Model(&product).Association("Categories").Replace(categories); err != nil {
			respondError(c, http.StatusInternalServerError, "failed to update categories")
			return
		}
	}

	if err := s.db.Save(&product).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to update product")
		return
	}

	if product.Price != originalPrice || product.Stock != originalStock {
		if err := s.recordHistory(product.ID, product.Price, product.Stock); err != nil {
			respondError(c, http.StatusInternalServerError, "failed to record history")
			return
		}
	}

	s.wsHub.Broadcast(NewWSMessage("product.updated", product))

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (s *Server) deleteProduct(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	if err := s.db.Delete(&models.Product{}, id).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to delete product")
		return
	}

	s.wsHub.Broadcast(NewWSMessage("product.deleted", gin.H{"id": id}))

	c.Status(http.StatusNoContent)
}

func (s *Server) productHistory(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var query HistoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "invalid query params")
		return
	}

	db := s.db.Where("product_id = ?", id)

	if !query.Start.IsZero() {
		db = db.Where("changed_at >= ?", query.Start)
	}
	if !query.End.IsZero() {
		db = db.Where("changed_at <= ?", query.End)
	}

	var history []models.ProductHistory
	if err := db.Order("changed_at desc").Find(&history).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to fetch history")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}
