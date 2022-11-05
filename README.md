# user-mgmt-service-api

Work in progress

L&D project containing a user management REST API exposing CRUD endpoints, written in Go with a Postgres DB

For request/response/error models please see [types](internal/api/types.go) or see [example output below](https://github.com/michaelprice232/user-mgmt-service-api#example-output). Currently supported endpoints:

| Endpoint                   | Description                                                               | Query Strings                                                                         | Request Payload Type | Response Payload Type | 
|----------------------------|---------------------------------------------------------------------------|---------------------------------------------------------------------------------------|----------------------|-----------------------|
| GET /users                 | List the users in the database. Supports pagination and filtering by name | **per_page**: how many users to display in each returned page                         | N/A (no payload)     | UsersResponse         |
|                            |                                                                           | **page**: page number to return                                                       |                      |                       |
|                            |                                                                           | **name_filter**: return users which have a full_name which match this wildcard search |                      |                       |
| POST /users                | Add a new user. User logon_name must be unique                            | N/A                                                                                   | User                 | User                  |
| DELETE /users/<logon_name> | Delete a user from the database based on their logon_name                 | N/A                                                                                   | N/A                  | N/A                   |


## How to run

Pre-reqs:
- Docker & docker-compose installed locally
- Go installed (v1.18 or above)

Steps:
```shell
# Start the Postgres DB (seeds DB & records during startup), builds and starts the Go webserver
make run-webserver

# Run some curl commands for testing the endpoints. Will be updated as more endpoints are added
make test-endpoints

# Stop docker-compose Postgres database and remove the Docker volume so that the DB init scripts are run next time
# Make sure to also stop the Go webserver process. It is not running as a Docker container yet
make stop-database-and-delete-volume
```

## Unit Tests
```shell
make test
```

## Example Output

```shell
% url='http://localhost:8080'

# Add a new user
% curl -s -X POST "${url}/users" \
  -H 'Content-Type: application/json' \
  -d '{"logon_name":"testuser1","full_name":"Test User 1","email":"test1@email.com"}' | jq
{
  "user_id": 11,
  "logon_name": "testuser1",
  "full_name": "Test User 1",
  "email": "test1@email.com"
}

#  Delete a user
% curl -s -i -X DELETE "${url}/users/susan9"
HTTP/1.1 204 No Content

# Listing all users 
% curl --silent "${url}/users" | jq
{
  "Users": [
    {
      "user_id": 1,
      "logon_name": "mike1",
      "full_name": "mike",
      "email": "mike@email.com"
    },
    {
      "user_id": 2,
      "logon_name": "bob44",
      "full_name": "bob",
      "email": "bob@email.com"
    }
  ],
  "total_pages": 5,
  "current_page": 1,
  "more_pages": true
}

# Pagination
% curl --silent "${url}/users?per_page=4&page=2" | jq
{
  "Users": [
    {
      "user_id": 5,
      "logon_name": "susan9",
      "full_name": "susan",
      "email": "susan@email.com"
    },
    {
      "user_id": 6,
      "logon_name": "holly0",
      "full_name": "holly",
      "email": "holly@email.com"
    },
    {
      "user_id": 7,
      "logon_name": "bobby8",
      "full_name": "bobby",
      "email": "bobby@email.com"
    },
    {
      "user_id": 8,
      "logon_name": "clive88",
      "full_name": "clive",
      "email": "clive@email.com"
    }
  ],
  "total_pages": 3,
  "current_page": 2,
  "more_pages": true
}



# Filtering
% curl --silent "${url}/users?name_filter=bob" | jq
{
  "Users": [
    {
      "user_id": 2,
      "logon_name": "bob44",
      "full_name": "bob",
      "email": "bob@email.com"
    },
    {
      "user_id": 7,
      "logon_name": "bobby8",
      "full_name": "bobby",
      "email": "bobby@email.com"
    }
  ],
  "total_pages": 1,
  "current_page": 1,
  "more_pages": false
}


# Example Validation
% curl --silent "${url}/users?per_page=2000" | jq
{
  "Code": 400,
  "Message": "processing query parameters: per_page query string must be an integer between 1->5"
}

% curl -s -X POST "${url}/users" \
  -H 'Content-Type: application/json' \
  -d '{"logon_name":"testuser1","full_name":"Test User 1","email":"test1@email.com"}' | jq
{
  "Code": 400,
  "Message": "logon_name 'testuser1' already taken. Please choose another one"
}
```

## Remaining Tasks
- [x] Add GET /users
- [x] Add POST /users
- [x] Add DELETE /users/<user> endpoint
- [ ] Add PUT /users/<user> endpoint
- [ ] Add health endpoint suitable for K8s
- [ ] Enable graceful shutdowns of HTTP server suitable for K8s
- [ ] Instrument with Prometheus library
- [ ] Instrument with OpenTelemetry client
- [ ] Integrate with GitHub Actions for running unit tests, linters & security scanner
- [ ] Add Terraform for deploying into K8s cluster
- [ ] Add Terratest smoke tests for validating deployment
