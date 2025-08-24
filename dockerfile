# Use the official Go image as a base
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final container
FROM gcr.io/distroless/base

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/main

# Command to run the executable
CMD ["/app/main"]
# Expose the port the app runs on
EXPOSE 8080