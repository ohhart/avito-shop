# avito-shop

Для запуска проекта нужно выполнить команду `docker-compose up`.
После этого сервис будет доступен на порту `:8080`

## Стек
Проект написан на языке Golang + Postgres.

## Функционал
Были реализованы все методы из предоставленной API.

## Тестирование

!!!Написаны интеграционные тесты для всех сценариев. Находятся в папке tests/integration.
Интеграционные тесты успешно проверяют основную функциональность приложения, включая:
- Аутентификацию
- Валидацию входных данных
- Работу с БД
- Обработку ошибок
- Защиту endpoints
- на сценарий покупки мерча
- сценарий передачи монеток другим сотрудникам
- и тд
Также проведены тесты для негативных сценариев (отрицательный баланс, пустой username и тд.).

!!!Также проведены юнит тесты для основной бизнес логики, находятся в папке internal/services.
Последний раз покрытие было 31,7%. Там основная недоработка в базе данных, нужно было написать обходной способ для работы с бд. Не успела доработать, к сожалению. 

## Для запуска тестов нужно выполнить следующие команды:

docker-compose up -d
затем: go test -v ./tests/integration - для интеграционных тестов 
и go test ./internal/services -v
(или флаг -cover чтобы проверить процент покрытия).

## Линтер 
Также была описана конфигурацию линтера (.golangci.yaml в корне проекта). Использовала линтер golangci-lint.
