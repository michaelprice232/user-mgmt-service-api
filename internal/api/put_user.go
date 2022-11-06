package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// putUser is a HTTP handler for PUT /users/<logon_name>
func (env *Env) putUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetLogonName := vars["logon_name"]
	log.Infof("Received PUT request for logon_name '%s'", targetLogonName)

	exists, err := checkLogonNameExists(targetLogonName, env)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("checking logon_name against database: %v", err))
		return
	}

	if exists {
		log.Infof("'%s' exists. Updating user in the DB", targetLogonName)

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
		log.Debugf("Unmarshaled payload: %#v", user)

		// Set LogonName based on URI so that the DB query can locate the user's record
		user.LogonName = targetLogonName

		err = validatePutRequestPayload(user, w, r)
		if err != nil {
			return
		}

		userResp, err := env.UsersDB.updateUser(user)
		if err != nil {
			jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("updating record for user '%s' in DB: %v", targetLogonName, err))
			return
		}

		// Return the updated record back to the client
		err = writeJSONHTTPResponse(w, 200, userResp)
		if err != nil {
			jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("writing HTTP response: %v", err))
			return
		}

		log.WithFields(log.Fields{
			"url":         getFullPathIncludingQueryParams(r.URL),
			"status_code": 200,
			"method":      r.Method,
			"logon_name":  user.LogonName,
		}).Infof("serving page")

	} else {
		jsonHTTPErrorResponseWriter(w, r, 404, fmt.Sprintf("'%s' does not exist. No action required", targetLogonName))
		return
	}
}

// validatePutRequestPayload validates the request payload of the PUT /users/<logon_name> operation
func validatePutRequestPayload(user User, w http.ResponseWriter, r *http.Request) error {
	err := validateFieldLengths(user)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("validating PUT request payload field lengths: %v", err))
		return err
	}

	if user.UserID != 0 {
		resp := "logon_name and user_id are not supported request body fields for this operation"
		jsonHTTPErrorResponseWriter(w, r, 400, resp)
		return fmt.Errorf(resp)
	}

	// optional field to pass
	if user.Email != "" {
		err = validateEmailField(user.Email)
		if err != nil {
			jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("validating email field format: %v", err))
			return err
		}
	}
	return nil
}
