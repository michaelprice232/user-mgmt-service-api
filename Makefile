BUILD_VERSION := $(shell git rev-parse --short HEAD)

run:
	HOSTPORT=8080 BUILD_VERSION=$(BUILD_VERSION) docker-compose up -d --build

down:
	HOSTPORT=8080 docker-compose down --volumes

lint:
	golangci-lint run

unit-tests:
	go test -v ./...

int-tests:
	go test -tags=integration -count=1 -v ./tests/integration

e2e-tests:
	export DOCKER_APP_IMAGE="633681147894.dkr.ecr.eu-west-2.amazonaws.com/user-mgmt-service-api:73a46c8ce278e6d205915f66b3e80b9ff61dc090"; \
	go test -tags=e2e -count=1 -v -timeout 60m ./tests/e2e

version:
	go run -ldflags="-X main.BuildVersion=$(BUILD_VERSION)" cmd/main.go --version

# Better to run the integration tests
test-endpoints:
	./scripts/client.sh