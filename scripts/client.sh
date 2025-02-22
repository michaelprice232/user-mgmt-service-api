#!/usr/bin/env bash

# Script for running ad-hoc manual tests of the CRUD endpoints. It's recommended to run the integration tests instead, which automate these using Terratest and Docker Compose

set -eu -o pipefail

url='http://localhost:8080'

## Happy paths ##
# GET /users
echo "GET /users (no query params)"
curl -s "${url}/users" | jq
echo

echo  "Test Pagination:  GET /users?per_page=5&page=2"
curl -s "${url}/users?per_page=5&page=2" | jq
echo

echo  "Test Filtering: GET /users?name_filter=bob"
curl -s "${url}/users?name_filter=bob" | jq
echo

echo  "Test Filtering & pagination: GET /users?name_filter=bob&per_page=1&page=2"
curl -s "${url}/users?name_filter=bob&per_page=1&page=1" | jq
echo

# POST /users
echo  "POST /users"
curl -s -X POST "${url}/users" \
  -H 'Content-Type: application/json' \
  -d '{"logon_name":"testuser1","full_name":"Test User 1","email":"test1@email.com"}' | jq
echo

# DELETE /users/<logon_name>
echo  "DELETE /users/clive88"
curl -s -i -X DELETE "${url}/users/clive88"
echo

# PUT /users/<logon_name>
echo  "PUT /users/<logon_name>"
curl -s -X PUT "${url}/users/holly0" \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"Holly Updated","email":"holly.updated@email.com"}' | jq
echo

## Exceptions ##
echo  "per_page param too large: GET /users?per_page=2000"
curl -s "${url}/users?per_page=2000" | jq
echo

echo  "page not found: GET /users?page=1000"
curl -s "${url}/users?page=1000" | jq
echo

echo  "logon_name already taken"
curl -s -X POST "${url}/users" \
  -H 'Content-Type: application/json' \
  -d '{"logon_name":"testuser1","full_name":"Test User 1","email":"test1@email.com"}' | jq
echo

echo  "full_name field too long"
long_full_name='qwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiopqwertyuiop'
curl -s -X POST "${url}/users" \
  -H 'Content-Type: application/json' \
  -d "{\"logon_name\":\"testuser2\",\"full_name\":\"${long_full_name}\",\"email\":\"test1@email.com\"}" | jq
echo

echo  "bad email format"
bad_email_format='email@'
curl -s -X POST "${url}/users" \
  -H 'Content-Type: application/json' \
  -d "{\"logon_name\":\"testuser3\",\"full_name\":\"Test User 3\",\"email\":\"${bad_email_format}\"}" | jq
echo


