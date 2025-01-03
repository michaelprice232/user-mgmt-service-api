#FROM --platform=$BUILDPLATFORM golang:1.23 AS build
FROM golang:1.23 AS build

## These are made available when using the --platform docker build parameter
#ARG TARGETOS
#ARG TARGETARCH

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /usr/local/bin/app ./cmd/main.go

FROM scratch

LABEL org.opencontainers.image.source=https://github.com/michaelprice232/user-mgmt-service-api

COPY --from=build /usr/local/bin/app /app

ENTRYPOINT ["/app"]