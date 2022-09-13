package api

type User struct {
	Name  string
	Email string
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
