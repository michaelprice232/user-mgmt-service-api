package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	defaultPageSize = 2
	maxPageSize     = 5
)

func listUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	var perPage = defaultPageSize
	var perPage64, page64 int64
	var page = 1
	var nameFilter string

	queryStrings := r.URL.Query()
	if perPageEnv := queryStrings.Get("per_page"); perPageEnv != "" {
		perPage64, err = strconv.ParseInt(perPageEnv, 10, 64)
		if err != nil || perPage64 <= 0 {
			jsonHTTPErrorResponseWriter(w, r, 400, "Invalid per_page value")
			return
		}
		perPage = int(perPage64)

		if perPage > maxPageSize {
			jsonHTTPErrorResponseWriter(w, r, 400, "per_page value larger than allowed max")
			return
		}
	}

	// todo: read from DB instead of this duplication
	var dbResults UsersResponse
	if nameFilter = queryStrings.Get("name_filter"); nameFilter != "" {
		for _, user := range usersDB.Users {
			if strings.Contains(user.Name, nameFilter) {
				dbResults.Users = append(dbResults.Users, user)
			}
		}
	} else {
		dbResults = usersDB
	}

	totalItems := len(dbResults.Users)
	numberOfPages := totalItems / perPage
	if totalItems%perPage != 0 {
		// Add a non-full page
		numberOfPages++
	}

	if pageEnv := queryStrings.Get("page"); pageEnv != "" {
		page64, err = strconv.ParseInt(pageEnv, 10, 64)
		if err != nil || page64 <= 0 || page64 > int64(numberOfPages) {
			jsonHTTPErrorResponseWriter(w, r, 404, "Page not found")
			return
		}
		page = int(page64)
	}

	response := UsersResponse{}
	var startingIndex, endIndex int
	if page == 1 {
		startingIndex = 0
	} else {
		startingIndex = (page * perPage) - perPage
	}

	if page == numberOfPages {
		// Last page
		endIndex = len(dbResults.Users)
	} else {
		endIndex = startingIndex + perPage
		response.MorePages = true
	}

	for i := startingIndex; i < endIndex; i++ {
		response.Users = append(response.Users, dbResults.Users[i])
	}

	response.TotalPages = numberOfPages
	response.CurrentPage = page

	var jsonResponse []byte
	jsonResponse, err = json.Marshal(response)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.WithFields(log.Fields{"url": r.URL.Path}).WithError(err).Error("writing HTTP response")
	}

	log.WithFields(log.Fields{
		"url":           r.URL.Path,
		"totalItems":    totalItems,
		"numberOfPages": numberOfPages,
		"perPage":       perPage,
		"page":          page,
		"statusCode":    200,
	}).Infof("serving page")
}

// todo: read from database
var usersDB = UsersResponse{
	Users: []User{
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
	},
}
