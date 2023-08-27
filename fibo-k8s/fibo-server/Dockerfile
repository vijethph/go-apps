# Use the official Go image as the base image
FROM golang:1.20-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use the official Alpine image as the base image
FROM alpine:3.18

# Install the CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the application binary from the build image
COPY --from=build /app/app .

# Expose the port
EXPOSE 5000

# Start the application
CMD ["./app"]
