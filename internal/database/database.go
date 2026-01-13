package database

import (
	"fmt"
	"log"

	"github.com/igorfazlyev/dm/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	logLevel := logger.Silent
	if cfg.Server.Environment == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&Patient{},
		&Study{},
		&PlanVersion{},
		&PlanItem{},
		&Clinic{},
		&PriceListItem{},
		&OfferRequest{},
		&Offer{},
		&Order{},
		&Slot{},
		&JobQueue{},
		&AuditLog{},
	)
}
