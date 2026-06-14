package response

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"gopkg.in/yaml.v3"
)

// Message es una respuesta sencilla para mensajes o errores.
type Message struct {
	XMLName xml.Name `json:"-" yaml:"-" xml:"respuesta"`
	Mensaje string   `json:"mensaje" xml:"mensaje" yaml:"mensaje"`
}

// Respond escribe la respuesta en el formato pedido. (Lecciones 203-205, 211-212)
// El formato se elige con ?format=xml o ?format=yaml (por defecto: JSON).
func Respond(w http.ResponseWriter, r *http.Request, status int, data any) {
	switch r.URL.Query().Get("format") {
	case "xml":
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(status)
		xml.NewEncoder(w).Encode(data)
	case "yaml":
		w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
		w.WriteHeader(status)
		yaml.NewEncoder(w).Encode(data)
	default:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(data)
	}
}

// Error es un atajo para responder con un mensaje de error.
func Error(w http.ResponseWriter, r *http.Request, status int, mensaje string) {
	Respond(w, r, status, Message{Mensaje: mensaje})
}
