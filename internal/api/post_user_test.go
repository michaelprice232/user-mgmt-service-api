package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockPostUserModel is used to mock the Postgres DB calls
type mockPostUserModel struct{}

func (m *mockPostUserModel) queryUsers(_, _ int, _ string) (users []User, err error) {
	return
}

func (m *mockPostUserModel) queryRecordCount(_, logonNameFilter string) (count int, err error) {
	switch logonNameFilter {
	case "testuser2":
		return 1, nil
	default:
		return 0, nil
	}
}

func (m *mockPostUserModel) addUser(user User) (User, error) {
	switch user.LogonName {
	case "testuser1":
		user.UserID = 11
	}
	return user, nil
}

func (m *mockPostUserModel) deleteUser(_ string) (err error) {
	return
}

func (m *mockPostUserModel) updateUser(_ User) (user User, err error) { return }

func setupMockPostUserHTTPHandler(body bytes.Buffer) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", &body)
	env := &Env{UsersDB: &mockPostUserModel{}}
	http.HandlerFunc(env.postUser).ServeHTTP(recorder, req)

	return recorder
}

func postRequestHelperSuccess(user User, t *testing.T) (*httptest.ResponseRecorder, User) {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(user)
	rec := setupMockPostUserHTTPHandler(buf)
	var resp User
	err := json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	return rec, resp
}

func postRequestHelperFailure(user User, t *testing.T) (*httptest.ResponseRecorder, JsonHTTPErrorResponse) {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(user)
	rec := setupMockPostUserHTTPHandler(buf)

	var resp JsonHTTPErrorResponse
	err := json.Unmarshal([]byte(rec.Body.String()), &resp)
	if err != nil {
		t.Fatal("unable to unmarshal JSON response")
	}
	return rec, resp
}

// TestAddUser tests adding a new user
func TestAddUser(t *testing.T) {
	user := User{
		LogonName: "testuser1",
		FullName:  "Test User 1",
		Email:     "test@email.com",
	}
	rec, resp := postRequestHelperSuccess(user, t)
	assert.Equal(t, 201, rec.Code)
	assert.Equal(t, 11, resp.UserID)
	assert.Equal(t, user.LogonName, resp.LogonName)
	assert.Equal(t, user.FullName, resp.FullName)
	assert.Equal(t, user.Email, resp.Email)
}

// TestAddUserLogonAlreadyTaken tests attempting to add a new user when the logon_name is already taken in the DB
func TestAddUserLogonAlreadyTaken(t *testing.T) {
	user := User{
		LogonName: "testuser2",
		FullName:  "Test User 2",
		Email:     "test2@email.com",
	}
	rec, resp := postRequestHelperFailure(user, t)

	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, fmt.Sprintf("logon_name '%s' already taken. Please choose another one", user.LogonName), resp.Message)
}

// TestAddUserFieldLengthTooLong tests that the validation around field lengths is working as expected
func TestAddUserFieldLengthTooLong(t *testing.T) {
	longFieldName := "qwertyuiopqwertyuiopqwertyuiop"
	user := User{
		LogonName: longFieldName,
		FullName:  "Test User 3",
		Email:     "test3@email.com",
	}
	rec, resp := postRequestHelperFailure(user, t)

	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Message, fmt.Sprintf("validating request payload field lengths: logon_name maximum lengh is 20. Currently %d", len(longFieldName)))
}

// TestAddUserInvalidEmailFieldFormat tests that the validation around email field format is working as expected
func TestAddUserInvalidEmailFieldFormat(t *testing.T) {
	badEmailFormat := "test3@"
	user := User{
		LogonName: "testuser4",
		FullName:  "Test User 4",
		Email:     badEmailFormat,
	}
	rec, resp := postRequestHelperFailure(user, t)

	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Message, fmt.Sprintf("'%s' not a valid email address field:", badEmailFormat))
}

// TestAddUserPassedTheUserLogonField tests that an unsupported field - user_id - is handled correctly
func TestAddUserPassedTheUserLogonField(t *testing.T) {
	user := User{
		UserID:    3, // not supported in the request payload
		LogonName: "testuser4",
		FullName:  "Test User 4",
		Email:     "test4@email.com",
	}
	rec, resp := postRequestHelperFailure(user, t)

	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "passing a user_id in the request payload is not supported", resp.Message)
}
