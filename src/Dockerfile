FROM golang:1.24-alpine AS builder

# # Установка tzdata
# RUN apt-get update && apt-get install -y tzdata

# # Установка переменной окружения для таймзоны (опционально)
# RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
nan
WORKDIR /app
COPY . .

RUN go mod tidy && go build -o main .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]