FROM golang:1.23.0 AS builder

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot main.go
RUN ls -l /opt/app  # Проверяем наличие файла bot

FROM alpine:latest

WORKDIR /opt/app

COPY --from=builder /opt/app/bot /opt/app/bot
RUN apk add --no-cache libc6-compat
RUN ls -l /opt/app  # Проверяем наличие и права файла bot
RUN chmod +x /opt/app/bot  # Устанавливаем права на выполнение

ENV REDIS_ADDR="redis:6379"
ENV REDIS_PASSWORD=""

ENTRYPOINT ["/opt/app/bot"]
CMD ["-tg-bot-token", "7483223949:AAE3QuIxCO7gXy1cG-qJzRbLocXF53WS3SQ"]
