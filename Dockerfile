# ── Build stage ────────────────────────────────────
FROM golang:1.25-alpine AS builder

# CGO is required for go-sqlite3
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate templ files, then build
RUN templ generate
RUN CGO_ENABLED=1 go build -o /app/server .

# ── Runtime stage ──────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary and static assets
COPY --from=builder /app/server .
COPY --from=builder /app/static ./static

# SQLite data lives here
RUN mkdir -p /app/data

EXPOSE 8080

ENTRYPOINT ["./server"]