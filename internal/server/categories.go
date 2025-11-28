package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

var categorySortOptions = map[string]string{
	"name_asc":  "name asc",
	"name_desc": "name desc",
	"newest":    "created_at desc",
	"oldest":    "created_at asc",
}

func (s *Server) listCategories(c *gin.Context) {
	var query CategoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "invalid query params")
		return
	}

	order := sanitizeSort(query.Sort, categorySortOptions, "created_at desc")

	db := s.db.Model(&models.Category{})

	if query.Query != "" {
		like := "%" + query.Query + "%"
		db = db.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	var categories []models.Category
	if err := db.Order(order).Find(&categories).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"total": len(categories),
	})
}

func (s *Server) createCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.db.Create(&category).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to create category")
		return
	}

	s.wsHub.Broadcast(NewWSMessage("category.created", category))

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

func (s *Server) updateCategory(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if errorsIs(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "category not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to fetch category")
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.db.Save(&category).Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to update category")
		return
	}

	s.wsHub.Broadcast(NewWSMessage("category.updated", category))

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (s *Server) deleteCategory(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	res := s.db.Delete(&models.Category{}, id)
	if err := res.Error; err != nil {
		respondError(c, http.StatusInternalServerError, "failed to delete category")
		return
	}

	if res.RowsAffected == 0 {
		respondError(c, http.StatusNotFound, "category not found")
		return
	}

	s.wsHub.Broadcast(NewWSMessage("category.deleted", gin.H{"id": id}))

	c.Status(http.StatusNoContent)
}
