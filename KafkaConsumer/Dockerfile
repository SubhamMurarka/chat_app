# Dockerfile for products service
FROM golang:1.22.0-alpine

# Set the working directory inside the container.
WORKDIR /app

# Install curl.
RUN apk add --no-cache curl

# Copy go.mod and go.sum files to the working directory.
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the application.
RUN go build -o main .

# Run the executable.
CMD ["./main"]
