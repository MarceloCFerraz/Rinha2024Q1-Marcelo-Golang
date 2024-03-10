# Use the official .NET SDK image as the base image
FROM golang:1.22.1-alpine3.19 AS builder

# Set working directory for the build context
WORKDIR /app

# Copy the dependency files first
COPY src/ ./

# Install dependencies
RUN go mod download && go mod verify

# Build the application (replace with your own build command)
# Assuming your main entry point is in cmd/main.go
RUN go build -o Rinha24MarceloGo .

# Use a smaller image for running the application
FROM alpine

# Copy the compiled binary from the builder stage
COPY --from=builder /app/Rinha24MarceloGo /app/Rinha24MarceloGo

# Set the working directory for the application
WORKDIR /app

# Set the entrypoint to 'som'
ENTRYPOINT ["./Rinha24MarceloGo"]

# Expose the port your application listens on (if applicable)
# EXPOSE 8080  # Example: Expose port 8080