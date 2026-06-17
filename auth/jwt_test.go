package auth

import "testing"

// Un token recién creado debe poder leerse y devolver el mismo user_id.
func TestGenerarYLeerToken(t *testing.T) {
	token, err := GenerateToken(42)
	if err != nil {
		t.Fatalf("no esperaba error generando el token: %v", err)
	}
	if token == "" {
		t.Fatal("el token no debería estar vacío")
	}

	id, err := ParseToken(token)
	if err != nil {
		t.Fatalf("no esperaba error leyendo el token: %v", err)
	}
	if id != 42 {
		t.Fatalf("esperaba user_id 42, obtuve %d", id)
	}
}

// Un texto que no es un token debe rechazarse.
func TestTokenBasuraSeRechaza(t *testing.T) {
	if _, err := ParseToken("esto-no-es-un-token"); err == nil {
		t.Fatal("un token basura debería dar error")
	}
}

// Si alguien manipula el token, la firma deja de cuadrar y debe rechazarse.
func TestTokenManipuladoSeRechaza(t *testing.T) {
	token, err := GenerateToken(7)
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	manipulado := token + "x" // rompemos la firma
	if _, err := ParseToken(manipulado); err == nil {
		t.Fatal("un token manipulado debería rechazarse")
	}
}

// Un token firmado con OTRA clave no debe aceptarse (seguridad).
func TestTokenConOtraClaveSeRechaza(t *testing.T) {
	t.Setenv("JWT_SECRET", "clave-A")
	token, err := GenerateToken(1)
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	t.Setenv("JWT_SECRET", "clave-B") // cambiamos la clave del servidor
	if _, err := ParseToken(token); err == nil {
		t.Fatal("un token firmado con otra clave debería rechazarse")
	}
}
