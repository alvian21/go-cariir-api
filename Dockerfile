# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum (if exists)
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ---- Run Stage ----
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy other necessary files (like templates, static assets, etc. if any)
# COPY --from=builder /app/public ./public

# Expose the correct port
EXPOSE 7310

# Run the application
CMD ["./main"]