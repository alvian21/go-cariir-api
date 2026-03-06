package main

import (
	"go-cariir-api/database"
	"go-cariir-api/database/seed"
)

func main() {
	database.DatabaseInit()
	if database.DB != nil {
		seed.RunSeed()
	}
}
