# ==============================
# Build Stage
# ==============================
FROM golang:1.24 AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o stock-scraper ./cmd/api/main.go

# Create non-root user
RUN useradd -r -u 1001 appuser

# ==============================
# Final Stage
# ==============================
FROM debian:bullseye-slim
WORKDIR /app

# Install CA certificates (for HTTPS)
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy binary and user info
COPY --from=build-stage /app/stock-scraper .
COPY --from=build-stage /etc/passwd /etc/group /etc/

# Run as non-root
USER appuser

CMD ["./stock-scraper"]
