package user

type CreateUserDTO struct {
	Name     string
	Email    string
	Password string
}

type LoginUserDTO struct {
	Email    string
	Password string
}
