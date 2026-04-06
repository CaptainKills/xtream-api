FROM ubuntu:latest

COPY --chmod=+x ["xtream-api", "/xtream-api"]

CMD ["./xtream-api"]
