#!/bin/bash

if [ -z "${HOSTNAME}" ]; then
  echo "Error: HOSTNAME is not set or empty" >&2
  exit 1
fi

if [ -z "${USERNAME}" ]; then
  echo "Error: USERNAME is not set or empty" >&2
  exit 1
fi

if [ -z "${PGPASSWORD}" ]; then
  echo "Error: PASSWORD is not set or empty" >&2
  exit 1
fi

if [ -z "${DB_NAME}" ]; then
  echo "Error: DB_NAME is not set or empty" >&2
  exit 1
fi

sleep 40

# Create table
if ! psql --host="${HOSTNAME}" --dbname="${DB_NAME}" --username="${USERNAME}" --file=/sql-scripts/01-create-table.sql; then
  echo "Problem creating SQL table"
fi

# Import sample rows to test against
if ! psql --host="${HOSTNAME}" --dbname="${DB_NAME}" --username="${USERNAME}" --file=/sql-scripts/02-insert-users.sql; then
  echo "Problem inserting rows into table"
fi

echo "SQL commands run successfully"