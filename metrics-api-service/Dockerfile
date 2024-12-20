# Use the official Golang image
FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o metrics-api-service

# Expose the port
EXPOSE 8001

# Run the application
CMD ["./metrics-api-service"]
