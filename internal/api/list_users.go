package api

import (
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

// listUsers is an HTTP handler got GET /users
func (env *Env) listUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	var params queryParameters
	var recordCount int
	var startingIndex int
	var dbResults []User

	queryStrings := r.URL.Query()
	params, err = extractAndValidateQueryParams(queryStrings)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 400, fmt.Sprintf("processing query parameters: %v", err))
		return
	}

	recordCount, err = env.UsersDB.queryRecordCount(params.nameFilter, "")
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("calculating the number of records in database: %v", err))
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
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("querying the users table: %v", err))
		return
	}

	response.Users = dbResults

	err = writeJSONHTTPResponse(w, response)
	if err != nil {
		jsonHTTPErrorResponseWriter(w, r, 500, fmt.Sprintf("writing HTTP response: %v", err))
		return
	}

	log.WithFields(log.Fields{
		"url":           getFullPathIncludingQueryParams(r.URL),
		"totalItems":    recordCount,
		"numberOfPages": numberOfPages,
		"perPage":       params.perPage,
		"page":          params.page,
		"statusCode":    200,
	}).Infof("serving page")
}

// getFullPathIncludingQueryParams returns either uri, or the uri including the query parameters if they are present
func getFullPathIncludingQueryParams(url *url.URL) string {
	if url.Query().Encode() != "" {
		return fmt.Sprintf("%s?%s", url.Path, url.Query().Encode())
	} else {
		return url.Path
	}
}

// extractAndValidateQueryParams extracts any query strings and validates them
func extractAndValidateQueryParams(queryStrings url.Values) (queryParameters, error) {
	var err error
	var perPage64, page64 int64
	var params queryParameters

	if perPageEnv := queryStrings.Get("per_page"); perPageEnv != "" {
		perPage64, err = strconv.ParseInt(perPageEnv, 10, 64)
		if err != nil || perPage64 <= 0 || perPage64 > maxPageSize {
			if err != nil {
				return params, fmt.Errorf("per_page query string must be an integer between 1->%d: %v", maxPageSize, err)
			} else {
				return params, fmt.Errorf("per_page query string must be an integer between 1->%d", maxPageSize)
			}

		}
		params.perPage = int(perPage64)
	} else {
		params.perPage = defaultPageSize
	}

	if pageEnv := queryStrings.Get("page"); pageEnv != "" {
		page64, err = strconv.ParseInt(pageEnv, 10, 64)
		if err != nil || page64 <= 0 {
			if err != nil {
				return params, fmt.Errorf("page query string must be an integer greater than 0: %v", err)
			} else {
				return params, fmt.Errorf("page query string must be an integer greater than 0")
			}

		}
		params.page = int(page64)
	} else {
		params.page = 1
	}

	params.nameFilter = queryStrings.Get("name_filter")

	return params, nil
}
