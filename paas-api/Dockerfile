# -------- Build Stage --------
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o paas-api

# -------- Runtime Stage --------
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy statically linked binary from builder
COPY --from=builder /app/paas-api .

# Expose API port
EXPOSE 8080

# Start the app
CMD ["./paas-api"]
