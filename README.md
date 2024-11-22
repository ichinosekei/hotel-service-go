# hotel-service-go
 бронирования отеля на языке go

# Сруктура проекта
```bash



hotel-service-go/
├── booking-service/            # Сервис для управления бронированием, не реализова
├── client-service/             # Сервис для управления данными клиентов, не реализован
├── deployments/                # Конфигурации для деплоя
│   └── docker-compose.yml      # Файл для настройки Docker Compose, для hotel-service
├── hotel-service/              # Основной сервис управления отелями
│   ├── api/                    # API спецификации
│   │   ├── openapi/            # OpenAPI спецификации
│   │   └── protobufs/          # Протобуф-контракты, нет
│   ├── configs/                # Конфигурационные файлы сервиса
│   │   └── hotelier-config.yaml # Основной конфигурационный файл
│   ├── hotelier-service/       # Код основного сервиса
│   │   ├── Dockerfile          # Dockerfile для сборки сервиса
│   │   ├── main.go             # Точка входа сервиса
│   │   └── service.go          # обработчик запросов
│   ├── hotelier_tests/         # Тесты для сервиса
│   │   └── service_test.go     # Тесты для модуля hotelier-service
│   ├── internal/               # Внутренние модули и библиотеки
│   │   └── hotelier/           # Реализация бизнес-логики сервиса
│   │       └── service.go      # Бизнес-логика отеля
│   ├── tracing/                # Настройки трассировки 
│   └── metrics/                # Реализация метрик
│       └── metrics.go          # потом сделаем
├── migrations/                 # SQL-скрипты для управления базой данных
│   ├── 001_create_tables.sql   # Скрипт для создания таблиц
│   └── 002_add_indexes.sql     # Скрипт для добавления индексов не реализован
├── test/                       # Тестирование
│   ├── booking_tests/          # Тесты для сервиса бронирования не реализовано
│   └── mocks/                  # Моки для тестов не реализовано
├── .env                        # Переменные окружения
├── go.mod                      # Модульный файл Go
└── README.md                   # Описание проекта



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
docker-compose up --build
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
