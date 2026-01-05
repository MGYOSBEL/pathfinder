# Services
MQTT_BROKER ?= hivemq

COMPOSE_FILES := \
		-f deploy/docker/docker-compose.yaml \
		-f deploy/docker/mqtt/$(MQTT_BROKER).yaml \

run-metrics-injector:
	go run cmd/metrics-injector/main.go

.PHONY: deploy
deploy:
	docker compose $(COMPOSE_FILES) up -d

.PHONY: stop
stop:
	docker compose $(COMPOSE_FILES) stop

.PHONY: down
down:
	docker compose $(COMPOSE_FILES) down

.PHONY: destroy
destroy:
	docker compose $(COMPOSE_FILES) down -v --remove-orphans

.PHONY: help
help:
	@echo "Available commands:"
	@echo "	deploy:		Start the project"
	@echo "	stop:		Stop the containers(keep them)"
	@echo "	down:		Delete the containers(keep data)"
	@echo "	destroy:	Delete containers and data"
