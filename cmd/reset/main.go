package main

import (
	"fmt"
	"go-cariir-api/database"
	"go-cariir-api/database/seed"
	"log"

	"github.com/pressly/goose/v3"
	_ "github.com/pressly/goose/v3/database"
)

func main() {
	database.DatabaseInit()

	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("Could not get sql.DB from gorm: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Could set goose dialect: %v", err)
	}

	migrationDir := "database/migration/sql"

	fmt.Println("Resetting database (Down all)...")
	if err := goose.DownTo(sqlDB, migrationDir, 0); err != nil {
		log.Fatalf("Could not run down migrations: %v", err)
	}

	fmt.Println("Running up migrations...")
	if err := goose.Up(sqlDB, migrationDir); err != nil {
		log.Fatalf("Could not run up migrations: %v", err)
	}

	fmt.Println("Running seeds...")
	seed.RunSeed()

	fmt.Println("Database reset and migration complete")
}
