BUILD_VERSION := $(shell git rev-parse --short HEAD)

run:
	HOSTPORT=8080 docker-compose up -d --build

down:
	HOSTPORT=8080 docker-compose down

unit-tests:
	go test ./...

# Do not use cached test results
int-tests:
	go test -tags=integration -count=1 ./tests/integration

version:
	go run -ldflags="-X 'main.BuildVersion=$(BUILD_VERSION)'" cmd/main.go --version

# Better to run the integration tests
test-endpoints:
	./scripts/client.sh