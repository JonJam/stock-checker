# Builder
FROM golang:1.15.6-buster AS builder
# Install Delve
RUN go get github.com/go-delve/delve/cmd/dlv
# Build src
COPY . /src
WORKDIR /src
RUN go build -o app -gcflags="all=-N -l"

# Base runner
FROM debian:10.7 AS base-runner
# Install chromium
RUN apt-get update && apt-get install -y \
    ca-certificates \
    chromium \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/app /
COPY config.yml /

# Dev runner
FROM base-runner AS dev-runner
EXPOSE 40000
COPY --from=builder /go/bin/dlv /
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app"]

# Prod runner
FROM base-runner AS prod-runner
CMD ["/app"]