package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockUserModel is used to mock the Postgres DB calls
type mockUserModel struct{}

func (m *mockUserModel) queryRecordCount(nameFilter string) (int, error) {
	if nameFilter == "bob" {
		return 2, nil
	} else {
		return 5, nil
	}

}

func (m *mockUserModel) queryUsers(offset, limit int, nameFilter string) ([]User, error) {
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

// setupMockGetUsersHTTPHandler is helper function to remove duplication in setting up the HTTP test handlers in the unit tests
func setupMockGetUsersHTTPHandler(url string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	env := &Env{UsersDB: &mockUserModel{}}
	http.HandlerFunc(env.listUsers).ServeHTTP(recorder, req)

	return recorder
}

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
