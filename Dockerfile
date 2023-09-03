FROM golang:1.21.0-alpine as builder
WORKDIR /build
COPY . .
RUN go build -o Weibo-To-Telegram main.go

FROM alpine:3.15.4
WORKDIR /app
COPY --from=builder /build/Weibo-To-Telegram .
ENTRYPOINT [ "./Weibo-To-Telegram" ]