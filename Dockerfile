FROM golang:1.20 AS build

WORKDIR /usr/src/app

COPY ./certs/ca-cert.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /usr/local/bin/app ./cmd/main.go

FROM scratch

COPY --from=build /usr/local/bin/app /app

ENTRYPOINT ["/app"]