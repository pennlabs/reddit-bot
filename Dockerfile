FROM golang:latest as build

LABEL maintainer="Peyton Walters <pwpon500@gmail.com>"

WORKDIR /app

# Copy over source
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build
RUN go build -o labs-bot .

FROM debian:buster-slim

COPY --from=build /app/labs-bot /labs-bot

RUN apt-get update \
    && apt-get install --no-install-recommends -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

CMD ["/labs-bot", "-config", "/config.yaml"]
