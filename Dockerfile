# Builder Image
FROM golang:1.26 AS builder

RUN git clone --depth=1 https://github.com/CaptainKills/xtream-api.git /application
WORKDIR /application

RUN go mod download
RUN go build -v -o /xtream-api

# Executable Image
FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder --chmod=+x /xtream-api /xtream-api

CMD ["./xtream-api"]
