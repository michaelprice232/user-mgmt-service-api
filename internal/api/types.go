package api

import "database/sql"

type Env struct {
	UsersDB interface {
		queryRecordCount(string) (int, error)
		queryUsers(int, int, string) ([]User, error)
	}
}

type UserModel struct {
	DB *sql.DB
}

type User struct {
	UserID int
	Name   string
	Email  string
}

type UsersResponse struct {
	Users       []User
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	MorePages   bool `json:"more_pages"`
}

type JsonHTTPErrorResponse struct {
	Code    int
	Message string
}

type queryParameters struct {
	perPage    int
	page       int
	nameFilter string
}
