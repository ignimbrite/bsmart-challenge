package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ignimbrite/bsmart-challenge/internal/config"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	conn, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Printf("warn: database ping failed: %v", err)
	}

	return conn, nil
}
