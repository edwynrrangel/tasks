FROM golang:1.20.3-alpine3.17 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main main.go

# Start fresh from a smaller image
FROM alpine:3.17.1 AS production

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/main /app/bin/main

# Add a non-root user
RUN adduser -D user

# Change ownership of the copied files to the non-root user
RUN chown -R user:user /app

# Change to non-root user
USER user

# Set environment variables
ENV PORT 3000

# Expose port
EXPOSE ${PORT}

# Set the Current Working Directory inside the container
WORKDIR /app

# Command to run the executable
CMD ["bin/main"]
