# Используем многоэтапную сборку для минимизации размера образа
FROM golang:1.21-alpine AS builder

# Установка необходимых зависимостей
RUN apk add --no-cache git gcc musl-dev

# Установка рабочей директории
WORKDIR /app

# Копирование файлов go.mod и go.sum
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o helpbot .

# Финальный образ
FROM alpine:3.16

# Установка рабочей директории
WORKDIR /app

# Копирование бинарного файла из предыдущего этапа
COPY --from=builder /app/bot .

CMD ["./bot"] 