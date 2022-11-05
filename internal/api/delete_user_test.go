package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockGetUsersModel is used to mock the Postgres DB calls
type mockDeleteUserModel struct{}

func (m *mockDeleteUserModel) queryUsers(_, _ int, _ string) (users []User, err error) {
	return
}

func (m *mockDeleteUserModel) addUser(_ User) (user User, err error) {
	return
}

func (m *mockDeleteUserModel) queryRecordCount(_, logonNameFilter string) (count int, err error) {
	switch logonNameFilter {
	case "testuser6":
		return 1, nil
	default:
		return 0, nil
	}
}

func (m *mockDeleteUserModel) deleteUser(_ string) error {
	return nil
}

func setupMockDeleteUserHTTPHandler(logonName string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%s", logonName), nil)
	env := &Env{UsersDB: &mockDeleteUserModel{}}

	// Need to create a router so that the URI parameters (logon_name) are picked up
	router := mux.NewRouter()
	router.HandleFunc("/users/{logon_name}", env.deleteUser)
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestDeleteUser(t *testing.T) {
	rec := setupMockDeleteUserHTTPHandler("testuser6")
	assert.Equal(t, 204, rec.Code)
}

func TestDeleteNotFoundUser(t *testing.T) {
	rec := setupMockDeleteUserHTTPHandler("testuser7")
	assert.Equal(t, 404, rec.Code)
}
