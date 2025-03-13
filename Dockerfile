# Step 1: Use the official Golang 1.23.5 image to build the app
FROM golang:1.23.5 AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the entire application into the container
COPY . .

# Step 5: Build the Go application
RUN go build -o myapp .

# Step 6: Use a slim base image (Debian) for the runtime environment
FROM debian:bookworm-slim

# Step 7: Install necessary libraries (including ca-certificates and libc6)
RUN apt-get update && apt-get install -y ca-certificates libc6

# Step 8: Set the working directory inside the container
WORKDIR /root/

# Step 9: Copy the Go binary from the builder stage
COPY --from=builder /app/myapp .

# Step 10: Copy the .env file into the container
COPY .env ./

# Step 11: Expose the port your application runs on (if applicable)
EXPOSE 8080

# Step 12: Define the entry point for the container to run the app
CMD ["./myapp"]