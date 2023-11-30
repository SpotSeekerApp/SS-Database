package types

type UserRequest struct {
	Id      string `json:"user_id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Born    string `json:"born"`
	Email   string `json:"email"`
}
