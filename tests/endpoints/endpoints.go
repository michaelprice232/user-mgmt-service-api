//go:build integration

package endpoints

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/http-helper"
)

// CheckEndpoints performs HTTP requests against the CRUD endpoints. These can be re-used between the integration and E2E tests.
func CheckEndpoints(t *testing.T, baseURL string, maxRetries int, timeBetweenRetries time.Duration) {
	t.Run("GET_/users_successfully", func(t *testing.T) {
		url := fmt.Sprintf("%s/users", strings.TrimSuffix(baseURL, "/"))
		expectedResponseBody := `{"Users":[{"user_id":1,"logon_name":"mike1","full_name":"mike","email":"mike@email.com"},{"user_id":2,"logon_name":"bob44","full_name":"bob","email":"bob@email.com"},{"user_id":3,"logon_name":"sarah485","full_name":"sarah","email":"sarah@email.com"},{"user_id":4,"logon_name":"eric2","full_name":"eric","email":"eric@email.com"}],"total_pages":3,"current_page":1,"more_pages":true}`
		http_helper.HttpGetWithRetry(t, url, &tls.Config{}, http.StatusOK, expectedResponseBody, maxRetries, timeBetweenRetries)
	})
}
