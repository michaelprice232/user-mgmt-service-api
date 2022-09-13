package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

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
		log.WithError(err).Errorf("marshalling response into JSON: %v", resp)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		// Log & continue
		log.WithError(err).Errorf("writing HTTP response: %v", jsonResp)
	}

	log.WithFields(log.Fields{
		"statusCode": statusCode, "message": message, "url": r.URL.Path}).Error("writing non-2xx HTTP response")
}
