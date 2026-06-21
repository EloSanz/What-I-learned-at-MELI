# Etapa de compilación
FROM golang:alpine AS builder
WORKDIR /app

# Copiamos el descriptor de dependencias
COPY go.mod ./

# Descargamos dependencias
RUN go mod download

# Copiamos el código
COPY . .

# Compilamos el binario de Go
RUN CGO_ENABLED=0 GOOS=linux go build -o orders-api cmd/api/main.go

# Etapa final súper ligera
FROM alpine:3.19
WORKDIR /app

# Copiamos el binario compilado
COPY --from=builder /app/orders-api .

# Exponemos el puerto de escucha
EXPOSE 8082

# Ejecutamos el microservicio
CMD ["./orders-api"]
