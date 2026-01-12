# --- STAGE 1: The Builder ---
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container's virtual filesystem
WORKDIR /app

# Copy the dependency files first. 
COPY go.mod go.sum ./
RUN go mod download

# Copy your actual source code into the container
COPY . .

# Compile the code. 
# CGO_ENABLED=0 creates a "static" binary (contains all libraries inside it).
# GOOS=linux ensures it's built for Linux (even if you were building on Windows).
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go


# --- STAGE 2: The Runner ---
FROM alpine:latest

# Set the working directory for the final app
WORKDIR /root/

# Copy ONLY the 'main' binary we just built in the first stage.
COPY --from=builder /app/main .


EXPOSE 8080

# This command runs when the container starts.
CMD ["./main"]