package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// deleteUser is an HTTP handler for DELETE /users/<user>
func (env *Env) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetLogonName := vars["logon_name"]
	log.Infof("Received delete request for logon_name '%s'", targetLogonName)

	exists, err := checkLogonNameExists(targetLogonName, env)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("checking logon_name in database: %v", err))
		return
	}

	if exists {
		log.Infof("'%s' exists. Deleting user from the DB", targetLogonName)
		err = env.UsersDB.deleteUser(targetLogonName)
		if err != nil {
			jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("deleting user from DB: %v", err))
			return
		}
		w.WriteHeader(204)
		return

	} else {
		jsonHTTPErrorResponseWriter(w, r, 404, fmt.Sprintf("'%s' does not exist. No deletion required", targetLogonName))
		return
	}
}

// checkLogonNameExists returns true if logonName already exists in the DB
func checkLogonNameExists(logonName string, env *Env) (bool, error) {
	count, err := env.UsersDB.queryRecordCount("", logonName)
	if err != nil {
		return false, fmt.Errorf("checking to ensure that logon_name '%s' exists in database: %v", logonName, err)
	}
	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
