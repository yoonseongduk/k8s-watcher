# Build stage
FROM golang:1.23-alpine AS builder
# Add git for go mod download
RUN apk add --no-cache git
# Set working directory
WORKDIR /app
# Copy go mod files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy source code
COPY main.go .
# Ensure go.sum is up to date
RUN go mod tidy
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o watcher .

# Final stage
FROM alpine:3.19

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Seoul
WORKDIR /root/
# Copy the binary from builder
COPY --from=builder /app/watcher .
# Copy kubeconfig if needed
# COPY kubeconfig /root/.kube/config
# Command to run
CMD ["./watcher"]