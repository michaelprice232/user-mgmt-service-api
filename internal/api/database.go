package api

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// queryRecordCount returns the total number of records in the users table, or number of records that match the nameFilter filter (if nameFilter is non-empty)
func queryRecordCount(nameFilter string) (int, error) {
	var count int
	var row *sql.Row

	// todo: extract into an ENVAR or AWS Secrets Manager
	connStr := "user=postgres password=test dbname=user-mgmt-db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.WithError(err).Fatal("closing database following record count")
		}
	}(db)

	if nameFilter != "" {
		row = db.QueryRow("SELECT COUNT(*) FROM users WHERE full_name like '%' || $1 || '%'", nameFilter)
	} else {
		row = db.QueryRow("SELECT COUNT(*) FROM users")
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// queryUsers returns a slice of Users from the users table based on the supplied offset, limit & nameFilter
func queryUsers(offset, limit int, nameFilter string) ([]User, error) {
	usersDBResponse := make([]User, 0)
	var err error
	var rows *sql.Rows

	// todo: extract into an ENVAR or AWS Secrets Manager
	connStr := "user=postgres password=test dbname=user-mgmt-db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if nameFilter != "" {
		rows, err = db.Query(`SELECT full_name, email FROM users WHERE full_name like '%' || $1 || '%' ORDER BY user_id OFFSET $2 LIMIT $3`, nameFilter, offset, limit)
	} else {
		rows, err = db.Query(`SELECT full_name, email FROM users ORDER BY user_id OFFSET $1 LIMIT $2`, offset, limit)
	}

	if err != nil {
		return usersDBResponse, fmt.Errorf("querying database: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.WithError(err).Fatal("closing DB rows response")
		}
	}(rows)

	user := User{}
	for rows.Next() {
		if err = rows.Scan(&user.Name, &user.Email); err != nil {
			return usersDBResponse, fmt.Errorf("scanning the DB results: %v", err)
		}
		usersDBResponse = append(usersDBResponse, user)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return usersDBResponse, fmt.Errorf("iterating over the DB results: %v", err)
	}

	return usersDBResponse, nil
}
