package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"productos-api/auth"
	"productos-api/models"
	"productos-api/repository"
	"productos-api/response"

	"golang.org/x/crypto/bcrypt"
)

// AuthHandler agrupa los handlers de autenticación (registro y login).
type AuthHandler struct {
	DB *sql.DB
}

// Register crea un usuario nuevo. POST /register  (Lección 213)
// Pasos:
//  1. Leer el email y la contraseña que manda el cliente.
//  2. Encriptar la contraseña con bcrypt (NUNCA se guarda en texto plano).
//  3. Guardar el usuario en la base de datos.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		response.Error(w, r, http.StatusBadRequest, "datos inválidos")
		return
	}

	if u.Email == "" || u.Password == "" {
		response.Error(w, r, http.StatusBadRequest, "email y password son obligatorios")
		return
	}

	// bcrypt convierte "miclave123" en algo como "$2a$10$N9qo8uLO...".
	// Es de una sola vía: se puede comprobar, pero NO se puede desencriptar.
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "no se pudo encriptar la contraseña")
		return
	}

	id, err := repository.CreateUser(h.DB, u.Email, string(hash))
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "no se pudo crear el usuario (¿email repetido?)")
		return
	}

	u.ID = int(id)
	u.Password = "" // borramos la contraseña para no devolverla en la respuesta
	response.Respond(w, r, http.StatusCreated, u)
}

// Login comprueba email + contraseña y devuelve un token JWT. POST /login (Lección 213)
// Pasos:
//  1. Buscar al usuario por su email.
//  2. Comparar la contraseña recibida con el hash guardado (bcrypt).
//  3. Si todo cuadra, generar y devolver el token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var datos models.User
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		response.Error(w, r, http.StatusBadRequest, "datos inválidos")
		return
	}

	// 1. Buscamos al usuario.
	usuario, err := repository.GetUserByEmail(h.DB, datos.Email)
	if err != nil {
		// Nota de seguridad: damos el MISMO mensaje si el email no existe o si
		// la contraseña está mal. Así no revelamos qué emails están registrados.
		response.Error(w, r, http.StatusUnauthorized, "email o contraseña incorrectos")
		return
	}

	// 2. bcrypt vuelve a encriptar la contraseña recibida y la compara con el
	//    hash guardado. Devuelve nil (sin error) solo si coinciden.
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(datos.Password)); err != nil {
		response.Error(w, r, http.StatusUnauthorized, "email o contraseña incorrectos")
		return
	}

	// 3. Credenciales correctas: fabricamos el token.
	token, err := auth.GenerateToken(usuario.ID)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "no se pudo generar el token")
		return
	}

	response.Respond(w, r, http.StatusOK, map[string]string{"token": token})
}
