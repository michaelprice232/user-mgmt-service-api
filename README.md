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

## Example Output

```shell
% url='http://localhost:8080'

# Listing all users 
% curl --silent "${url}/users" | jq
{
  "Users": [
    {
      "Name": "mike",
      "Email": "mike@email.com"
    },
    {
      "Name": "bob",
      "Email": "bob@email.com"
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
      "Name": "susan",
      "Email": "susan@email.com"
    },
    {
      "Name": "holly",
      "Email": "holly@email.com"
    },
    {
      "Name": "bobby",
      "Email": "bobby@email.com"
    },
    {
      "Name": "clive",
      "Email": "clive@email.com"
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
      "Name": "bob",
      "Email": "bob@email.com"
    },
    {
      "Name": "bobby",
      "Email": "bobby@email.com"
    }
  ],
  "total_pages": 1,
  "current_page": 1,
  "more_pages": false
}
```
