package main

import (
	"database/sql"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"user-mgmt-service-api/internal/api"
)

func init() {
	var level log.Level
	var err error

	log.SetOutput(os.Stderr)

	if _, isLocal := os.LookupEnv("RUNNING_LOCALLY"); isLocal {
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	llEnv, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		// default to warnings
		log.SetLevel(log.WarnLevel)
	} else {
		level, err = log.ParseLevel(llEnv)
		if err != nil {
			log.WithError(err).Fatal("parsing logrus log level from LOG_LEVEL envar")
		}
		log.SetLevel(level)
	}
	log.Infof("Log level: %v\n", level)

	dbName := RequireStringEnvar("database_name")
	dbUsername := RequireStringEnvar("database_username")
	dbPassword := RequireStringEnvar("database_password")
	dbSslMode := RequireStringEnvar("database_ssl_mode")

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUsername, dbPassword, dbName, dbSslMode))
	if err != nil {
		log.WithError(err).Fatal("opening DB connection pool")
	}

	api.EnvConfig = &api.Env{UsersDB: &api.UserModel{DB: db}}

}

func main() {
	api.RunAPIServer()
}

func RequireStringEnvar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("envar '%s' not set. Exiting", key)
	}
	return value
}
