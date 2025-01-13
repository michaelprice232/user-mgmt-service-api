#!/bin/bash

if [ -z "${RDS_ENDPOINT}" ]; then
  echo "Error: RDS_ENDPOINT is not set or empty" >&2
  exit 1
fi

if [ -z "${RDS_USERNAME}" ]; then
  echo "Error: RDS_USERNAME is not set or empty" >&2
  exit 1
fi

if [ -z "${PGPASSWORD}" ]; then
  echo "Error: PGPASSWORD is not set or empty" >&2
  exit 1
fi

if [ -z "${DB_NAME}" ]; then
  echo "Error: DB_NAME is not set or empty" >&2
  exit 1
fi


# Create table
if ! psql --host="${RDS_ENDPOINT}" --dbname="${DB_NAME}" --username="${RDS_USERNAME}" --file=/sql-scripts/01-create-table.sql; then
  echo "Problem creating SQL table"
fi

# Import sample rows to test against
if ! psql --host="${RDS_ENDPOINT}" --dbname="${DB_NAME}" --username="${RDS_USERNAME}" --file=/sql-scripts/02-insert-users.sql; then
  echo "Problem inserting rows into table"
fi

echo "SQL commands run successfully"