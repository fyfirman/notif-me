package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cron "notif-me/services/cron"
)

func doMigration(db *gorm.DB) {
	db.AutoMigrate(&cron.MangaUpdate{})
}

func ConnectDB() (*gorm.DB, error) {
	log.Info().Msg("Connecting DB...")
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := os.Getenv("POSTGRES_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	log.Info().Msg("DB successfully connected")
	doMigration(database)
	return database, nil
}
