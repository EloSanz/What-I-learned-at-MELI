# Build stage
FROM golang:alpine AS builder
WORKDIR /app

# Disable go.work to avoid missing module errors
ENV GOWORK=off

# Copy the go.mod files
COPY orders-service/go.mod orders-service/go.sum* ./orders-service/
COPY meli-sdk/go.mod meli-sdk/go.sum* ./meli-sdk/

# Download dependencies from inside orders-service
WORKDIR /app/orders-service
RUN go mod download

# Copy the rest of the source code
COPY orders-service/ ./
COPY meli-sdk/ ../meli-sdk/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o orders-api cmd/api/main.go

# Super lightweight final stage
FROM alpine:3.19
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/orders-service/orders-api .

# Expose the listening port
EXPOSE 8082

# Execute the microservice
CMD ["./orders-api"]
