package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	log.Println("Connecting DB...")
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	log.Println("DB successfully connected")
	return database, nil
}
