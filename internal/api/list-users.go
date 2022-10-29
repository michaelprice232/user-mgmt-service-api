package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
		if perPage64 > maxPageSize {
			jsonHTTPErrorResponseWriter(w, r, 400, "per_page value larger than allowed max")
			return
		}
		perPage = int(perPage64)
	}

	nameFilter = queryStrings.Get("name_filter")

	recordCount, err := queryRecordCount(nameFilter)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, "Internal Server Error")
		return
	}
	fmt.Printf("Total Records: %d\n", recordCount)

	totalItems := recordCount
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
	var startingIndex int
	if page == 1 {
		startingIndex = 0
	} else {
		startingIndex = (page * perPage) - perPage
	}

	if page != numberOfPages {
		response.MorePages = true
	}

	response.TotalPages = numberOfPages
	response.CurrentPage = page

	dbResults, err := queryUsers(startingIndex, perPage, nameFilter)
	if err != nil {
		log.WithError(err).Fatal("querying database")
	}

	response.Users = dbResults

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
