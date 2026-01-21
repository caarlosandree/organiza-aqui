package dto

// LoginRequest representa a requisição de login
// Permite login com email ou username
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // email ou username
	Password   string `json:"password" validate:"required,min=6"`
}

// RegisterRequest representa a requisição de registro
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=8,password_strength"`
	Name     string `json:"name" validate:"required,min=3,max=255"`
}

// UsernameSuggestionsResponse representa sugestões de username quando já está em uso
type UsernameSuggestionsResponse struct {
	Suggestions []string `json:"suggestions"`
}

// LoginResponse representa a resposta de login
type LoginResponse struct {
	Token string    `json:"token"`
	User  UserDTO   `json:"user"`
}

// UserDTO representa um usuário na camada de apresentação
type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
