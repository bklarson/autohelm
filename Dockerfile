# --- Stage 1: Build Go binary ---
FROM golang:1.22-alpine AS builder

ENV CGO_ENABLED=0
WORKDIR /app

COPY main.go .

RUN go build -o autohelm main.go

# --- Stage 2: Get Docker CLI and Compose plugin ---
FROM docker:25-cli AS docker-cli

# --- Stage 3: Final image ---
FROM gcr.io/distroless/static:latest

COPY --from=builder /app/autohelm /autohelm

COPY --from=docker-cli /usr/local/bin/docker /usr/local/bin/docker
COPY --from=docker-cli /usr/local/libexec/docker/cli-plugins/docker-compose /usr/local/libexec/docker/cli-plugins/docker-compose

ENTRYPOINT ["/autohelm"]
