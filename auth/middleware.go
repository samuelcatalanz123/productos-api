package auth

import (
	"net/http"
	"strings"

	"productos-api/response"
)

// RequireAuth es un MIDDLEWARE: recibe una función (la ruta real) y devuelve
// otra función que primero hace de "portero". (Lección 213)
//
// Se usa así en main.go:
//   mux.HandleFunc("POST /productos", auth.RequireAuth(h.Create))
//
// Si el token es válido, deja pasar a h.Create. Si no, corta y responde 401.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Leemos la cabecera. El cliente la manda así:
		//    Authorization: Bearer eyJhbGciOi...
		header := r.Header.Get("Authorization")
		if header == "" {
			response.Error(w, r, http.StatusUnauthorized, "falta el token (cabecera Authorization)")
			return
		}

		// 2. Separamos "Bearer" del token en sí.
		partes := strings.SplitN(header, " ", 2)
		if len(partes) != 2 || partes[0] != "Bearer" {
			response.Error(w, r, http.StatusUnauthorized, "formato inválido (usa: Bearer <token>)")
			return
		}

		// 3. Verificamos el token con la función del Paso anterior.
		if _, err := ParseToken(partes[1]); err != nil {
			response.Error(w, r, http.StatusUnauthorized, "token inválido o caducado")
			return
		}

		// 4. Todo en orden: dejamos pasar a la ruta real.
		next(w, r)
	}
}
