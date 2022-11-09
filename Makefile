BUILD_VERSION := $(shell git rev-parse --short HEAD)

start-local-db:
	docker-compose up -d

run-webserver: start-local-db
	LOG_LEVEL=debug RUNNING_LOCALLY=true database_host_name=localhost database_port=5432 database_name=user-mgmt-db database_username=postgres database_password=test database_ssl_mode=disable go run -ldflags="-X 'main.BuildVersion=$(BUILD_VERSION)'" cmd/main.go

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
