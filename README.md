# API de Productos (Go + MySQL + JWT)

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8-4479A1?logo=mysql&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)

API REST de productos hecha en **Go**, con **arquitectura en capas**,
base de datos **MySQL** y **autenticación con JWT**. Las rutas de lectura son
públicas; crear, editar y borrar requieren iniciar sesión.

Responde en **JSON** (por defecto), **XML** o **YAML** según el parámetro
`?format=`.

## Características

- **CRUD completo** de productos (Crear, Leer, Actualizar, Borrar).
- **Autenticación JWT**: registro, login y rutas protegidas.
- **Contraseñas encriptadas** con `bcrypt` (nunca se guardan en texto plano).
- **Middleware** propio que protege las rutas de escritura.
- **Tres formatos** de respuesta: JSON, XML y YAML.
- **SQL con parámetros** (`?`) para evitar inyección SQL.

## Requisitos

- Go 1.26+
- MySQL en marcha con una base de datos llamada `tienda`.

```sql
CREATE DATABASE tienda;
```

Las tablas (`productos`, `usuarios`) se crean solas al arrancar.

## Uso

```bash
go run .
```

La API escucha en **http://localhost:8080**.

La conexión por defecto es `root@tcp(localhost:3306)/tienda`. La clave para
firmar los tokens se puede cambiar con la variable de entorno `JWT_SECRET`.

## Endpoints

| Método | Ruta              | ¿Token? | Descripción              |
|--------|-------------------|---------|--------------------------|
| POST   | `/register`       | No      | Crear un usuario         |
| POST   | `/login`          | No      | Iniciar sesión → token   |
| GET    | `/productos`      | No      | Listar productos         |
| GET    | `/productos/{id}` | No      | Ver un producto          |
| POST   | `/productos`      | **Sí**  | Crear un producto        |
| PUT    | `/productos/{id}` | **Sí**  | Actualizar un producto   |
| DELETE | `/productos/{id}` | **Sí**  | Borrar un producto       |

## Ejemplo de uso (con `curl`)

```bash
# 1. Registrar un usuario
curl -X POST localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"samuel@correo.com","password":"miclave123"}'

# 2. Iniciar sesión y guardar el token
TOKEN=$(curl -s -X POST localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"samuel@correo.com","password":"miclave123"}' \
  | sed 's/.*"token":"//;s/".*//')

# 3. Crear un producto (con el token)
curl -X POST localhost:8080/productos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nombre":"Laptop","precio":999.99}'

# 4. Listar productos (no necesita token)
curl localhost:8080/productos

# En XML o YAML:
curl "localhost:8080/productos?format=xml"
curl "localhost:8080/productos?format=yaml"
```

Sin token, las rutas de escritura responden `401 Unauthorized`.

## Cómo funciona la autenticación

1. **Registro** (`POST /register`): la contraseña se encripta con `bcrypt` y se
   guarda el hash, nunca la contraseña real.
2. **Login** (`POST /login`): se compara la contraseña con el hash. Si coincide,
   se devuelve un **token JWT** firmado, con el `user_id` dentro y caducidad de
   24 horas.
3. **Rutas protegidas**: el middleware `auth.RequireAuth` lee la cabecera
   `Authorization: Bearer <token>`, verifica la firma y la caducidad, y solo
   entonces deja pasar.

## Estructura

```
main.go              arranque: conecta a MySQL, crea tablas y define rutas
auth/                JWT (crear/verificar tokens) y middleware RequireAuth
handlers/            handlers HTTP de productos y de autenticación
repository/          acceso a la base de datos (SQL) de productos y usuarios
models/              estructuras Product y User
response/            responder en JSON / XML / YAML
database/            conexión a MySQL
```

## Stack

Go (net/http, database/sql) · MySQL (go-sql-driver/mysql) ·
JWT (golang-jwt/jwt) · bcrypt (golang.org/x/crypto) · YAML (gopkg.in/yaml.v3).
