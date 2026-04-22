FROM golang:1.25.0-alpine AS server-builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/apiserver ./cmd/apiserver

FROM alpine:latest

# Create application directory
RUN mkdir -p /app /app/site

WORKDIR /app

ENV ENVIRONMENT="production"

# Copying the Go binary
COPY --from=server-builder /app/bin ./bin

EXPOSE 3000

CMD ["./bin/apiserver"]
