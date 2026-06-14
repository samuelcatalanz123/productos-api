package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// secret devuelve la clave con la que firmamos los tokens. (Lección 213)
// Es como el "sello de fábrica" de la pulsera: solo el servidor la conoce.
// En producción debe venir de una variable de entorno (JWT_SECRET), ser larga
// y aleatoria, y NUNCA subirse a Git. Aquí dejamos una de respaldo para
// desarrollo.
func secret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "clave-secreta-de-desarrollo-cambiar-en-produccion"
	}
	return []byte(s)
}

// GenerateToken crea un token JWT para un usuario. Caduca en 24 horas.
// Dentro del token guardamos el id del usuario ("claims" = lo que afirma el token).
func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // caducidad
	}

	// HS256 es el método de firma. Mezcla los datos + la clave secreta.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret())
}

// ParseToken verifica un token y, si es válido, devuelve el id del usuario
// que lleva dentro. Es lo contrario de GenerateToken. (Lección 213)
func ParseToken(tokenString string) (int, error) {
	// jwt.Parse revisa la firma. La función que le pasamos le dice cuál es la
	// clave secreta para comprobarla.
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// Comprobamos que la firma sea del tipo que esperamos (HMAC/HS256).
		// Esto evita un ataque clásico donde alguien cambia el método de firma.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return secret(), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("token inválido o caducado")
	}

	// "claims" son los datos que guardamos dentro del token al crearlo.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("token inválido")
	}

	// Ojo: los números dentro de un JSON llegan como float64, así que el
	// user_id hay que convertirlo a int.
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("token sin user_id")
	}
	return int(userID), nil
}
