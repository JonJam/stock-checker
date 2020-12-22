# Builder
FROM golang:1.15.6-buster AS base-builder
ADD . /src
WORKDIR /src


# Dev Builder
FROM base-builder AS dev-builder
# Install Delve
RUN go get github.com/go-delve/delve/cmd/dlv
# Build without optimisations. See https://github.com/golang/vscode-go/blob/master/docs/debugging.md#try-building-your-binary-without-compiler-optimizations
RUN go build -gcflags="all=-N -l" -o app


# Prod Builder
FROM base-builder AS prod-builder
RUN go build -o app


# Base runner
FROM debian:10.7 AS base-runner
# Install chromium
RUN apt-get update && apt-get install -y \
    ca-certificates \
    chromium \
    && rm -rf /var/lib/apt/lists/*


# Dev runner
FROM dev-builder AS dev-runner
# COPY --from=dev-builder /src /
COPY --from=dev-builder /go/bin/dlv /src
EXPOSE 40000
WORKDIR /src
CMD ["/src/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/src/app"]


# TODO Fix this
# Prod runner
FROM base-runner AS prod-runner
COPY --from=prod-builder /src/app /server
COPY config.yml /server
WORKDIR /server
CMD ["./app"]