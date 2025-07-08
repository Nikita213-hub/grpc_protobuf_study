# ---- Builder Stage ----
    FROM golang:1.23-alpine AS builder

    WORKDIR /workspace
    
    # Copy go.work and all go.mod/sum files first for better caching
    COPY go.work.auth ./go.work
    COPY go.mod go.sum ./
    COPY auth-service/go.mod auth-service/go.mod
    
    # Download dependencies for the workspace
    RUN go mod download
    
    # Copy the rest of the source code
    COPY . .
    
    # Build the auth-service binary
    RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service-binary ./auth-service/cmd/main.go
    
    # ---- Final Stage ----
    FROM alpine:latest
    
    WORKDIR /app
    
    # Copy config file (if needed by your app)
    COPY auth-service/cmd/local.yaml .
    
    # Copy the built binary from the builder stage
    COPY --from=builder /auth-service-binary .
    
    EXPOSE 44044
    
    CMD ["./auth-service-binary"]
