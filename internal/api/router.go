package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var EnvConfig *Env

// RunAPIServer starts an HTTP server after setting up any dependencies using *Env
func RunAPIServer() {
	serverAddr := "0.0.0.0:8080"

	// todo: add a /health endpoint
	r := mux.NewRouter()
	r.HandleFunc("/users", EnvConfig.listUsers).Methods("GET")
	r.HandleFunc("/users", EnvConfig.postUser).Methods("POST")

	// todo: enable graceful shutdowns from the appropriate OS signals
	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Infof("Running webserver on: %s\n", serverAddr)
	log.Fatal(srv.ListenAndServe())
}
