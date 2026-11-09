package main

import (
	"log"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/config"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/database"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Open(cfg.DatabaseDSN, cfg.Environment)
	if err != nil {
		log.Fatal(err)
	}
	if err := domain.Migrate(db); err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("database migration completed")
}
