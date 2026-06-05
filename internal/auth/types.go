package auth

type RegisterRequest struct {
	Email    string
	Password string
}

type LoginRequest struct {
	Email string
	Password string 
}

