package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			{Name: "bob", Email: "bob@email.com"},
			{Name: "bobby", Email: "bobby@email.com"},
		}
	} else {
		if offset == 3 && limit == 3 {
			users = []User{
				{Name: "jayne", Email: "jayne@email.com"},
				{Name: "mike", Email: "mike@email.com"},
			}
		} else {
			users = []User{
				{Name: "mark", Email: "mark@email.com"},
				{Name: "bob", Email: "bob@email.com"},
				{Name: "bobby", Email: "bobby@email.com"},
				{Name: "jayne", Email: "jayne@email.com"},
				{Name: "mike", Email: "mike@email.com"},
			}
		}

	}

	return users, nil
}

func TestListUsersWithoutQueryParams(t *testing.T) {
	var err error
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	env := &Env{UsersDB: &mockUserModel{}}
	http.HandlerFunc(env.listUsers).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp UsersResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	assert.Equal(t, 5, len(resp.Users))
	assert.Equal(t, "mark", resp.Users[0].Name)
	assert.Equal(t, "mike@email.com", resp.Users[4].Email)
}

func TestListUsersWithNameFilter(t *testing.T) {
	var err error
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users?name_filter=bob", nil)
	env := &Env{UsersDB: &mockUserModel{}}
	http.HandlerFunc(env.listUsers).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp UsersResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	assert.Equal(t, 2, len(resp.Users))
	assert.Equal(t, "bobby", resp.Users[1].Name)
	assert.Equal(t, "bob@email.com", resp.Users[0].Email)
}

func TestListUsersWithPagination(t *testing.T) {
	var err error
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users?per_page=3&page=2", nil)
	env := &Env{UsersDB: &mockUserModel{}}
	http.HandlerFunc(env.listUsers).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp UsersResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	assert.Equal(t, 2, len(resp.Users))
	assert.Equal(t, "jayne", resp.Users[0].Name)
	assert.Equal(t, "mike@email.com", resp.Users[1].Email)
}
