.PHONY: run dev test migrate reset

run:
	go run main.go

dev:
	air

test:
	go test -v ./model/entity/...

# Usage: make migrate-make name=migration_name
migrate-make:
	go run cmd/migration/main.go $(name)

reset:
	go run cmd/reset/main.go
