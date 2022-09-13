package main

import (
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
}

func main() {
	api.RunAPIServer()
}
