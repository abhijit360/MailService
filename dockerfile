# Dockerfile

# Step 1: Build the Go application
FROM golang:1.20 AS builder
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Step 2: Run the application
FROM gcr.io/distroless/base-debian10
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
