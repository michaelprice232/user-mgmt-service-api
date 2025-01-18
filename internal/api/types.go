package api

import "database/sql"

type Env struct {
	UsersDB interface {
		queryRecordCount(string, string) (int, error)
		queryUsers(int, int, string) ([]User, error)
		addUser(User) (User, error)
		deleteUser(string) error
		updateUser(User) (User, error)
	}
	DBCredentials DBCredentials
	BuildVersion  string
}

type DBCredentials struct {
	HostName   string
	Port       int64
	DBName     string
	DBUsername string
	DBPassword string
	SSLMode    string
}

type UserModel struct {
	DB *sql.DB
}

type User struct {
	UserID    int    `json:"user_id,omitempty"`
	LogonName string `json:"logon_name"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
}

type UsersResponse struct {
	Users       []User
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	MorePages   bool `json:"more_pages"`
}

type JSONHTTPErrorResponse struct {
	Code    int
	Message string
}

type queryParameters struct {
	perPage    int
	page       int
	nameFilter string
}
