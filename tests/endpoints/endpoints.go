//go:build integration

package endpoints

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/michaelprice232/user-mgmt-service-api/internal/api"

	"github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/stretchr/testify/assert"
)

func unmarshalJSONUsersResponse(t *testing.T, input string) api.UsersResponse {
	resp := api.UsersResponse{}
	err := json.Unmarshal([]byte(input), &resp)
	assert.NoError(t, err)
	return resp
}

func unmarshalJSONErrorResponse(t *testing.T, input string) api.JsonHTTPErrorResponse {
	resp := api.JsonHTTPErrorResponse{}
	err := json.Unmarshal([]byte(input), &resp)
	assert.NoError(t, err)
	return resp
}

func unmarshalJSONUserResponse(t *testing.T, input string) api.User {
	resp := api.User{}
	err := json.Unmarshal([]byte(input), &resp)
	assert.NoError(t, err)
	return resp
}

// CheckEndpoints performs HTTP requests against the CRUD endpoints. These can be re-used between the integration and E2E tests.
func CheckEndpoints(t *testing.T, baseURL string, maxRetries int, timeBetweenRetries time.Duration) {
	baseURLFormatted := strings.TrimSuffix(baseURL, "/")

	// Successful GET requests
	t.Run("GET /users no params", func(t *testing.T) {
		url := fmt.Sprintf("%s/users", baseURLFormatted)
		http_helper.HttpGetWithRetryWithCustomValidation(t, url, &tls.Config{}, maxRetries, timeBetweenRetries, func(statusCode int, responseBody string) bool {
			if statusCode != http.StatusOK {
				return false
			}
			resp := unmarshalJSONUsersResponse(t, responseBody)
			assert.Equal(t, 4, len(resp.Users), "Expected 4 users to be returned")
			assert.Equal(t, "bob@email.com", resp.Users[1].Email, "Expected 2nd returned user to be Bob")
			assert.Equal(t, 3, resp.TotalPages, "Expected 3 total pages to be available")
			return true
		})
	})

	t.Run("GET /users with pagination", func(t *testing.T) {
		url := fmt.Sprintf("%s/users?per_page=5&page=2", baseURLFormatted)
		http_helper.HttpGetWithRetryWithCustomValidation(t, url, &tls.Config{}, maxRetries, timeBetweenRetries, func(statusCode int, responseBody string) bool {
			if statusCode != http.StatusOK {
				return false
			}
			resp := unmarshalJSONUsersResponse(t, responseBody)
			assert.Equal(t, 5, len(resp.Users), "Expected 5 users to be returned")
			assert.Equal(t, "jayne@email.com", resp.Users[len(resp.Users)-1].Email, "Expected the last user to be Jayne")
			assert.Equal(t, 2, resp.CurrentPage, "Expected the current returned page to be 2")
			return true
		})
	})

	t.Run("GET /users with filtering", func(t *testing.T) {
		url := fmt.Sprintf("%s/users?name_filter=bob", baseURLFormatted)
		http_helper.HttpGetWithRetryWithCustomValidation(t, url, &tls.Config{}, maxRetries, timeBetweenRetries, func(statusCode int, responseBody string) bool {
			if statusCode != http.StatusOK {
				return false
			}
			resp := unmarshalJSONUsersResponse(t, responseBody)
			assert.Equal(t, 2, len(resp.Users), "Expected 2 users to be returned")
			assert.Equal(t, "bob44", resp.Users[0].LogonName, "Expected the first user to be Bob")
			assert.Equal(t, 1, resp.TotalPages, "Expected the total pages to be 1")
			return true
		})
	})

	t.Run("GET /users with filtering and pagination", func(t *testing.T) {
		url := fmt.Sprintf("%s/users?name_filter=bob&per_page=1&page=1", baseURLFormatted)
		http_helper.HttpGetWithRetryWithCustomValidation(t, url, &tls.Config{}, maxRetries, timeBetweenRetries, func(statusCode int, responseBody string) bool {
			if statusCode != http.StatusOK {
				return false
			}
			resp := unmarshalJSONUsersResponse(t, responseBody)
			assert.Equal(t, 1, len(resp.Users), "Expected 1 users to be returned")
			assert.Equal(t, 2, resp.Users[0].UserID, "Expected the first user to be Bob")
			assert.Equal(t, 2, resp.TotalPages, "Expected the total pages to be 2")
			assert.True(t, resp.MorePages, "Expected more pages to be available")
			return true
		})
	})

	// Exception based GET requests
	t.Run("GET /users and per_page too large", func(t *testing.T) {
		url := fmt.Sprintf("%s/users?per_page=2000", baseURLFormatted)
		http_helper.HttpGetWithRetryWithCustomValidation(t, url, &tls.Config{}, maxRetries, timeBetweenRetries, func(statusCode int, responseBody string) bool {
			if statusCode != http.StatusBadRequest {
				return false
			}
			resp := unmarshalJSONErrorResponse(t, responseBody)
			assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected bad response code to be in the response body")
			assert.Contains(t, resp.Message, "per_page query string must be an integer between", "Expected details to be in the error message")
			return true
		})
	})

	t.Run("GET /users and page not found", func(t *testing.T) {
		url := fmt.Sprintf("%s/users?page=1000", baseURLFormatted)
		http_helper.HttpGetWithRetryWithCustomValidation(t, url, &tls.Config{}, maxRetries, timeBetweenRetries, func(statusCode int, responseBody string) bool {
			if statusCode != http.StatusNotFound {
				return false
			}
			resp := unmarshalJSONErrorResponse(t, responseBody)
			assert.Equal(t, http.StatusNotFound, resp.Code, "Expected bad response code to be in the response body")
			assert.Contains(t, resp.Message, "not found", "Expected details to be in the error message")
			return true
		})
	})

	// Successful POST requests
	t.Run("POST /users", func(t *testing.T) {
		url := fmt.Sprintf("%s/users", baseURLFormatted)
		input := api.User{
			LogonName: "testuser1",
			FullName:  "Test User 1",
			Email:     "test1@email.com",
		}
		body, err := json.Marshal(input)
		assert.NoError(t, err)
		bodyInput := bytes.NewReader(body)
		http_helper.HTTPDoWithCustomValidation(t, "POST", url, bodyInput, map[string]string{"Content-Type": "application/json"}, func(statusCode int, responseBody string) bool {
			if statusCode != http.StatusCreated {
				return false
			}
			resp := unmarshalJSONUserResponse(t, responseBody)
			assert.Equal(t, "testuser1", resp.LogonName, "Expected the returned user to have a logon name of testuser1")
			assert.Equal(t, "Test User 1", resp.FullName, "Expected the returned user to have a full name of Test User 1")
			return true
		}, &tls.Config{})

		// Check we can retrieve the newly inserted user
		// We have to url encode the whitespace in the query string
		url = fmt.Sprintf("%s/users?name_filter=Test+User+1", baseURLFormatted)
		http_helper.HttpGetWithRetryWithCustomValidation(t, url, &tls.Config{}, maxRetries, timeBetweenRetries, func(statusCode int, responseBody string) bool {
			if statusCode != http.StatusOK {
				return false
			}
			resp := unmarshalJSONUsersResponse(t, responseBody)
			assert.Equal(t, 1, len(resp.Users), "Expected 1 user to be returned")
			assert.Equal(t, "test1@email.com", resp.Users[0].Email, "Expected the correct email address to be returned")
			return true
		})
	})

}
