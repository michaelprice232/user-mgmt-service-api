#!/usr/bin/env bash

set -eu -o pipefail

url='http://localhost:8080'

echo "Query GET /users:"
curl --silent "${url}/users" | jq
echo

echo  "Test Pagination (query params: per_page=4&page=2)"
curl --silent "${url}/users?per_page=4&page=2" | jq
echo

echo  "Test Filtering (query params: name_filter=bob)"
curl --silent "${url}/users?name_filter=bob" | jq
echo

echo  "Test Filtering & pagination (query params: name_filter=bob&per_page=1&page=2)"
curl --silent "${url}/users?name_filter=bob&per_page=1&page=2" | jq
echo


# Bad paths
echo  "per_page param too large"
curl --silent "${url}/users?per_page=2000" | jq
echo