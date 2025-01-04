//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/michaelprice232/user-mgmt-service-api/tests/endpoints"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestUsingDockerCompose(t *testing.T) {
	maxRetries := 5
	timeBetweenRetries := 2 * time.Second
	hostPost := "8081"
	baseURL := fmt.Sprintf("http://localhost:%s", hostPost)

	buildOptions := docker.Options{
		WorkingDir:  "../..",
		ProjectName: fmt.Sprintf("user-mgmt-service-api-%s", random.UniqueId()),

		// Run on a different host port than the docker-compose instance instantiated locally from the Makefile, to avoid clashes
		EnvVars: map[string]string{"HOSTPORT": hostPost},
	}

	defer docker.RunDockerCompose(t, &buildOptions, "down")
	docker.RunDockerCompose(t, &buildOptions, "up", "-d")

	endpoints.CheckEndpoints(t, baseURL, maxRetries, timeBetweenRetries)
}
