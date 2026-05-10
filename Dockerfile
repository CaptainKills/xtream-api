# Builder Image
FROM golang:1.26 AS builder

RUN git clone --depth=1 https://github.com/CaptainKills/xtream-api.git /application
WORKDIR /application

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-X main.Branch=$(git branch --show-current)" -v -o /xtream-api

# Executable Image
FROM alpine:latest

RUN apk add --no-cache ca-certificates libc6-compat

COPY --from=builder --chmod=+x /xtream-api /xtream-api

CMD ["./xtream-api"]
