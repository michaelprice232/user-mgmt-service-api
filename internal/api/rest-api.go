package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Name  string
	Email string
}

type Users []User

type JsonHTTPResponse struct {
	Code    int
	Message string
}

const (
	defaultPageSize = 2
	maxPageSize     = 5
)

var usersDB = Users{
	{Name: "jayne 0", Email: "michaelprice232@outlook.com"},
	{Name: "jayne 1", Email: "jaynefreer30031984@yahoo.co.uk"},
	{Name: "jayne 2", Email: "jaynefreer30031984@yahoo.co.uk"},
	{Name: "jayne 3", Email: "jaynefreer30031984@yahoo.co.uk"},
	{Name: "bob 1", Email: "jaynefreer30031984@yahoo.co.uk"},
	{Name: "jayne 5", Email: "jaynefreer30031984@yahoo.co.uk"},
	{Name: "jayne 6", Email: "jaynefreer30031984@yahoo.co.uk"},
	{Name: "bob 2", Email: "jaynefreer30031984@yahoo.co.uk"},
	{Name: "jayne 8", Email: "jaynefreer30031984@yahoo.co.uk"},
	{Name: "jayne 9", Email: "jaynefreer30031984@yahoo.co.uk"},
}

func RunAPIServer() {
	serverAddr := "0.0.0.0:8080"
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/users", listUsers).Methods("GET")

	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Infof("Running webserver on: %s\n", serverAddr)
	log.Fatal(srv.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Root Page!\n"))
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	var perPage = defaultPageSize
	var perPage64, page64 int64
	var page = 1
	var nameFilter string

	queryStrings := r.URL.Query()
	if perPageEnv := queryStrings.Get("per_page"); perPageEnv != "" {
		perPage64, err = strconv.ParseInt(perPageEnv, 10, 64)
		if err != nil {
			perPage64 = defaultPageSize
		}
		perPage = int(perPage64)

		if perPage > maxPageSize {
			perPage = maxPageSize
		}
	}

	// todo: read from DB instead of this duplication
	var dbResults Users
	if nameFilter = queryStrings.Get("name_filter"); nameFilter != "" {
		for _, user := range usersDB {
			if strings.Contains(user.Name, nameFilter) {
				dbResults = append(dbResults, user)
			}
		}
	} else {
		dbResults = usersDB
	}

	totalItems := len(dbResults)
	numberOfPages := totalItems / perPage
	if totalItems%perPage != 0 {
		// Add a non-full page
		numberOfPages++
	}

	if pageEnv := queryStrings.Get("page"); pageEnv != "" {
		page64, err = strconv.ParseInt(pageEnv, 10, 64)
		if err != nil || page64 < 0 || page64 > int64(numberOfPages) {
			jsonHTTPResponseWriter(w, 404, "Page not found")
			return
		}
		page = int(page64)
	}

	log.WithFields(log.Fields{
		"url":           r.URL.Path,
		"totalItems":    totalItems,
		"perPage":       perPage,
		"numberOfPages": numberOfPages,
		"page":          numberOfPages,
	}).Infof("serving page")

	response := Users{}
	var startingIndex int
	var endIndex int
	if page == 1 {
		startingIndex = 0
	} else {
		startingIndex = (page * perPage) - perPage
	}

	if page == numberOfPages {
		// Last page
		endIndex = len(dbResults)
	} else {
		endIndex = startingIndex + perPage
	}

	for i := startingIndex; i < endIndex; i++ {
		response = append(response, dbResults[i])
	}

	var jsonResponse []byte
	jsonResponse, err = json.Marshal(response)
	if err != nil {
		jsonHTTPResponseWriter(w, 500, "Internal Server Error")
		log.WithError(err).WithFields(log.Fields{"url": r.URL.Path}).Error("marshalling response into JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResponse)
}

func jsonHTTPResponseWriter(w http.ResponseWriter, statusCode int, message string) {
	resp := JsonHTTPResponse{
		Code:    statusCode,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResp, _ := json.Marshal(resp)
	_, _ = w.Write(jsonResp)
}
