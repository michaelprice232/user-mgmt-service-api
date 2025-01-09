package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/michaelprice232/user-mgmt-service-api/internal/api"

	log "github.com/sirupsen/logrus"
)

var BuildVersion string // Set the git commit version from linker flags at build time

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
	log.Infof("Log level: %v", level)

	version := flag.Bool("version", false, "Returns the version of user-mgmt-service-api binary")
	flag.Parse()
	if *version {
		fmt.Printf("user-mgmt-service-api version: %s (OS: %s) (Arch: %s)\n", BuildVersion, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	api.EnvConfig, err = api.OpenDBConnection()
	if err != nil {
		log.WithError(err).Fatal("opening DB connection")
	}

	api.EnvConfig.BuildVersion = BuildVersion
}

func main() {
	api.RunAPIServer()
}
