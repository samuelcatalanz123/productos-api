package models

// User representa un usuario que puede iniciar sesión. (Lección 213 - JWT)
// La etiqueta `json:"password,omitempty"` hace que NUNCA enviemos la
// contraseña de vuelta al cliente cuando esté vacía.
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}
