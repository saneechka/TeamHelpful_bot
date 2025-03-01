# Используем многоэтапную сборку для минимизации размера образа
FROM golang:latest

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
FROM alpine:latest

# Установка необходимых пакетов для SQLite
RUN apk add --no-cache ca-certificates tzdata sqlite

# Создание пользователя без прав root
RUN adduser -D -g '' appuser

# Установка рабочей директории
WORKDIR /app

# Копирование бинарного файла из предыдущего этапа
COPY --from=builder /app/helpbot .

# Копирование файла базы данных, если он существует
COPY --from=builder /app/users.db ./users.db


USER appuser


CMD ["./helpbot"] 