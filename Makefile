BUILD_VERSION := $(shell git rev-parse --short HEAD)

run:
	HOSTPORT=8080 BUILD_VERSION=$(BUILD_VERSION) docker-compose up -d --build

down:
	HOSTPORT=8080 docker-compose down --volumes

unit-tests:
	go test -v ./...

int-tests:
	go test -tags=integration -count=1 -v ./tests/integration

e2e-tests:
	go test -tags=e2e -count=1 -v -timeout 30m ./tests/e2e

version:
	go run -ldflags="-X main.BuildVersion=$(BUILD_VERSION)" cmd/main.go --version

# Better to run the integration tests
test-endpoints:
	./scripts/client.sh