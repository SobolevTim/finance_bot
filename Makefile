.PHONY: migrate-up migrate-down up down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

up:
	docker-compose -f deployments/docker/docker-compose.yml up -d

down:
	docker-compose -f deployments/docker/docker-compose.yml down
	docker rmi test-server:latest