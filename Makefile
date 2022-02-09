DEV_ENV_SETUP_FOLDER ?= ./dev-env-setup
DOCKER_COMPOSE_FILE ?= $(DEV_ENV_SETUP_FOLDER)/docker-compose.yml

help: 
		@echo "make env-up to build up environment by docker-compose"
		@echo "make env-down to stop/close environment by docker-compose"
		@echo "make start to start go-api server"
		@echo "make test to run the functional test"
		@echo "make create-migrations to create migration files(up&down). Enter migration_name from standard input."
		@echo "make upgrade-schema to the latest version."
		@echo "make downgrade-schema to remove all the migrations. Use with CAUTION."
		@echo "make goto-schema to go to the specific version. Enter schema_version from standard input. "

env-up:
		@cp ./membership_user.sql $(DEV_ENV_SETUP_FOLDER)/mysql/initdb.sql
		@docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d

env-down: 
		@docker-compose -f $(DOCKER_COMPOSE_FILE) down

start:
		@go run main.go

test: 
		@go test $$(glide novendor) 

DB_USER ?= test_membership
DB_PASSWORD ?= test_membership
DB_NAME ?= test_membership
DB_ADDRESS ?= 127.0.0.1
DB_PORT ?= 3306

# Migration
MIGRATION_NAME ?= $(shell read -p "Migration name: " migration_name; echo $$migration_name)
MIGRATION_EXT ?= "sql"
MIGRATION_DIR ?= "migrations"
SCHEMA_VERSION ?= $(shell read -p "Schema version: " schema_version; echo $$schema_version)

DB_CONN = mysql://$(DB_USER):$(DB_PASSWORD)@tcp\($(DB_HOST):$(DB_PORT)\)/$(DB_NAME)


################################################################################
# Database Migration
# ################################################################################

create-migrations: check-migrate
		@migrate create -ext $(MIGRATION_EXT) -dir $(MIGRATION_DIR) -seq $(MIGRATION_NAME)

upgrade-schema: check-migrate
		@echo Upgrade UP schema or to latest version
		migrate -database $(DB_CONN) -path $(MIGRATION_DIR) up $(UP)

downgrade-schema: check-migrate
	        @echo CAUTION: Remove DOWN schema or all the schema
		@migrate -database $(DB_CONN) -path $(MIGRATION_DIR) down $(DOWN)

goto-schema: check-migrate
		@echo "Roll up or down to the specified version"
		@migrate -database $(DB_CONN) -path $(MIGRATION_DIR) goto $(SCHEMA_VERSION)

force-schema: check-migrate
		@echo "Force to specified version without actually migrate"
		@migrate -database $(DB_CONN) -path $(MIGRATION_DIR) force $(SCHEMA_VERSION)

check-version: check-migrate
		@echo "Check current migrate version"
		@migrate -database $(DB_CONN) -path $(MIGRATION_DIR) version

check-migrate:
		@printf "Check if migrate CLI (https://github.com/golang-migrate/migrate/tree/master/cli) is installed."
		@type migrate > /dev/null
		@echo ....OK

.PHONY: help env-up env-down start test
