# user-mgmt-service-api

Work in progress

L&D project containing a user management REST API exposing CRUD endpoints, written in Go with a Postgres DB

Currently consists of only a single endpoint `GET /users` which supports `pagination` and `filtering` by name

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
```
