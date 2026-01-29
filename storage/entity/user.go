package entity

type User struct {
	Id int `json:"id"`
	Password []byte `json:"-"`
	Email string `json:"email"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	CreatedAt string `json:"created_at"`
}