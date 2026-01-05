# Services
MQTT_BROKER ?= hivemq
BROKER ?= rabbitmq

COMPOSE_FILES := \
		-f deploy/docker/docker-compose.yaml \
		-f deploy/docker/mqtt/$(MQTT_BROKER).yaml \
		-f deploy/docker/broker/$(BROKER).yaml \

# Validations
VALID_MQTT_BROKERS := hivemq vernemq
VALID_BROKERS := rabbitmq

ifneq ($(filter $(MQTT_BROKER),$(VALID_MQTT_BROKERS)),$(MQTT_BROKER))
$(error Invalid MQTT_BROKER '$(MQTT_BROKER)'. Valid: $(VALID_MQTT_BROKERS))
endif

ifneq ($(filter $(BROKER),$(VALID_BROKERS)),$(BROKER))
$(error Invalid BROKER '$(BROKER)'. Valid: $(VALID_BROKERS))
endif

.PHONY: deploy
deploy:
	docker compose $(COMPOSE_FILES) up -d

.PHONY: config
config:
	docker compose $(COMPOSE_FILES) config

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
	@echo "	config:		Parse the project files"
	@echo "	stop:		Stop the containers(keep them)"
	@echo "	down:		Delete the containers(keep data)"
	@echo "	destroy:	Delete containers and data"
