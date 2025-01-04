package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/assert"
)

func TestUsingDockerCompose(t *testing.T) {
	buildOptions := docker.Options{
		WorkingDir: "../..",
		ProjectName: "user-mgmt-service-api",
	}
	output := docker.RunDockerCompose(t, &buildOptions, "up")

	assert.NotEmpty(t, output)

	t.Logf("output = %s", output)

	//tag := "gruntwork/docker-hello-world-example"
	//buildOptions := &docker.BuildOptions{
	//	Tags: []string{tag},
	//}
	//
	//docker.Build(t, "../../", buildOptions)
	//
	//opts := &docker.RunOptions{Command: []string{"cat", "/test.txt"}}
	//output := docker.Run(t, tag, opts)
	//assert.Equal(t, "Hello, World!", output)
}
