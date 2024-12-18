# Use the official Golang image
FROM golang:1.21

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN go mod init rule-engine-api
RUN go mod tidy

# Build the Go application
RUN go build -o rule-engine-api

# Expose the port
EXPOSE 8003

# Run the application
CMD ["./rule-engine-api"]
