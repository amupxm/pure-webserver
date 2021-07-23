FROM golang:1.15.6-alpine3.12 AS builder

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .
CMD ["./pure-webserver"]
