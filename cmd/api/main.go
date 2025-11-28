package main

import (
	"log"
	"time"

	"github.com/ignimbrite/bsmart-challenge/internal/config"
	appdb "github.com/ignimbrite/bsmart-challenge/internal/db"
	"github.com/ignimbrite/bsmart-challenge/internal/models"
	"github.com/ignimbrite/bsmart-challenge/internal/seed"
	"github.com/ignimbrite/bsmart-challenge/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := appdb.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql DB: %v", err)
	}
	defer sqlDB.Close()

	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	if cfg.AppEnv == "development" || cfg.SeedOnStart {
		if err := seed.Run(db); err != nil {
			log.Fatalf("seed failed: %v", err)
		}
	}

	tokenTTL, err := time.ParseDuration(cfg.JWTExpiration)
	if err != nil {
		log.Fatalf("invalid JWT_EXPIRATION: %v", err)
	}

	srv := server.New(cfg, db, []byte(cfg.JWTSecret), tokenTTL)

	log.Printf("starting api server on :%s (env: %s)", cfg.HTTPPort, cfg.AppEnv)

	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
