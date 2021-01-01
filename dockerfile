FROM golang:1.15.2-alpine AS build

COPY . /app

WORKDIR /app

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./bin/demo

ENTRYPOINT ["./bin/demo"]