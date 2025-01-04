BUILD_VERSION := $(shell git rev-parse --short HEAD)

run:
	HOSTPORT=8080 docker-compose up -d

down:
	docker-compose down

cleanup:
	docker-compose down
	docker-compose rm --force
	docker volume rm user-mgmt-service-api_db-data --force

test-endpoints:
	./scripts/client.sh

prune-docker:
	docker system prune --all --volumes --force

unit-tests:
	go test ./...

int-tests:
	go test -tags=integration ./tests/integration

version:
	go run -ldflags="-X 'main.BuildVersion=$(BUILD_VERSION)'" cmd/main.go --version
