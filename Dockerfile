# Используем официальный образ Go в качестве базового образа
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые пакеты для работы с C
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект в контейнер
COPY . .

# Собираем приложение с включенным CGO
ENV CGO_ENABLED=1
RUN go build -o main ./cmd/webserver/main.go

# Используем легковесный образ для запуска приложения
FROM alpine:latest

# Копируем скомпилированный бинарник из сборочного образа
COPY --from=builder /app/main /app/main

# Копируем веб-ресурсы
COPY --from=builder /app/web /web

# Копируем базу данных в контейнер
COPY --from=builder /app/storage/scheduler.db /storage/scheduler.db

# Указываем порт, который будет использоваться приложением
EXPOSE 7540

# Устанавливаем переменные окружения
ENV TODO_DBFILE="/storage/scheduler.db"
ENV TODO_PORT="7540"

# Запускаем приложение
CMD ["/app/main"]
