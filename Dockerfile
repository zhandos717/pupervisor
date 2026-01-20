# Используем официальный образ Go для сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Второй этап - создание минимального образа для запуска
FROM alpine:latest

# Устанавливаем зависимости для работы приложения
RUN apk --no-cache add ca-certificates

# Создаем пользователя для запуска приложения
RUN addgroup -g 65532 nonroot && adduser -D -u 6532 -G nonroot nonroot

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из builder-контейнера
COPY --from=builder /app/main .

# Копируем веб-файлы
COPY web/ ./web/

# Устанавливаем права доступа
RUN chown -R nonroot:nonroot /app
USER nonroot

# Открываем порт
EXPOSE 8080

# Команда запуска
CMD ["./main"]