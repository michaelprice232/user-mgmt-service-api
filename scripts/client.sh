#!/usr/bin/env bash

set -eu -o pipefail

url='http://localhost:8080'

echo "Query root:"
curl --silent "${url}/"
echo

echo "Query /users:"
curl --silent "${url}/users" | jq
echo

echo  "Pagination"
curl --silent "${url}/users?per_page=4&page=2" | jq
echo

echo  "Filtering"
curl --silent "${url}/users?name_filter=bob" | jq
echo


