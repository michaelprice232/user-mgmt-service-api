BUILD_VERSION := $(shell git rev-parse --short HEAD)

run:
	docker-compose up -d

stop-database-and-delete-volume:
	docker-compose stop
	docker-compose rm --force
	docker volume rm user-mgmt-service-api_db-data --force

test-endpoints:
	./scripts/client.sh

prune-docker:
	docker system prune --all --volumes --force

test:
	go test user-mgmt-service-api/internal/api

version:
	go run -ldflags="-X 'main.BuildVersion=$(BUILD_VERSION)'" cmd/main.go --version
