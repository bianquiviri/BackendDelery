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

# Install Swagger tool and generate documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go
RUN go mod tidy

# Build the binary.
# Use TARGETARCH to automatically build for the correct processor (ARM or x86).
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /go/bin/api cmd/api/main.go

# Use a clean, minimal Alpine image for the final stage.
FROM alpine:3.19

# Import the user and group files from the builder.
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

# Copy the binary and static files to the production image from the builder stage.
COPY --from=builder /go/bin/api /app/api
COPY public /app/public
# Also copy the example env file. The actual env should be provided at runtime.
COPY .env.example /app/.env

# Run the web service on container startup.
EXPOSE 8084
ENTRYPOINT ["/app/api"]
