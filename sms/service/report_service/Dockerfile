# Multi-stage build for Go application
FROM golang:1.24.1-alpine AS builder

# Set necessary environmet variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -a -installsuffix cgo -o main cmd/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /dist/main .

# Expose port 8802
EXPOSE 8802

# Command to run
CMD ["./main"]