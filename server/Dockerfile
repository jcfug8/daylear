# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o ./bin/daylearctl ./daylearctl

# Stage 2: Run the Go binary
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Install imagemagick
RUN apk add --no-cache imagemagick
RUN apk add --no-cache imagemagick-heic
RUN apk add --no-cache imagemagick-jpeg
RUN apk add --no-cache imagemagick-jxl
RUN apk add --no-cache imagemagick-pdf
RUN apk add --no-cache imagemagick-raw
RUN apk add --no-cache imagemagick-svg
RUN apk add --no-cache imagemagick-tiff
RUN apk add --no-cache imagemagick-webp

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/daylearctl .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./daylearctl"]