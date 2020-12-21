# Builder
FROM golang:1.15.6-buster AS builder

COPY . /src

WORKDIR /src

RUN go build -o app

# App
FROM debian:10.7

RUN apt-get update && apt-get install -y \
    ca-certificates \
    chromium \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/app /

COPY config.yml /

CMD ["/app"]