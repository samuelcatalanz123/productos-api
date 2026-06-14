# ---- Etapa 1: compilar ----
# Usamos la imagen oficial de Go (grande, con todo el compilador) solo para
# construir el programa.
FROM golang:1.26 AS builder

WORKDIR /app

# Copiamos primero los archivos de dependencias y las descargamos. Docker
# guarda esta capa en caché: si no cambian, no se vuelven a descargar.
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código y compilamos un binario estático.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api .

# ---- Etapa 2: imagen final (mínima) ----
# alpine es una imagen de Linux diminuta. La app final pesa poquísimo.
FROM alpine:latest

WORKDIR /app

# Copiamos SOLO el binario compilado desde la etapa anterior.
COPY --from=builder /api /app/api

EXPOSE 8080
CMD ["/app/api"]
