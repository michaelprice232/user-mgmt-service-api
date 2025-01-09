# Cross compilation multi-architecture build
# https://docs.docker.com/build/building/multi-platform/#cross-compiling-a-go-application
FROM --platform=$BUILDPLATFORM golang:1.23 AS build

# These are made available when using the --platform docker build parameter, along with BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Add default for when the version is not passed to the Docker build
ARG BUILD_VERSION=dev

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build -ldflags="-X main.BuildVersion=${BUILD_VERSION}" -o /usr/local/bin/app ./cmd/main.go

FROM scratch

ARG BUILD_VERSION

LABEL org.opencontainers.image.title="User management REST API"
LABEL org.opencontainers.image.description="A simple CRUD API"
LABEL org.opencontainers.image.source=https://github.com/michaelprice232/user-mgmt-service-api
LABEL org.opencontainers.image.revision=${BUILD_VERSION}

COPY --from=build /usr/local/bin/app /app

ENTRYPOINT ["/app"]