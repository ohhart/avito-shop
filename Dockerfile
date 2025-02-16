# Используем официальный образ Golang
FROM golang:1.22

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . /app

# Скачиваем зависимости
RUN go mod tidy

# Компилируем бинарник сервиса
RUN go build -o /build ./cmd/main.go \
    && go clean -cache -modcache

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/build"]