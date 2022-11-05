package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// jsonHTTPErrorResponseWriter writes non-2xx JSON responses back to the HTTP client as well as raising an error level log
func jsonHTTPErrorResponseWriter(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	var jsonResp []byte
	var err error
	resp := JsonHTTPErrorResponse{
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

	log.WithFields(log.Fields{
		"statusCode": statusCode, "message": message, "url": getFullPathIncludingQueryParams(r.URL)}).Error("writing non-2xx HTTP response")
}

// writeJSONHTTPResponse writes data as a JSON payload back to the HTTP client
func writeJSONHTTPResponse(w http.ResponseWriter, payload interface{}) error {
	var err error
	var jsonResponse []byte

	jsonResponse, err = json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("marshalling JSON in preparation for HTTP response")
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		return fmt.Errorf("writing JSON formatted HTTP response")
	}
	return nil
}
