DOCKER_COMPOSE_FILE ?= ./dev-env-setup/docker-compose.yml

help: 
		@echo "make env-up to build up environment by docker-compose"
		@echo "make env-down to stop/close environment by docker-compose"

env-up:
		@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

env-down: 
		@docker-compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: help env-up env-down
