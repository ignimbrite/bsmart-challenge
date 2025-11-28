package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ignimbrite/bsmart-challenge/internal/config"
)

type Server struct {
	cfg            config.Config
	db             *gorm.DB
	engine         *gin.Engine
	tokenSecret    []byte
	tokenTTL       time.Duration
	wsHub          *Hub
	allowedOrigins []string
}

func New(cfg config.Config, db *gorm.DB, tokenSecret []byte, tokenTTL time.Duration) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	hub := NewHub()
	go hub.Run()

	srv := &Server{
		cfg:            cfg,
		db:             db,
		engine:         engine,
		tokenSecret:    tokenSecret,
		tokenTTL:       tokenTTL,
		wsHub:          hub,
		allowedOrigins: cfg.WSAllowed,
	}

	engine.Use(corsMiddleware(srv.allowedOrigins))
	srv.registerRoutes()

	return srv
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.engine.StaticFS("/web", gin.Dir("docs", false))
	s.engine.StaticFS("/docs", gin.Dir("docs", false))

	s.engine.GET("/ws", s.authMiddleware("admin", "client"), s.handleWebSocket)

	api := s.engine.Group("/api")

	api.POST("/auth/login", s.login)

	protected := api.Group("/")
	protected.Use(s.authMiddleware("admin", "client"))
	protected.GET("/products", s.listProducts)
	protected.GET("/products/:id", s.getProduct)
	protected.GET("/products/:id/history", s.productHistory)
	protected.GET("/categories", s.listCategories)
	protected.GET("/search", s.search)

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

func corsMiddleware(allowed []string) gin.HandlerFunc {
	normalized := make(map[string]struct{})
	for _, o := range allowed {
		if o == "*" {
			normalized["*"] = struct{}{}
			continue
		}
		o = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(o)), "/")
		if o != "" {
			normalized[o] = struct{}{}
		}
	}

	isAllowed := func(origin string) bool {
		if len(normalized) == 0 {
			return false
		}
		if _, ok := normalized["*"]; ok {
			return true
		}
		origin = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(origin)), "/")
		_, ok := normalized[origin]
		return ok
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowOrigin := isAllowed(origin)
		if allowOrigin {
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else if _, ok := normalized["*"]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
			c.Writer.Header().Set("Vary", "Origin")
		}
		if allowOrigin || c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Accept,ngrok-skip-browser-warning")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
