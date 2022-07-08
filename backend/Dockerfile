# syntax=docker/dockerfile:1

FROM golang:1.18-alpine as builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/yummy-bot .

EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["./yummy-bot"]