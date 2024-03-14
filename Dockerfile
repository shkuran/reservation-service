# Use the official Golang image as the base image
FROM golang:1.22 AS build-stage

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o reservation-service .

# Create a minimal runtime image
FROM alpine:3.19

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the builder image
COPY --from=build-stage /app/reservation-service /app/reservation-service
COPY /config.yaml /app/config.yaml

# Expose the port on which the Go application will run
EXPOSE 8082

# Define the command to run the Go application
CMD ["./reservation-service"]