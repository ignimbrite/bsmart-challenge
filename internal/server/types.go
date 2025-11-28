package server

import "time"

type PaginationQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Sort     string `form:"sort"`
	Query    string `form:"q"`
}

type ProductQuery struct {
	PaginationQuery
	CategoryID uint `form:"category_id"`
}

type SearchQuery struct {
	Type string `form:"type" binding:"required,oneof=product category"`
	PaginationQuery
}

type CategoryQuery struct {
	Sort  string `form:"sort"`
	Query string `form:"q"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Description string  `json:"description" binding:"omitempty,max=2000"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	CategoryIDs []uint  `json:"category_ids" binding:"required,dive,gt=0"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=2,max=255"`
	Description *string  `json:"description" binding:"omitempty,max=2000"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
	CategoryIDs []uint   `json:"category_ids" binding:"omitempty,dive,gt=0"`
}

type HistoryQuery struct {
	Start time.Time `form:"start" time_format:"2006-01-02" time_utc:"1"`
	End   time.Time `form:"end" time_format:"2006-01-02" time_utc:"1"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
