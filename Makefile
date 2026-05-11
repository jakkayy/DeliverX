.PHONY: help infra-up infra-down infra-logs backend-build backend-up backend-down logs db-migrate mobile-run tidy

COMPOSE_DIR = infra/docker
COMPOSE_FILE = $(COMPOSE_DIR)/docker-compose.yml
ENV_FILE = $(COMPOSE_DIR)/.env

help:
	@echo "Available commands:"
	@echo "  make infra-up        Start infrastructure (PostgreSQL, Redis, Kafka)"
	@echo "  make infra-down      Stop infrastructure"
	@echo "  make infra-logs      Tail infrastructure logs"
	@echo "  make backend-build   Build all backend Docker images"
	@echo "  make backend-up      Start all services"
	@echo "  make backend-down    Stop all services"
	@echo "  make logs            Tail all service logs"
	@echo "  make db-migrate      Run database schema"
	@echo "  make mobile-run      Run Flutter app"
	@echo "  make clean           Remove all containers and volumes"
	@echo "  make tidy            Run go mod tidy for all backend services"

infra-up:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d postgres redis zookeeper kafka kafka-ui pgadmin
	@echo "Infrastructure is up!"
	@echo "  pgAdmin   -> http://localhost:5050"
	@echo "  Kafka UI  -> http://localhost:8090"

infra-down:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) stop postgres redis zookeeper kafka kafka-ui pgadmin

infra-logs:
	docker compose -f $(COMPOSE_FILE) logs -f postgres redis kafka

backend-build:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) build

backend-up:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d
	@echo "All services are up!"
	@echo "  API Gateway   -> http://localhost:8080"
	@echo "  Auth Service  -> http://localhost:8081"
	@echo "  User Service  -> http://localhost:8082"
	@echo "  Order Service -> http://localhost:8083"
	@echo "  Tracking      -> http://localhost:8084"
	@echo "  Payment       -> http://localhost:8085"

backend-down:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) stop api-gateway auth-service user-service order-service tracking-service payment-service notification-service

logs:
	docker compose -f $(COMPOSE_FILE) logs -f

db-migrate:
	docker exec -i grab-postgres psql -U grab_user -d grab_db < database/schema.sql

mobile-run:
	cd mobile && flutter run

tidy:
	@for dir in auth-service user-service order-service tracking-service payment-service notification-service api-gateway; do \
		echo "==> $$dir"; \
		cd backend/$$dir && go mod tidy && cd ../..; \
	done

clean:
	docker compose -f $(COMPOSE_FILE) down -v --remove-orphans
