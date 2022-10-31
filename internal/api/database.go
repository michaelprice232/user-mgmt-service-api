package api

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// queryRecordCount returns the total number of records in the users table, or number of records that match the nameFilter filter (if nameFilter is non-empty)
func (m *UserModel) queryRecordCount(nameFilter string) (int, error) {
	var count int
	var row *sql.Row
	var err error

	if nameFilter != "" {
		row = m.DB.QueryRow("SELECT COUNT(*) FROM users WHERE full_name like '%' || $1 || '%'", nameFilter)
	} else {
		row = m.DB.QueryRow("SELECT COUNT(*) FROM users")
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// queryUsers returns a slice of Users from the users table based on the supplied offset, limit & nameFilter
func (m *UserModel) queryUsers(offset, limit int, nameFilter string) ([]User, error) {
	usersDBResponse := make([]User, 0)
	var err error
	var rows *sql.Rows

	if nameFilter != "" {
		rows, err = m.DB.Query(`SELECT user_id, full_name, email FROM users WHERE full_name like '%' || $1 || '%' ORDER BY user_id OFFSET $2 LIMIT $3`, nameFilter, offset, limit)
	} else {
		rows, err = m.DB.Query(`SELECT user_id, full_name, email FROM users ORDER BY user_id OFFSET $1 LIMIT $2`, offset, limit)
	}
	if err != nil {
		return usersDBResponse, fmt.Errorf("querying database for users: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.WithError(err).Error("closing DB rows response")
		}
	}(rows)

	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.UserID, &user.Name, &user.Email); err != nil {
			return usersDBResponse, fmt.Errorf("scanning over the DB results: %v", err)
		}
		usersDBResponse = append(usersDBResponse, user)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return usersDBResponse, fmt.Errorf("iterating over the DB results: %v", err)
	}

	return usersDBResponse, nil
}
