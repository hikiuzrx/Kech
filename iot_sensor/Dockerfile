FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# We use cmd/device as entry point
RUN CGO_ENABLED=0 GOOS=linux go build -o iot-sensor ./cmd/device

# Final stage
FROM alpine:latest

WORKDIR /app

# Install run dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/iot-sensor .

# Env vars will be supplied by docker-compose
ENV SIMULATION_MODE=true

# Run the binary
CMD ["./iot-sensor"]
