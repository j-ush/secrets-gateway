# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Compile with security flags: disable cgo, strip symbols, and use PIE
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -buildmode=pie -o gateway ./cmd/gateway/main.go

# Stage 2: Hardened Runtime
FROM scratch
WORKDIR /
# Copy the binary from the builder
COPY --from=builder /app/gateway /gateway
# Copy CA certificates for secure TLS communication with Vault/AWS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run as a non-privileged user (assigned at runtime via K8s securityContext)
USER 65532:65532

ENTRYPOINT ["/gateway"]

