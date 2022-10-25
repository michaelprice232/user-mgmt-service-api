package api

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func queryAllUsers() ([]User, error) {
	usersDBResponse := make([]User, 0)
	var err error

	// todo: pull username/password from AWS Secrets Manager
	connStr := "user=postgres password=test dbname=user-mgmt-db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT full_name, email FROM users")
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
