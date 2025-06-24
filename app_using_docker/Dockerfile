FROM golang:1.24.1

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the binary
COPY . .

# Set the entrypoint to the compiled binary
CMD ["go", "build", "app_using_docker"]