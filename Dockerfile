FROM golang:1.20 AS build

WORKDIR /usr/src/app

COPY ./certs/ca-cert.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /usr/local/bin/app ./cmd/main.go

FROM scratch

LABEL org.opencontainers.image.source=https://github.com/michaelprice232/user-mgmt-service-api

COPY --from=build /usr/local/bin/app /app

ENTRYPOINT ["/app"]