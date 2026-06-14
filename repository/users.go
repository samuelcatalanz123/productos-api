package repository

import (
	"database/sql"

	"productos-api/models"
)

// CreateUsersTable crea la tabla de usuarios si no existe. (Lección 213)
// El email es UNIQUE: no puede haber dos usuarios con el mismo email.
// La columna password guarda el HASH (texto encriptado), no la contraseña real,
// por eso necesita 255 caracteres.
func CreateUsersTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS usuarios (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL
		)`
	_, err := db.Exec(query)
	return err
}

// CreateUser inserta un usuario nuevo y devuelve su id.
// Recibe el hash ya encriptado (el handler se encarga de encriptar).
func CreateUser(db *sql.DB, email, hash string) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO usuarios (email, password) VALUES (?, ?)",
		email, hash,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserByEmail busca un usuario por su email (lo usaremos en el login).
func GetUserByEmail(db *sql.DB, email string) (models.User, error) {
	var u models.User
	err := db.QueryRow(
		"SELECT id, email, password FROM usuarios WHERE email = ?", email,
	).Scan(&u.ID, &u.Email, &u.Password)
	return u, err
}
