package main

import (
	"go-cariir-api/database"
	"go-cariir-api/database/migration"
	"go-cariir-api/database/seed"
	"go-cariir-api/route"

	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// INITIAL DATABASE
	database.DatabaseInit()
	// RUN MIGRATION
	migration.RunMigration()
	// RUN SEED
	seed.RunSeed()

	app := fiber.New(fiber.Config{
		ReadBufferSize: 16384,
	})

	// INITIAL ROUTE
	route.RouteInit(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
