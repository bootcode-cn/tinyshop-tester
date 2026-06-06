# Stage 1: Build the Go binary
FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o tinyshop-tester \
    -ldflags="-s -w" \
    .

# Stage 2: Runtime image with Python 3
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    python3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -m -s /bin/bash tester

COPY --from=builder /app/tinyshop-tester /usr/local/bin/tinyshop-tester

WORKDIR /workspace

USER tester

ENTRYPOINT ["tinyshop-tester"]
CMD ["--help"]
