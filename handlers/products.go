package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"productos-api/models"
	"productos-api/repository"
	"productos-api/response"
)

// ProductHandler agrupa los handlers y guarda la conexión a la base de datos.
type ProductHandler struct {
	DB *sql.DB
}

// List responde con todos los productos. GET /productos  (Lección 202)
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	productos, err := repository.List(h.DB)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "no se pudo listar")
		return
	}
	response.Respond(w, r, http.StatusOK, models.ProductList{Productos: productos})
}

// Get responde con un producto por su id. GET /productos/{id}  (Lección 206)
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromURL(w, r)
	if !ok {
		return
	}

	producto, err := repository.Get(h.DB, id)
	if err == sql.ErrNoRows {
		response.Error(w, r, http.StatusNotFound, "producto no encontrado")
		return
	}
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "error al obtener")
		return
	}
	response.Respond(w, r, http.StatusOK, producto)
}

// Create crea un producto nuevo. POST /productos  (Lección 207)
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.Error(w, r, http.StatusBadRequest, "datos inválidos")
		return
	}

	id, err := repository.Create(h.DB, p)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "no se pudo crear")
		return
	}
	p.ID = int(id)
	response.Respond(w, r, http.StatusCreated, p)
}

// Update actualiza un producto. PUT /productos/{id}  (Lección 208)
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromURL(w, r)
	if !ok {
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.Error(w, r, http.StatusBadRequest, "datos inválidos")
		return
	}
	p.ID = id

	filas, err := repository.Update(h.DB, p)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "no se pudo actualizar")
		return
	}
	if filas == 0 {
		response.Error(w, r, http.StatusNotFound, "producto no encontrado")
		return
	}
	response.Respond(w, r, http.StatusOK, p)
}

// Delete borra un producto. DELETE /productos/{id}  (Lección 209)
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromURL(w, r)
	if !ok {
		return
	}

	filas, err := repository.Delete(h.DB, id)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "no se pudo eliminar")
		return
	}
	if filas == 0 {
		response.Error(w, r, http.StatusNotFound, "producto no encontrado")
		return
	}
	response.Error(w, r, http.StatusOK, "producto eliminado")
}

// idFromURL saca el id de la ruta (/productos/{id}) y lo convierte a número.
func idFromURL(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "id inválido")
		return 0, false
	}
	return id, true
}
