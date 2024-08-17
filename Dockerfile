# Start from the latest golang base image
FROM golang:1.23.0-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY config/config.dev.yaml ./
COPY docs/swagger.json ./
COPY i18n/tr.json ./
COPY i18n/en.json ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/ticket_service/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o service cmd/ticket_service/main.go

# Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# This lib is required
RUN apk add libc6-compat
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/service .
COPY --from=builder /app/config/config.dev.yaml ./config/config.dev.yaml
COPY --from=builder /app/docs/swagger.json ./doc/swagger.json
COPY --from=builder /app/i18n/tr.json ./i18n/tr.json
COPY --from=builder /app/i18n/en.json ./i18n/en.json

ENV APP_ENV=dev

EXPOSE 8080:8080

# Command to run the executable
ENTRYPOINT ["./service"]
