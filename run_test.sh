#!/bin/bash
# run-tests.sh

# Поднимаем окружение
docker-compose up -d

# Ждем, пока сервис станет доступен
while ! curl -s http://localhost:8080/health > /dev/null; do
    sleep 1
done

# Запускаем тесты
go test -v ./tests/integration

# Очищаем окружение
docker-compose down