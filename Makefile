run-metrics-injector:
	go run cmd/metrics-injector/main.go

.PHONY: deploy
deploy:
	docker compose \
		-f deploy/docker/docker-compose.yaml \
		-f deploy/docker/mqtt/hivemq.yaml \
		up -d

drain:
	docker compose \
		-f deploy/docker/docker-compose.yaml \
		-f deploy/docker/mqtt/hivemq.yaml \
		down

