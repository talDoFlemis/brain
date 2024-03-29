package domain

type CreateUserRequest struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}

type UpdateUserRequest struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}
