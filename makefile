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


# Example usage:
# make migration name=initial_tables
# make migrate-up
# make migrate-down N=1
# make migrate-force V=0
# make migrate-reset