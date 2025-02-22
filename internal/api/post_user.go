package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// postUser is an HTTP handler for POST /users
func (env *Env) postUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("reading http request body: %v", err))
		return
	}
	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("unmarshalling http request body: %v", err))
		return
	}
	log.Infof("Unmarshaled payload: %#v", user)

	err = validateRequestPayload(user, env, w, r)
	if err != nil {
		return
	}

	user, err = env.UsersDB.addUser(user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("adding user to DB users table: %v", err))
		return
	}

	err = writeJSONHTTPResponse(w, 201, user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("writing HTTP response: %v", err))
		return
	}

	log.WithFields(log.Fields{
		"url":         getFullPathIncludingQueryParams(r.URL),
		"status_code": 201,
		"method":      r.Method,
		"logon_name":  user.LogonName,
	}).Infof("serving page")
}

// validateRequestPayload validates the request payload of the POST /users/<logon_name> operation
func validateRequestPayload(user User, env *Env, w http.ResponseWriter, r *http.Request) error {
	err := validateFieldLengths(user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("validating request payload field lengths: %v", err))
		return err
	}

	if user.UserID != 0 {
		resp1 := "passing a user_id in the request payload is not supported"
		jsonHTTPErrorResponseWriter(w, r, 400, resp1)
		return fmt.Errorf("%s", resp1)
	}

	err = validateEmailField(user.Email)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("validating email field format: %v", err))
		return err
	}

	found, err := checkForUniqueLogonName(user.LogonName, env)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("validating logon_name uniqueness: %v", err))
		return err
	} else if found {
		resp2 := fmt.Sprintf("logon_name '%s' already taken. Please choose another one", user.LogonName)
		jsonHTTPErrorResponseWriter(w, r, 400, resp2)
		return fmt.Errorf("%s", resp2)
	}

	return nil
}

// checkForUniqueLogonName queries the database to see if the logon_name is already present in the users table
func checkForUniqueLogonName(logonName string, env *Env) (bool, error) {
	count, err := env.UsersDB.queryRecordCount("", logonName)
	if err != nil {
		return false, fmt.Errorf("checking database for unique logon_name '%s': %v", logonName, err)
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}
