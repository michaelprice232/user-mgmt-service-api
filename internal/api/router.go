package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var dbConfig DBConfig

func RunAPIServer(c DBConfig) {
	serverAddr := "0.0.0.0:8080"
	dbConfig = c

	r := mux.NewRouter()
	r.HandleFunc("/users", listUsers).Methods("GET")

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
