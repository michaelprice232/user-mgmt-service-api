//go:build integration

package integration

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestUsingDockerCompose(t *testing.T) {
	buildOptions := docker.Options{
		WorkingDir:  "../..",
		ProjectName: fmt.Sprintf("user-mgmt-service-api-%s", random.UniqueId()),

		// Run on a different host port than the docker-compose instance instantiated locally from the Makefile, to avoid clashes
		EnvVars: map[string]string{"HOSTPORT": "8081"},
	}

	defer docker.RunDockerCompose(t, &buildOptions, "down")
	docker.RunDockerCompose(t, &buildOptions, "up", "-d")

	maxRetries := 5
	timeBetweenRetries := 2 * time.Second
	url := "http://localhost:8081/users"

	responseBody := `{"Users":[{"user_id":1,"logon_name":"mike1","full_name":"mike","email":"mike@email.com"},{"user_id":2,"logon_name":"bob44","full_name":"bob","email":"bob@email.com"},{"user_id":3,"logon_name":"sarah485","full_name":"sarah","email":"sarah@email.com"},{"user_id":4,"logon_name":"eric2","full_name":"eric","email":"eric@email.com"}],"total_pages":3,"current_page":1,"more_pages":true`
	http_helper.HttpGetWithRetry(t, url, &tls.Config{}, http.StatusOK, responseBody, maxRetries, timeBetweenRetries)
}
