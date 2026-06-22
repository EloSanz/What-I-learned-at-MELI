# Build stage
FROM golang:alpine AS builder
WORKDIR /app

# Copy the dependency descriptor
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Compile the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o orders-api cmd/api/main.go

# Super lightweight final stage
FROM alpine:3.19
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/orders-api .

# Expose the listening port
EXPOSE 8082

# Execute the microservice
CMD ["./orders-api"]
