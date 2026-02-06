FROM golang:1.25-alpine AS build-env


WORKDIR /app


COPY go.* ./
RUN go mod download -x
COPY *.go ./
RUN go build -o alerter main.go

FROM alpine:3.17.0
LABEL maintainer="slazzaris@gmail.com"

COPY --from=build-env /app/alerter /usr/local/bin

ENTRYPOINT ["/usr/local/bin/alerter"]
