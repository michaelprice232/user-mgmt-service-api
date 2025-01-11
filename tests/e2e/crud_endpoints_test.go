//go:build e2e

package e2e

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestUsingAWS(t *testing.T) {
	//maxRetries := 5
	//timeBetweenRetries := 2 * time.Second

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../terraform",

		Vars: map[string]interface{}{
			"unique_identifier_prefix": strings.ToLower(random.UniqueId()),
		},
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "service_endpoint")

	url := fmt.Sprintf("%s/health", output)
	status, body := http_helper.HttpGet(t, url, &tls.Config{})
	assert.Equal(t, status, 200)
	t.Logf("Body result: %s", body)

	//endpoints.CheckEndpoints(t, baseURL, maxRetries, timeBetweenRetries)
}

// todo: seed database with ad-hoc task
