package migration

import (
	"embed"
	"fmt"
	"go-cariir-api/database"
	"log"

	"github.com/pressly/goose/v3"
	_ "github.com/pressly/goose/v3/database"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

func RunMigration() {
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("Could not get sql.DB from gorm: %v", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Could not set goose dialect: %v", err)
	}

	if err := goose.Up(sqlDB, "sql"); err != nil {
		log.Fatalf("Could not run up migrations: %v", err)
	}

	fmt.Println("Migration complete")
}
