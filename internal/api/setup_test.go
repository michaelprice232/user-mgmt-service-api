package api

import (
	log "github.com/sirupsen/logrus"

	"testing"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.ErrorLevel)
}
