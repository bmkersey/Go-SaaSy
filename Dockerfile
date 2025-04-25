# Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

# Install goose (if needed at runtime)
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build main app binary
RUN go build -o saasy ./cmd/saasy

# Build migration runner
RUN go build -o migrate ./cmd/migrate

EXPOSE 8080

CMD ["./saasy"]
