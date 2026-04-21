FROM golang:1.25.0-alpine AS server-builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/worker ./cmd/worker

FROM alpine:latest

# Install required packages for user management and su-exec
RUN apk add --no-cache ca-certificates shadow su-exec

# Copying the entrypoint script from previous stage
COPY --from=server-builder /app/script/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Create application directory
RUN mkdir -p /app /app/site

WORKDIR /app

ENV ENVIRONMENT="production"

# Copying the Go binary
COPY --from=server-builder /app/bin ./bin

# DO NOT set USER here - we need to run as root initially to create users and fix permissions
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["./bin/worker"]
