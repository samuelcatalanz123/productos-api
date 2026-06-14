package repository

import (
	"database/sql"

	"productos-api/models"
)

// CreateTable crea la tabla productos si no existe.
func CreateTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS productos (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL,
			precio DECIMAL(10,2) NOT NULL
		)`
	_, err := db.Exec(query)
	return err
}

// List devuelve todos los productos.
func List(db *sql.DB) ([]models.Product, error) {
	rows, err := db.Query("SELECT id, nombre, precio FROM productos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productos := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Precio); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, nil
}

// Get busca un producto por su id.
func Get(db *sql.DB, id int) (models.Product, error) {
	var p models.Product
	err := db.QueryRow(
		"SELECT id, nombre, precio FROM productos WHERE id = ?", id,
	).Scan(&p.ID, &p.Nombre, &p.Precio)
	return p, err
}

// Create inserta un producto y devuelve su nuevo id.
func Create(db *sql.DB, p models.Product) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO productos (nombre, precio) VALUES (?, ?)",
		p.Nombre, p.Precio,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update modifica un producto existente. Devuelve cuántas filas cambiaron.
func Update(db *sql.DB, p models.Product) (int64, error) {
	result, err := db.Exec(
		"UPDATE productos SET nombre = ?, precio = ? WHERE id = ?",
		p.Nombre, p.Precio, p.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete borra un producto por su id. Devuelve cuántas filas se borraron.
func Delete(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM productos WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
