.PHONY: migrate-up migrate-down compose-up compose-down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

compose-up:
	docker-compose -f deployments/docker/docker-compose.yml up -d

compose-down:
	docker-compose -f deployments/docker/docker-compose.yml down