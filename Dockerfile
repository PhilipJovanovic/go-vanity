# syntax=docker/dockerfile:1

FROM golang:1.24.0 AS build

ENV DEBIAN_FRONTEND noninteractive

# Set destination for COPY
WORKDIR /app

COPY go.mod ./
RUN go mod download

RUN apt-get update

COPY . .

RUN go build -o vanity .

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 1337

# Run
CMD ["./vanity"]