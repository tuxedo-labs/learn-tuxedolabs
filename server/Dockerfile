ARG GO_VERSION=1.23.5
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# First copy module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 8000
CMD ["./main"]
