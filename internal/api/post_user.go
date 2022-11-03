package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// postUser is an HTTP handler fot POST /users
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

	// todo: validate each field
	// logon_name not already present in DB & max length
	// full_name max length
	// email_address correct format & max length

	found, err := checkForUniqueLogonName(user.LogonName, env)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("validating logon_name uniqueness: %v", err))
		return
	} else if found {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("logon_name '%s' already present. Please choose another one", user.LogonName))
		return
	}

	user, err = env.UsersDB.addUser(user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("adding user to DB users table: %v", err))
		return
	}

	err = writeJSONHTTPResponse(w, user)
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
