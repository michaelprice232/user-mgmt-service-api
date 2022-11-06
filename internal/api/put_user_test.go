package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockPutUserModel is used to mock the Postgres DB calls
type mockPutUserModel struct{}

func (m *mockPutUserModel) queryUsers(_, _ int, _ string) (users []User, err error) {
	return
}

func (m *mockPutUserModel) addUser(_ User) (user User, err error) {
	return
}

func (m *mockPutUserModel) deleteUser(_ string) (err error) {
	return
}

func (m *mockPutUserModel) queryRecordCount(_, logonNameFilter string) (count int, err error) {
	switch logonNameFilter {
	case "testuser8":
		return 1, nil
	case "testuser9":
		return 1, nil
	case "testuser10":
		return 1, nil
	default:
		return 0, nil
	}
}
func (m *mockPutUserModel) updateUser(user User) (User, error) {
	if user.LogonName == "testuser8" {
		// Both full_name and email being updated
		user.Email = "testuser8.updated@email.com"
		user.FullName = "Test User 8 Updated"
		user.UserID = 10
		user.LogonName = "testuser8"
	} else if user.LogonName == "testuser9" {
		// Just the email address updated
		user.Email = "testuser9.updated@email.com"
		user.FullName = "Test User 9"
		user.UserID = 11
		user.LogonName = "testuser9"
	} else if user.LogonName == "testuser10" {
		// Just the email address updated
		user.Email = "testuser10@email.com"
		user.FullName = "Test User 10"
		user.UserID = 12
		user.LogonName = "testuser10"
	}

	return user, nil
}

func setupMockPutUserHTTPHandler(logonName string, body bytes.Buffer) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", logonName), &body)
	env := &Env{UsersDB: &mockPutUserModel{}}

	// Need to create a router so that the URI parameters (logon_name) are picked up
	router := mux.NewRouter()
	router.HandleFunc("/users/{logon_name}", env.putUser)
	router.ServeHTTP(recorder, req)
	return recorder
}

func putRequestHelperSuccess(user User, logonName string, t *testing.T) (*httptest.ResponseRecorder, User) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(user)
	if err != nil {
		t.Fatal("unable to encode into buffer")
	}
	rec := setupMockPutUserHTTPHandler(logonName, buf)
	var resp User
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	return rec, resp
}

func putRequestHelperFailure(user User, logonName string, t *testing.T) (*httptest.ResponseRecorder, JsonHTTPErrorResponse) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(user)
	if err != nil {
		t.Fatal("unable to encode into buffer")
	}
	rec := setupMockPutUserHTTPHandler(logonName, buf)
	var resp JsonHTTPErrorResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	return rec, resp
}

// TestPutUser tests updating the email & full_name fields of an existing resource
func TestPutUser(t *testing.T) {
	// logon_name is extracted from the URI for PUT requests, so passing separate from the User object
	logonName := "testuser8"
	user := User{
		Email:    "testuser8.updated@email.com",
		FullName: "Test User 8 Updated",
	}
	rec, respUser := putRequestHelperSuccess(user, logonName, t)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, user.Email, respUser.Email)
	assert.Equal(t, user.FullName, respUser.FullName)
	assert.Equal(t, logonName, respUser.LogonName)
}

// TestPutUser tests updating just the email field of an existing resource
func TestPutUserJustEmail(t *testing.T) {
	// logon_name is extracted from the URI for PUT requests, so passing separate from the User object
	logonName := "testuser9"
	user := User{
		Email: "testuser9.updated@email.com",
	}
	rec, respUser := putRequestHelperSuccess(user, logonName, t)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, user.Email, respUser.Email)
	assert.Equal(t, "Test User 9", respUser.FullName)
	assert.Equal(t, logonName, respUser.LogonName)
}

// TestPutUser tests updating just the full_name field of an existing resource
func TestPutUserJustFullName(t *testing.T) {
	// logon_name is extracted from the URI for PUT requests, so passing separate from the User object
	logonName := "testuser10"
	user := User{
		FullName: "Test User 10",
	}
	rec, respUser := putRequestHelperSuccess(user, logonName, t)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "testuser10@email.com", respUser.Email)
	assert.Equal(t, user.FullName, respUser.FullName)
	assert.Equal(t, logonName, respUser.LogonName)
}

// TestPutUserBadUser tests trying to update a user which is not present in the DB
func TestPutUserBadUser(t *testing.T) {
	// logon_name is extracted from the URI for PUT requests, so passing separate from the User object
	logonName := "baduser"
	user := User{
		Email:    "bad.user@email.com",
		FullName: "Bad User 1",
	}
	rec, respUser := putRequestHelperFailure(user, logonName, t)
	assert.Equal(t, 404, rec.Code)
	assert.Equal(t, 404, respUser.Code)
	assert.Equal(t, fmt.Sprintf("'%s' does not exist. No action required", logonName), respUser.Message)
}

// TestPutUserBadUser tests trying to update a user with an email address format which is invalid
func TestPutUserBadEmailAddressFormat(t *testing.T) {
	// logon_name is extracted from the URI for PUT requests, so passing separate from the User object
	logonName := "testuser10"
	user := User{
		Email: "bad.email@",
	}
	rec, respUser := putRequestHelperFailure(user, logonName, t)
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, 400, respUser.Code)
	assert.Contains(t, respUser.Message, "validating email field format")
}

// TestPutUserBadUser tests trying to update a user with a full_name which exceeds the limits
func TestPutUserTooLongField(t *testing.T) {
	// logon_name is extracted from the URI for PUT requests, so passing separate from the User object
	logonName := "testuser10"
	tooLongFieldName := "qwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiop"
	user := User{
		FullName: tooLongFieldName,
	}
	rec, respUser := putRequestHelperFailure(user, logonName, t)
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, 400, respUser.Code)
	assert.Contains(t, respUser.Message, "full_name maximum length is 100")
}

// TestPutUserInvalidPayloadFields tests trying to update a user with request payload fields which are not supported (user_id & logon_name)
func TestPutUserInvalidPayloadFields(t *testing.T) {
	// logon_name is extracted from the URI for PUT requests, so passing separate from the User object
	logonName := "testuser8"
	user := User{
		Email:     "testuser8.updated@email.com",
		FullName:  "Test User 8 Updated",
		LogonName: "testuser8",
		UserID:    2,
	}
	rec, respUser := putRequestHelperFailure(user, logonName, t)
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, 400, respUser.Code)
	assert.Contains(t, respUser.Message, "logon_name and user_id are not supported request body fields for this operation")
}
