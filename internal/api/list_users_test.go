package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockGetUsersModel is used to mock the Postgres DB calls
type mockGetUsersModel struct{}

func (m *mockGetUsersModel) queryRecordCount(nameFilter, _ string) (int, error) {
	if nameFilter == "bob" {
		return 2, nil
	} else {
		return 5, nil
	}

}

func (m *mockGetUsersModel) queryUsers(offset, limit int, nameFilter string) ([]User, error) {
	var users []User

	if nameFilter == "bob" {
		users = []User{
			{UserID: 2, LogonName: "bob44", FullName: "bob", Email: "bob@email.com"},
			{UserID: 3, LogonName: "bobby8", FullName: "bobby", Email: "bobby@email.com"},
		}
	} else {
		if offset == 3 && limit == 3 {
			users = []User{
				{UserID: 4, LogonName: "jayne2234", FullName: "jayne", Email: "jayne@email.com"},
				{UserID: 5, LogonName: "mike1", FullName: "mike", Email: "mike@email.com"},
			}
		} else {
			users = []User{
				{UserID: 1, LogonName: "mark9", FullName: "mark", Email: "mark@email.com"},
				{UserID: 2, LogonName: "bob44", FullName: "bob", Email: "bob@email.com"},
				{UserID: 3, LogonName: "bobby8", FullName: "bobby", Email: "bobby@email.com"},
				{UserID: 4, LogonName: "jayne2234", FullName: "jayne", Email: "jayne@email.com"},
				{UserID: 5, LogonName: "mike1", FullName: "mike", Email: "mike@email.com"},
			}
		}

	}

	return users, nil
}

func (m *mockGetUsersModel) addUser(_ User) (user User, err error) {
	return
}
func (m *mockGetUsersModel) deleteUser(_ string) (err error) {
	return
}

func (m *mockGetUsersModel) updateUser(_ User) (user User, err error) { return }

// setupMockGetUsersHTTPHandler is helper function to remove duplication in setting up the HTTP test handlers in the unit tests
func setupMockGetUsersHTTPHandler(url string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	env := &Env{UsersDB: &mockGetUsersModel{}}
	http.HandlerFunc(env.listUsers).ServeHTTP(recorder, req)

	return recorder
}

// TestListUsersWithoutQueryParams tests listing users with no query params (using the paging defaults)
func TestListUsersWithoutQueryParams(t *testing.T) {
	var err error
	rec := setupMockGetUsersHTTPHandler("/users")

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp UsersResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	assert.Equal(t, 5, len(resp.Users))
	assert.Equal(t, "mark", resp.Users[0].FullName)
	assert.Equal(t, "mike@email.com", resp.Users[4].Email)
	assert.Equal(t, "mark9", resp.Users[0].LogonName)
	assert.Equal(t, 1, resp.Users[0].UserID)
}

// TestListUsersWithNameFilter tests listing users when a name_filter query parameter has been applied
func TestListUsersWithNameFilter(t *testing.T) {
	var err error
	rec := setupMockGetUsersHTTPHandler("/users?name_filter=bob")

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp UsersResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	assert.Equal(t, 2, len(resp.Users))
	assert.Equal(t, "bobby", resp.Users[1].FullName)
	assert.Equal(t, "bob@email.com", resp.Users[0].Email)
	assert.Equal(t, "bob44", resp.Users[0].LogonName)
	assert.Equal(t, 3, resp.Users[1].UserID)
}

// TestListUsersWithPagination tests using the pagination query parameters whilst listing users
func TestListUsersWithPagination(t *testing.T) {
	var err error
	rec := setupMockGetUsersHTTPHandler("/users?per_page=3&page=2")

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp UsersResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	assert.Equal(t, 2, len(resp.Users))
	assert.Equal(t, "jayne", resp.Users[0].FullName)
	assert.Equal(t, "mike@email.com", resp.Users[1].Email)
	assert.Equal(t, "mike1", resp.Users[1].LogonName)
	assert.Equal(t, 4, resp.Users[0].UserID)
}

// TestListUsersPerPageTooLarge tests for when a page size has been requested which exceeds the built in limit
func TestListUsersPerPageTooLarge(t *testing.T) {
	var err error
	rec := setupMockGetUsersHTTPHandler("/users?per_page=3000")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp JsonHTTPErrorResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Message, "per_page query string must be an integer between")
}

// TestListUsersPageNotFound tests for when a page has been requested which exceeds the number of user resources stored (based on total count)
func TestListUsersPageNotFound(t *testing.T) {
	var err error
	rec := setupMockGetUsersHTTPHandler("/users?page=1000")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	var resp JsonHTTPErrorResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Message, fmt.Sprintf("page %d not found", 1000))
}
