# user-mgmt-service-api

[![CI](https://github.com/michaelprice232/user-mgmt-service-api/actions/workflows/ci.yml/badge.svg)](https://github.com/michaelprice232/user-mgmt-service-api/actions/workflows/ci.yml) [![PR against main](https://github.com/michaelprice232/user-mgmt-service-api/actions/workflows/pull-request.yml/badge.svg)](https://github.com/michaelprice232/user-mgmt-service-api/actions/workflows/pull-request.yml)

L&D project containing a user management REST API exposing CRUD endpoints, written in Go with a Postgres DB.

## Running locally

Pre-reqs:
- Docker & docker-compose installed
- Go installed (v1.23 or above)

Steps:
```shell
# Run the app & database using Docker Compose
# You can then curl the endpoints e.g. http://localhost:8080/users
make run

# Cleanup the above running containers
make down

# Run unit tests
make unit-tests

# Run integration tests. Uses Docker Compose to deploy the app and then Terratest
make int-tests

# Run the E2E tests. Deploys the infra using Terraform into AWS and then uses Terratest to verify it
# Requires that one the AWS credentials chain is present e.g. AWS_PROFILE env var is set
make e2e-tests
```

## CI (GitHub Actions)

- Push to any branch will trigger the linter (TODO), unit tests and integration tests (Docker Compose)
- PR against the main branch will run the above tests as well as the E2E tests (deploying into AWS). Docker images in ECR will be based off the Git branch SHAs
- When the PR is merged into the main branch this will auto create a release PR using a Semantic Version tag via [release please](https://github.com/googleapis/release-please-action)
- Once the release PR is merged the GitHub release will be created, build binaries attached and Docker images built in ECR using the SemVer tags

## Endpoints

| Endpoint                   | Description                                                                                                                                                       | Query Strings                                                                         | Request Payload Type | Response Payload Type                    | 
|----------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|----------------------|------------------------------------------|
| GET /users                 | List the users in the database. Supports pagination and filtering by name                                                                                         | **per_page**: how many users to display in each returned page                         | N/A (no payload)     | UsersResponse                            |
|                            |                                                                                                                                                                   | **page**: page number to return                                                       |                      |                                          |
|                            |                                                                                                                                                                   | **name_filter**: return users which have a full_name which match this wildcard search |                      |                                          |
| POST /users                | Add a new user. User logon_name must be unique. user_id is auto generated and cannot be passed in the request payload                                             | N/A                                                                                   | User                 | User                                     |
| DELETE /users/<logon_name> | Delete a user from the database based on their logon_name                                                                                                         | N/A                                                                                   | N/A                  | N/A                                      |
| PUT /users/<logon_name>    | Update an existing user. Supports the full_name & email fields or both                                                                                            | N/A                                                                                   | User                 | User                                     |
| GET /health                | Health endpoint for use by K8s readiness/liveness probes. Currently polls the database. Utilises the [health-go library](https://github.com/hellofresh/health-go) | N/A                                                                                   | N/A                  | github.com/hellofresh/health-go/v5/Check |


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

# Update a user
% curl -s -X PUT "${url}/users/holly0" \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"Holly Updated","email":"holly.updated@email.com"}' | jq
{
  "user_id": 6,
  "logon_name": "holly0",
  "full_name": "Holly Updated",
  "email": "holly.updated@email.com"
}

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

## TODO
- [ ] Add GitHub workflow for auto releasing on merge into main
- [ ] Add a GET endpoint for an individual user
- [ ] Use primary keys in the REST URI's rather than logon names
- [ ] Instrument with Prometheus library
- [ ] Instrument with OpenTelemetry client
- [ ] Add OpenAPI docs
- [ ] Replace Gorilla Mux module with standard library HTTP routing functionality