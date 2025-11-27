package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/config"
)

type Server struct {
	cfg         config.Config
	db          *gorm.DB
	engine      *gin.Engine
	tokenSecret []byte
	tokenTTL    time.Duration
	wsHub       *Hub
}

func New(cfg config.Config, db *gorm.DB, tokenSecret []byte, tokenTTL time.Duration) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	hub := NewHub()
	go hub.Run()

	srv := &Server{
		cfg:         cfg,
		db:          db,
		engine:      engine,
		tokenSecret: tokenSecret,
		tokenTTL:    tokenTTL,
		wsHub:       hub,
	}

	srv.registerRoutes()

	return srv
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.engine.GET("/ws", s.handleWebSocket)

	api := s.engine.Group("/api")

	api.POST("/auth/login", s.login)

	api.GET("/products", s.listProducts)
	api.GET("/products/:id", s.getProduct)
	api.GET("/products/:id/history", s.productHistory)

	api.GET("/categories", s.listCategories)

	api.GET("/search", s.search)

	admin := api.Group("/")
	admin.Use(s.authMiddleware("admin"))
	admin.POST("/products", s.createProduct)
	admin.PUT("/products/:id", s.updateProduct)
	admin.DELETE("/products/:id", s.deleteProduct)

	admin.POST("/categories", s.createCategory)
	admin.PUT("/categories/:id", s.updateCategory)
	admin.DELETE("/categories/:id", s.deleteCategory)
}

func (s *Server) Run() error {
	address := fmt.Sprintf(":%s", s.cfg.HTTPPort)
	return s.engine.Run(address)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
