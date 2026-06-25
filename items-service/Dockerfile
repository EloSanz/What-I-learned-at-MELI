# Build stage
FROM golang:alpine AS builder
WORKDIR /app

# Disable go.work to avoid missing module errors
ENV GOWORK=off

# Copy the go.mod files
COPY items-service/go.mod items-service/go.sum* ./items-service/
COPY meli-sdk/go.mod meli-sdk/go.sum* ./meli-sdk/

# Download dependencies from inside items-service
WORKDIR /app/items-service
RUN go mod download

# Copy the rest of the source code
COPY items-service/ ./
COPY meli-sdk/ ../meli-sdk/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o items-api cmd/api/main.go

# Super lightweight final stage
FROM alpine:3.19
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/items-service/items-api .

# Expose the default listening port
EXPOSE 8081

# Execute the microservice
CMD ["./items-api"]
