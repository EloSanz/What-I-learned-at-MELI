# Etapa de compilación
FROM golang:alpine AS builder
WORKDIR /app

# Copiamos el descriptor de dependencias
COPY go.mod ./

# Descargamos dependencias (si aplica)
RUN go mod download

# Copiamos el código
COPY . .

# Compilamos el binario estático de Go
RUN CGO_ENABLED=0 GOOS=linux go build -o items-api cmd/api/main.go

# Etapa final súper ligera
FROM alpine:3.19
WORKDIR /app

# Copiamos el binario compilado
COPY --from=builder /app/items-api .

# Exponemos el puerto de escucha por defecto del microservicio
EXPOSE 8081

# Ejecutamos el microservicio
CMD ["./items-api"]
