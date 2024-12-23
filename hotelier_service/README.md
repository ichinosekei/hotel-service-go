# hotel-service-go
 бронирования отеля на языке go

# Сруктура проекта
```bash

hotel-service-go/
├── booking-service/            # Сервис для управления бронированием (ещё не реализован)
├── client-service/             # Сервис для управления клиентами (ещё не реализован)
├── deployments/                # Конфигурации для развёртывания сервисов
│   └── docker-compose.yml      # Конфигурация Docker Compose для запуска сервисов
├── hotelier_service/         # Основной сервис управления отелями
│   ├── api/                    # API спецификации
│   │   ├── openapi/            # OpenAPI спецификации (документирование API)
│   │   └── protobufs/          # Протобафс-контракты для GRPC (ещё не реализованы)
│   ├── configs/                # Конфигурационные файлы для hotelier-сервиса
│   │   └── hotelier-config.yaml # Основной YAML-файл конфигурации сервиса
│   ├── hotelier_service/     # Код основного сервиса
│   │   ├── Dockerfile          # Файл сборки Docker-образа для сервиса
│   │   ├── main.go             # Точка входа в приложение (основной файл сервера)
│   │   └── config.go           # Логика для обработки конфигурационных файлов
│   ├── internal/               # Внутренние компоненты сервиса
│   │   ├── pkg/                # Пакеты для реализации логики сервиса
│   │   │   ├── api/            # Логика обработки HTTP-запросов
│   │   │   │   └── server.go   # Реализация API эндпоинтов
│   │   │   ├── metrics/        # Реализация метрик
│   │   │   │   └── metrics.go  # Код для сбора и регистрации метрик Prometheus не реализован
│   │   │   ├── repository/     # Бизнес-логика и работа с базой данных
│   │   │   │   ├── repository.go # Реализация операций с данными 
│   │   │   │   └── repository_test.go # тесты для логики репозитория
│   │   │   └── tracing/        # Реализация распределённой трассировки
│   │   │       └── tracing.go  # Код для интеграции с библиотекой OpenTelemetry
├── migrations/                 # SQL-скрипты для управления базой данных
│   ├── 001_create_tables.sql    # Скрипт для создания таблиц
│   └── 002_add_indexes.sql     # Скрипт для добавления индексов не реализован
├── test/                       # Модульные и интеграционные тесты
│   └── mocks/                   # Моки для тестов не реализовано
└── .env.dev.dev                        # Переменные окружения для сервиса


```

### План реализации hotelier service
#### 1. Контракты и API:
Контракт сервиса содержит эндпоинты для управления отелями и номерамb:
{
в hotelier.yaml подробнее
}
#### 2. Миграции (migrations/):
   SQL-файлы для создания таблиц Hotels и Rooms, а также их отношений.
#### 3. Основная логика (hotel-service/internal/hotelier/):
   Основные функции, такие как добавление отеля и номера, получение доступных номеров.
#### 4. Тесты (hotel-service/test/hotelier_tests/):
   Тесты для проверки эндпоинтов и бизнес-логики.
#### Конфигурация (configs/hotelier-config.yaml):
 Настройки подключения к базе данных, логирования и переменные окружения.


* Запуск
  Для запуска hotelier service и базы данных в Docker:
``` bash
docker-compose --env-file ../.env --env-file ../hotelier_service/deployments/.env up --build

```
## Немного о запросах 

Получение всех комнат:
```bash
curl -X GET http://localhost:8080/rooms
```
Получение комнат для конкретного отеля:
```bash
curl -X GET "http://localhost:8080/rooms?hotel_id=1"
```

Если вы пытаетесь создать комнату для отеля, который не существует:
Запрос:
```bash
curl -X POST -H "Content-Type: application/json" \
-d '{"hotel_id": 99, "room_number": "101", "price": 150.00}' \
http://localhost:8080/rooms
```
Ответ (JSON):
```json
{
"error": "hotel with the given ID does not exist"
}
```
## подробнее в hotel-service/api/openapi
рг