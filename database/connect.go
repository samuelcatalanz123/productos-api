package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Connect abre y comprueba la conexión con la base de datos 'tienda'.
// La dirección de conexión (DSN) se puede cambiar con la variable de entorno
// DB_DSN; así funciona igual en local y dentro de Docker. Si no está, usa la
// de desarrollo local.
func Connect() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root@tcp(localhost:3306)/tienda"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
