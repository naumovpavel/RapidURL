package entity

type User struct {
	Id       int
	Name     string
	Password string
	Salt     string
	Email    string
}
