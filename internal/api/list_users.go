package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	defaultPageSize = 2
	maxPageSize     = 5
)

func (env *Env) listUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	var params queryParameters
	var recordCount int
	var startingIndex int
	var dbResults []User

	queryStrings := r.URL.Query()
	params, err = extractAndValidateQueryParams(w, r, queryStrings)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, "processing query parameters")
		return
	}

	recordCount, err = env.UsersDB.queryRecordCount(params.nameFilter)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, "calculating the number of records in database")
		return
	}

	numberOfPages := recordCount / params.perPage
	if recordCount%params.perPage != 0 {
		// Add a non-full page
		numberOfPages++
	}

	// Can only be performed after the number of records is obtained and so can't be part of the extractAndValidateQueryParams function
	if params.page > numberOfPages {
		jsonHTTPErrorResponseWriter(w, r, 404, fmt.Sprintf("page %d not found", params.page))
		return
	}

	response := UsersResponse{}
	if params.page == 1 {
		startingIndex = 0
	} else {
		startingIndex = (params.page * params.perPage) - params.perPage
	}

	if params.page != numberOfPages {
		response.MorePages = true
	}

	response.TotalPages = numberOfPages
	response.CurrentPage = params.page
	dbResults, err = env.UsersDB.queryUsers(startingIndex, params.perPage, params.nameFilter)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, "querying the users table")
		return
	}

	response.Users = dbResults

	err = writeJSONHTTPResponse(w, response)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("writing HTTP response: %v", err))
		return
	}

	log.WithFields(log.Fields{
		"url":           r.URL.Path,
		"totalItems":    recordCount,
		"numberOfPages": numberOfPages,
		"perPage":       params.perPage,
		"page":          params.page,
		"statusCode":    200,
	}).Infof("serving page")
}

func extractAndValidateQueryParams(w http.ResponseWriter, r *http.Request, queryStrings url.Values) (queryParameters, error) {
	var err error
	var perPage64, page64 int64
	var params queryParameters

	if perPageEnv := queryStrings.Get("per_page"); perPageEnv != "" {
		perPage64, err = strconv.ParseInt(perPageEnv, 10, 64)
		if err != nil || perPage64 <= 0 || perPage64 > maxPageSize {
			jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("Invalid per_page value. Must be an integer between 1 -> %d", maxPageSize))
			return params, err
		}
		params.perPage = int(perPage64)
	} else {
		params.perPage = defaultPageSize
	}

	if pageEnv := queryStrings.Get("page"); pageEnv != "" {
		page64, err = strconv.ParseInt(pageEnv, 10, 64)
		if err != nil || page64 <= 0 {
			jsonHTTPErrorResponseWriter(w, r, 404, fmt.Sprintf("page %s not found", pageEnv))
			return params, err
		}
		params.page = int(page64)
	} else {
		params.page = 1
	}

	params.nameFilter = queryStrings.Get("name_filter")

	return params, nil
}

func writeJSONHTTPResponse(w http.ResponseWriter, payload UsersResponse) error {
	var err error
	var jsonResponse []byte

	jsonResponse, err = json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling JSON in preparation for HTTP response")
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		return err
	}
	return nil
}
