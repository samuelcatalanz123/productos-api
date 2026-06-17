package response

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Por defecto (sin ?format) responde en JSON, con el status correcto.
func TestRespondJSONPorDefecto(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	Respond(rec, req, http.StatusCreated, Message{Mensaje: "hola"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("esperaba JSON, content-type fue %q", ct)
	}
	if !strings.Contains(rec.Body.String(), `"mensaje":"hola"`) {
		t.Fatalf("el cuerpo no tiene el mensaje en JSON: %s", rec.Body.String())
	}
}

// Con ?format=xml responde en XML.
func TestRespondXML(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?format=xml", nil)
	rec := httptest.NewRecorder()

	Respond(rec, req, http.StatusOK, Message{Mensaje: "hola"})

	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "xml") {
		t.Fatalf("esperaba XML, content-type fue %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "<mensaje>hola</mensaje>") {
		t.Fatalf("el cuerpo no tiene el mensaje en XML: %s", rec.Body.String())
	}
}

// Con ?format=yaml responde en YAML.
func TestRespondYAML(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?format=yaml", nil)
	rec := httptest.NewRecorder()

	Respond(rec, req, http.StatusOK, Message{Mensaje: "hola"})

	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "yaml") {
		t.Fatalf("esperaba YAML, content-type fue %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "mensaje: hola") {
		t.Fatalf("el cuerpo no tiene el mensaje en YAML: %s", rec.Body.String())
	}
}

// El atajo Error responde con el status y el mensaje dados.
func TestErrorHelper(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	Error(rec, req, http.StatusBadRequest, "algo falló")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "algo falló") {
		t.Fatalf("el cuerpo no tiene el mensaje de error: %s", rec.Body.String())
	}
}
