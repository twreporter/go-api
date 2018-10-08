DEV_ENV_SETUP_FOLDER ?= ./dev-env-setup
DOCKER_COMPOSE_FILE ?= $(DEV_ENV_SETUP_FOLDER)/docker-compose.yml

help: 
		@echo "make env-up to build up environment by docker-compose"
		@echo "make env-down to stop/close environment by docker-compose"
		@echo "make start to start go-api server"
		@echo "make test to run the functional test"

env-up:
		@cp ./membership_user.sql $(DEV_ENV_SETUP_FOLDER)/mysql/initdb.sql
		@docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d

env-down: 
		@docker-compose -f $(DOCKER_COMPOSE_FILE) down

start:
		@go run main.go

test: 
		@go test $$(glide novendor) 

.PHONY: help env-up env-down start test
