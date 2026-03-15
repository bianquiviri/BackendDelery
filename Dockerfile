# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.25-alpine as builder

# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.mod go.sum ./
RUN go mod download

# Copy local code to the container image.
COPY . .

# Run tests (Resilience & correctness check within the container context)
RUN go test -v ./...

# Build the binary.
# -ldflags="-w -s" removes debugging info, reducing binary size (Performance).
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/api cmd/api/main.go

# Use a clean, minimal Alpine image for the final stage.
FROM alpine:3.19

# Import the user and group files from the builder.
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/api /app/api
# Also copy the example env file. The actual env should be provided at runtime.
COPY .env.example /app/.env

# Run the web service on container startup.
EXPOSE 8084
ENTRYPOINT ["/app/api"]
