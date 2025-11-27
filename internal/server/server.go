package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/config"
	"github.com/ignimbrite/bsmart-challenge/internal/models"
)

type Server struct {
	cfg    config.Config
	db     *gorm.DB
	engine *gin.Engine
}

func New(cfg config.Config, db *gorm.DB) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	srv := &Server{
		cfg:    cfg,
		db:     db,
		engine: engine,
	}

	srv.registerRoutes()

	return srv
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := s.engine.Group("/api")
	api.GET("/products", s.listProducts)
}

func (s *Server) listProducts(c *gin.Context) {
	var products []models.Product
	if err := s.db.Preload("Categories").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (s *Server) Run() error {
	address := fmt.Sprintf(":%s", s.cfg.HTTPPort)
	return s.engine.Run(address)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
