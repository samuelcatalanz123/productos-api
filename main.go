package main

import (
	"fmt"
	"log"
	"net/http"

	"productos-api/auth"
	"productos-api/database"
	"productos-api/handlers"
	"productos-api/repository"
)

func main() {
	// Conexión a MySQL.
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}
	defer db.Close()

	// Nos aseguramos de que las tablas existan.
	if err := repository.CreateTable(db); err != nil {
		log.Fatal("Error al crear la tabla de productos:", err)
	}
	if err := repository.CreateUsersTable(db); err != nil {
		log.Fatal("Error al crear la tabla de usuarios:", err)
	}

	// Creamos los handlers con la conexión dentro.
	h := &handlers.ProductHandler{DB: db}
	authHandler := &handlers.AuthHandler{DB: db}

	// Definimos las rutas (endpoints) de la API. (Lección 201)
	mux := http.NewServeMux()
	// Leer está abierto a todos.
	mux.HandleFunc("GET /productos", h.List)
	mux.HandleFunc("GET /productos/{id}", h.Get)

	// Modificar requiere token: envolvemos el handler con auth.RequireAuth.
	mux.HandleFunc("POST /productos", auth.RequireAuth(h.Create))
	mux.HandleFunc("PUT /productos/{id}", auth.RequireAuth(h.Update))
	mux.HandleFunc("DELETE /productos/{id}", auth.RequireAuth(h.Delete))

	// Rutas de autenticación. (Lección 213 - JWT)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	fmt.Println("🌐 API escuchando en http://localhost:8080")
	fmt.Println("   Prueba: http://localhost:8080/productos")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
