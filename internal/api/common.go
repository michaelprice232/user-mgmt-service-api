package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const ServiceName = "user-mgmt-service-api"

// jsonHTTPErrorResponseWriter writes non-2xx JSON responses back to the HTTP client as well as raising an error level log
func jsonHTTPErrorResponseWriter(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	var jsonResp []byte
	var err error
	resp := JSONHTTPErrorResponse{
		Code:    statusCode,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResp, err = json.Marshal(resp)
	if err != nil {
		// Log & continue
		log.WithError(err).Errorf("marshalling error response into JSON: %v", resp)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		// Log & continue
		log.WithError(err).Errorf("writing HTTP error response: %v", jsonResp)
	}

	// Only log at INFO if not a 5xx error
	if statusCode >= 500 {
		log.WithFields(log.Fields{
			"status_code": statusCode, "method": r.Method, "message": message, "url": getFullPathIncludingQueryParams(r.URL)}).Error("writing non-2xx HTTP response")
	} else {
		log.WithFields(log.Fields{
			"status_code": statusCode, "method": r.Method, "message": message, "url": getFullPathIncludingQueryParams(r.URL)}).Infof("writing non-2xx HTTP response")

	}
}

// writeJSONHTTPResponse writes data as a JSON payload back to the HTTP client
func writeJSONHTTPResponse(w http.ResponseWriter, responseCode int, payload interface{}) error {
	var err error
	var jsonResponse []byte

	jsonResponse, err = json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("marshalling JSON in preparation for HTTP response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	_, err = w.Write(jsonResponse)
	if err != nil {
		return fmt.Errorf("writing JSON formatted HTTP response")
	}
	return nil
}

// checkLogonNameExists returns true if logonName already exists in the DB
func checkLogonNameExists(logonName string, env *Env) (bool, error) {
	count, err := env.UsersDB.queryRecordCount("", logonName)
	if err != nil {
		return false, fmt.Errorf("checking to ensure that logon_name '%s' exists in database: %v", logonName, err)
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
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
		return fmt.Errorf("email maximum length is 100. Currently %d", len(user.Email))
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

// RequireStringEnvar returns a string envar and fatally exits if not set
func RequireStringEnvar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("envar '%s' not set. Exiting", key)
	}
	return value
}

// RequireIntEnvar returns an int envar and fatally exits if not yet set
func RequireIntEnvar(key string) int64 {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("envar '%s' not set. Exiting", key)
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("unable to convert envar '%s' into an integer. Exiting", key)
	}
	return i
}
