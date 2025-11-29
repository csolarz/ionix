# ---------------------------------------
# 1. Build Stage (Go 1.25.1)
# ---------------------------------------
FROM golang:1.25.1 AS builder

# Set necessary variables for a static build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy go.mod + go.sum first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -ldflags="-s -w" -o /app/app .

# ---------------------------------------
# 2. Final Stage: Distroless
# ---------------------------------------
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /app/app /app

USER nonroot:nonroot

ENTRYPOINT ["/app"]
