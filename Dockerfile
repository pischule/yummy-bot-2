# syntax=docker/dockerfile:1

FROM golang:1.18-alpine as builder

# build the backend

RUN apk add --no-cache gcc musl-dev g++ leptonica-dev tesseract-ocr-dev opencv-dev npm

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY *.go ./

RUN go build -a -o yummy-bot-2

# build the frontend

ENV NODE_ENV=production

WORKDIR /app/frontend

COPY ["frontend/package.json", "frontend/package-lock.json*", "./"]

RUN npm install --omit=dev

COPY frontend /app/frontend

RUN npm run build


FROM alpine:3.16

RUN apk add --no-cache leptonica tesseract-ocr opencv tesseract-ocr-data-rus

WORKDIR /app

COPY --from=builder /app/yummy-bot-2 ./yummy-bot-2
COPY --from=builder /app/frontend/build ./frontend/build

COPY ./rects-tool /app/rects-tool

EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["./yummy-bot-2"]