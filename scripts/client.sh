#!/usr/bin/env bash

set -eu -o pipefail

url='http://localhost:8080'

echo "Query GET /users:"
curl --silent "${url}/users" | jq
echo

echo  "Test Pagination"
curl --silent "${url}/users?per_page=4&page=2" | jq
echo

echo  "Test Filtering"
curl --silent "${url}/users?name_filter=bob" | jq
echo


