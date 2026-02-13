include .env
export


# Create a new migration file
migration:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Run migration up
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

# Run down N migrations (default 1)
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down $(N)


# Force database to a specific version (helpful for dirty state)
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(V)

# Reset DB (DEV only: force version 0 then run all migrations)
migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force 0
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

migration-version:
	@echo "Checking current migration version..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose version

migrate-drop-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down -all

# make migrate-drop-all
# make migrate-up


migrate-reset-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down -all
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up


# Example usage:
# make migration name=initial_tables
# make migrate-up
# make migrate-down N=1
# make migrate-force V=0
# make migrate-reset


# Note:
# Makefiles require TAB for indentation. Spaces will break it silently.