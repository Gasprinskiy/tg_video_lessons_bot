FROM golang:1.24-alpine AS builder

# Устанавливаем tzdata
RUN apk add --no-cache tzdata

# Устанавливаем локальную зону
ENV TZ=Asia/Tashkent
RUN cp /usr/share/zoneinfo/Asia/Tashkent /etc/localtime \
    && echo "Asia/Tashkent" > /etc/timezone

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o main .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]