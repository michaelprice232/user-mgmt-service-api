package api

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// queryRecordCount returns the count of records based one 1 of 3 filters (only 1 can be used at once)
// 1) records which have a full_name which have a wildcard match against nameFilter
// 2) records which have a logon_name which has an exact match against logonNameFilter
// 3) all records in the users table (no filters)
func (m *UserModel) queryRecordCount(nameFilter, logonNameFilter string) (int, error) {
	var count int
	var row *sql.Row
	var err error

	if nameFilter != "" && logonNameFilter != "" {
		return 0, fmt.Errorf("cannot define both nameFilter and logonNameFilter for queryRecordCount function")
	}

	if nameFilter != "" {
		row = m.DB.QueryRow("SELECT COUNT(*) FROM users WHERE full_name like '%' || $1 || '%'", nameFilter)
	} else if logonNameFilter != "" {
		row = m.DB.QueryRow("SELECT COUNT(*) FROM users WHERE logon_name = $1", logonNameFilter)
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
		rows, err = m.DB.Query(`SELECT user_id, logon_name, full_name, email FROM users WHERE full_name like '%' || $1 || '%' ORDER BY user_id OFFSET $2 LIMIT $3`, nameFilter, offset, limit)
	} else {
		rows, err = m.DB.Query(`SELECT user_id, logon_name, full_name, email FROM users ORDER BY user_id OFFSET $1 LIMIT $2`, offset, limit)
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
		if err = rows.Scan(&user.UserID, &user.LogonName, &user.FullName, &user.Email); err != nil {
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

// addUser adds a new user to the users table
func (m *UserModel) addUser(user User) (User, error) {
	err := m.DB.QueryRow(`INSERT INTO users(logon_name, full_name, email) VALUES ($1, $2, $3) RETURNING user_id`, user.LogonName, user.FullName, user.Email).Scan(&user.UserID)
	if err != nil {
		return user, fmt.Errorf("inserting logon_name '%s' into users table: %v", user.LogonName, err)
	}

	return user, nil
}

// deleteUser deletes a user from the users table
func (m *UserModel) deleteUser(logonName string) error {
	_, err := m.DB.Exec(`DELETE FROM users WHERE logon_name = $1`, logonName)
	if err != nil {
		return fmt.Errorf("deleting record with logon_name = '%s' from users table: %v", logonName, err)
	}

	return nil
}

// updateUser updates a single record in the users table based on the logon_name
// Supports updating email or logon_name fields or both
func (m *UserModel) updateUser(user User) (User, error) {
	log.Debugf("user: %#v", user)
	var err error
	if user.Email != "" && user.FullName != "" {
		err = m.DB.QueryRow(fmt.Sprintf(`UPDATE users SET email = $1, full_name = $2 WHERE logon_name = $3 RETURNING *`), user.Email, user.FullName, user.LogonName).Scan(&user.UserID, &user.LogonName, &user.FullName, &user.Email)
	} else if user.Email != "" {
		err = m.DB.QueryRow(fmt.Sprintf(`UPDATE users SET email = $1 WHERE logon_name = $2 RETURNING *`), user.Email, user.LogonName).Scan(&user.UserID, &user.LogonName, &user.FullName, &user.Email)
	} else if user.FullName != "" {
		err = m.DB.QueryRow(fmt.Sprintf(`UPDATE users SET full_name = $1 WHERE logon_name = $2 RETURNING *`), user.FullName, user.LogonName).Scan(&user.UserID, &user.LogonName, &user.FullName, &user.Email)
	} else {
		return user, fmt.Errorf("email and/or full_name fields need to be set in the user object")
	}

	if err != nil {
		return user, fmt.Errorf("updating record: %v", err)
	}

	return user, nil
}

// OpenDBConnection opens a Postgres DB connection pool
func OpenDBConnection() (*Env, error) {
	EnvConfig = &Env{DBCredentials: DBCredentials{
		HostName:   RequireStringEnvar("database_host_name"),
		Port:       uint(RequireIntEnvar("database_port")),
		DBName:     RequireStringEnvar("database_name"),
		DBUsername: RequireStringEnvar("database_username"),
		DBPassword: RequireStringEnvar("database_password"),
		SSLMode:    RequireStringEnvar("database_ssl_mode")}}

	sqlConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		EnvConfig.DBCredentials.HostName, EnvConfig.DBCredentials.Port, EnvConfig.DBCredentials.DBUsername,
		EnvConfig.DBCredentials.DBPassword, EnvConfig.DBCredentials.DBName, EnvConfig.DBCredentials.SSLMode)

	db, err := sql.Open("postgres", sqlConnection)
	if err != nil {
		return EnvConfig, fmt.Errorf("opening DB connection: %v", err)
	}
	EnvConfig.UsersDB = &UserModel{DB: db}

	return EnvConfig, nil
}
