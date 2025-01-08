package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/hellofresh/health-go/v5"
	healthPg "github.com/hellofresh/health-go/v5/checks/postgres"
	log "github.com/sirupsen/logrus"
)

const (
	dbHealthCheckTimeOutSeconds = 5
	webServerPort               = 8080
	gracefulShutdownTime        = time.Second * 10
)

var EnvConfig *Env

// RunAPIServer starts an HTTP server after setting up any dependencies using *Env
func RunAPIServer() {
	serverAddr := fmt.Sprintf("0.0.0.0:%d", webServerPort)

	// Register a /health endpoint which polls the Postgres DB. Also display system & git build info
	h, err := health.New(health.WithSystemInfo(), health.WithComponent(health.Component{
		Name:    "user-mgmt-service-api",
		Version: EnvConfig.BuildVersion,
	}))
	if err != nil {
		log.WithError(err).Fatalf("unable to load health container")
	}
	err = h.Register(health.Config{
		Name:      "postgres-check",
		Timeout:   time.Second * dbHealthCheckTimeOutSeconds,
		SkipOnErr: false,
		Check: healthPg.New(healthPg.Config{
			DSN: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
				// Escape to avoid issues with special characters as we cannot quote in URLs
				url.QueryEscape(EnvConfig.DBCredentials.DBUsername), url.QueryEscape(EnvConfig.DBCredentials.DBPassword),
				EnvConfig.DBCredentials.HostName,
				EnvConfig.DBCredentials.Port, EnvConfig.DBCredentials.DBName, EnvConfig.DBCredentials.SSLMode),
		}),
	})

	r := mux.NewRouter()
	r.HandleFunc("/users", EnvConfig.listUsers).Methods("GET")
	r.HandleFunc("/users", EnvConfig.postUser).Methods("POST")
	r.HandleFunc("/users/{logon_name}", EnvConfig.deleteUser).Methods("DELETE")
	r.HandleFunc("/users/{logon_name}", EnvConfig.putUser).Methods("PUT")
	r.HandleFunc("/health", h.HandlerFunc)

	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Infof("Running webserver on: %s\n", serverAddr)

	go func() {
		if err = srv.ListenAndServe(); err != http.ErrServerClosed {
			log.WithError(err).Error("Problems shutting down HTTP server")
		} else {
			log.Infof("HTTP server shutdown ok")
		}
	}()

	// Graceful shutdown for SIGTERM & SIGINT OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	signalReceived := <-c
	log.Infof("OS signal received: %v", signalReceived)
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTime)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
