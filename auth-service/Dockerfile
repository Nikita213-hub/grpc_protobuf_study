
    FROM golang:1.23-alpine AS builder

    WORKDIR /workspace
    

    COPY go.work.auth ./go.work
    COPY go.mod go.sum ./
    COPY auth-service/go.mod auth-service/go.mod
    

    RUN go mod download
    

    COPY . .
    

    RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service-binary ./auth-service/cmd/main.go
    

    FROM alpine:latest
    
    WORKDIR /app
    

    COPY auth-service/cmd/local.yaml .

    COPY --from=builder /auth-service-binary .
    
    EXPOSE 44044
    
    CMD ["./auth-service-binary"]
