# Builder
FROM golang:1.15.6-buster AS builder
COPY . /src
WORKDIR /src
# Install Delve
RUN go get github.com/go-delve/delve/cmd/dlv
# Build without optimisations. See https://github.com/golang/vscode-go/blob/master/docs/debugging.md#try-building-your-binary-without-compiler-optimizations
RUN go build -gcflags="all=-N -l" -o app-dev
RUN go build -o app

# Base runner
FROM debian:10.7 AS base-runner
# Install chromium
RUN apt-get update && apt-get install -y \
    ca-certificates \
    chromium \
    && rm -rf /var/lib/apt/lists/*


# Dev runner
FROM base-runner AS dev-runner
WORKDIR /server
COPY config.yml /server
COPY --from=builder /src/app-dev /server/app
COPY --from=builder /go/bin/dlv /server
EXPOSE 40000
CMD ["/server/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server/app"]


# Prod runner
FROM base-runner AS prod-runner
WORKDIR /server
COPY config.yml /server
COPY --from=builder /src/app /server
CMD ["./app"]