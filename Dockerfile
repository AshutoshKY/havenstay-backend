FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server/main.go

# Create a minimal production image
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app/

# Copy the built binary
COPY --from=builder /app/bin/server /app/server

# Expose HTTP and gRPC ports
EXPOSE 8080
EXPOSE 50051

CMD ["/app/server"]
