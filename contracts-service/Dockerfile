FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.work.contracts ./go.work
COPY go.mod go.sum ./
COPY contracts-service/go.mod contracts-service/go.mod

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /contracts-service-binary ./contracts-service/cmd/main.go

FROM alpine:latest
    
WORKDIR /app

COPY --from=builder /contracts-service-binary .

EXPOSE 44045

CMD ["./contracts-service-binary"]
