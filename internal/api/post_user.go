package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
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

	err = validateRequestPayload(user, env, w, r)
	if err != nil {
		return
	}

	user, err = env.UsersDB.addUser(user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("adding user to DB users table: %v", err))
		return
	}

	err = writeJSONHTTPResponse(w, user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("writing HTTP response: %v", err))
		return
	}
}

// validateRequestPayload validates the JSON request payload that has been sent by the client
func validateRequestPayload(user User, env *Env, w http.ResponseWriter, r *http.Request) error {
	err := validateFieldLengths(user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("validating request payload field lengths: %v", err))
		return err
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
		resp := fmt.Sprintf("logon_name '%s' already taken. Please choose another one", user.LogonName)
		jsonHTTPErrorResponseWriter(w, r, 400, resp)
		return fmt.Errorf(resp)
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
	} else {
		return false, nil
	}
}

// validateFieldLengths validates that each of the User fields do not exceed the database table limits
func validateFieldLengths(user User) error {
	if len(user.LogonName) > 20 {
		return fmt.Errorf("logon_name maximum lengh is 20. Currently %d", len(user.LogonName))
	}
	if len(user.FullName) > 100 {
		return fmt.Errorf("full_name maximum length is 100. Currently %d", len(user.FullName))
	}
	if len(user.Email) > 100 {
		return fmt.Errorf("email maxium length is 100. Currently %d", len(user.Email))
	}

	return nil
}

// validateEmailField validates that the parameter is in a valid email address format
func validateEmailField(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("'%s' not a valid email address field: %v", email, err)
	}
	return nil
}
