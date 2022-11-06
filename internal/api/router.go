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

	r := mux.NewRouter()
	r.HandleFunc("/users", EnvConfig.listUsers).Methods("GET")
	r.HandleFunc("/users", EnvConfig.postUser).Methods("POST")
	r.HandleFunc("/users/{logon_name}", EnvConfig.deleteUser).Methods("DELETE")
	r.HandleFunc("/users/{logon_name}", EnvConfig.putUser).Methods("PUT")

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
