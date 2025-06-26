# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server cmd/server/main.go
RUN go build -o share cmd/cli/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binaries from builder
COPY --from=builder /app/server .
COPY --from=builder /app/share /usr/local/bin/

# Copy web assets if needed
# COPY --from=builder /app/web/dist ./web/dist

EXPOSE 8080

CMD ["./server"]