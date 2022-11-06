package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"user-mgmt-service-api/internal/api"
)

var BuildVersion string

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

	// Load DB credentials and start SQL DB connection pool
	api.EnvConfig = &api.Env{DBCredentials: api.DBCredentials{
		HostName:   RequireStringEnvar("database_host_name"),
		Port:       uint(RequireIntEnvar("database_port")),
		DBName:     RequireStringEnvar("database_name"),
		DBUsername: RequireStringEnvar("database_username"),
		DBPassword: RequireStringEnvar("database_password"),
		SSLMode:    RequireStringEnvar("database_ssl_mode")}}

	sqlConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		api.EnvConfig.DBCredentials.HostName, api.EnvConfig.DBCredentials.Port, api.EnvConfig.DBCredentials.DBUsername,
		api.EnvConfig.DBCredentials.DBPassword, api.EnvConfig.DBCredentials.DBName, api.EnvConfig.DBCredentials.SSLMode)

	db, err := sql.Open("postgres", sqlConnection)
	if err != nil {
		log.WithError(err).Fatal("opening DB connection pool")
	}

	// Set the git commit version from linker flags at build time
	api.EnvConfig.BuildVersion = BuildVersion
	api.EnvConfig.UsersDB = &api.UserModel{DB: db}
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

func RequireIntEnvar(key string) int64 {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("envar '%s' not set. Exiting", key)
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("unable to convert envar '%s' into an integer. Exiting", key)
	}
	return i
}
