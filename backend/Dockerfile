# Build stage
FROM golang:1.24.4 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

# Copy binary
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]