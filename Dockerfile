# Step 1: Build the Go binary
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/generate/main.go

# Step 2: Run the Go binary in a small container
FROM alpine:latest

COPY --from=builder /app/main /app/main

WORKDIR /root/

# EXPOSE 8080

CMD ["/app/main"]
