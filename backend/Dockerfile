# syntax=docker/dockerfile:1

FROM golang:1.18-alpine as builder

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/yummy-bot .

EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["./yummy-bot"]