package main

import (
	"log"

	"github.com/ignimbrite/bsmart-challenge/internal/config"
	"github.com/ignimbrite/bsmart-challenge/internal/server"
)

func main() {
	cfg := config.Load()

	srv := server.New(cfg)

	log.Printf("starting api server on :%s (env: %s)", cfg.HTTPPort, cfg.AppEnv)

	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
