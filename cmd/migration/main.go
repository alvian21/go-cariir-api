package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migration/main.go <migration_name>")
		return
	}

	migrationName := os.Args[1]
	migrationDir := "database/migration/sql"

	// Ensure directory exists
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		err := os.MkdirAll(migrationDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create migration directory: %v", err)
		}
	}

	// Get next sequence number
	nextSeq := getNextSequence(migrationDir)
	version := fmt.Sprintf("%06d", nextSeq)

	fileName := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.sql", version, migrationName))

	createMigrationFile(fileName)

	fmt.Printf("Migration created:\n - %s\n", fileName)
}

func getNextSequence(dir string) int {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 1
	}

	var sequences []int
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, ".sql") {
			parts := strings.Split(name, "_")
			if len(parts) > 0 {
				seq, err := strconv.Atoi(parts[0])
				if err == nil {
					sequences = append(sequences, seq)
				}
			}
		}
	}

	if len(sequences) == 0 {
		return 1
	}

	sort.Ints(sequences)
	return sequences[len(sequences)-1] + 1
}

func createMigrationFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", path, err)
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	content := fmt.Sprintf("-- +goose Up\n-- Migration: %s\n-- Created at: %s\n\n\n-- +goose Down\n", filepath.Base(path), timestamp)
	_, _ = file.WriteString(content)
}
