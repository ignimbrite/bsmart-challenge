package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ignimbrite/bsmart-challenge/internal/config"
)

type Server struct {
	cfg    config.Config
	engine *gin.Engine
}

func New(cfg config.Config) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	srv := &Server{
		cfg:    cfg,
		engine: engine,
	}

	srv.registerRoutes()

	return srv
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func (s *Server) Run() error {
	address := fmt.Sprintf(":%s", s.cfg.HTTPPort)
	return s.engine.Run(address)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
