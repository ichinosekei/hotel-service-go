# hotel-service-go
 бронирования отеля на языке go

# Сруктура проекта
```bash
final-project/
├── api/                    # Папка для API контрактов
│   ├── protobufs/          # Контракты API в формате Protocol Buffers (gRPC или REST)
│   └── openapi/            # Контракты API в формате OpenAPI/Swagger
│
├── cmd/                    # Папка для запуска сервисов
│   ├── booking-service/    # Основной сервис бронирования
│   │   └── main.go         # Точка входа для сервиса бронирования
│   ├── hotelier-service/   # Сервис для отельеров
│   │   └── main.go         # Точка входа для сервиса отельеров
│   └── client-service/     # Клиентский сервис
│       └── main.go         # Точка входа для клиентского сервиса
│
├── configs/                # Конфигурации для всех сервисов
│   ├── booking-config.yaml # Конфигурация сервиса бронирования
│   ├── hotelier-config.yaml# Конфигурация сервиса отельеров
│   └── client-config.yaml  # Конфигурация клиентского сервиса
│
├── deployments/            # Файлы для развертывания сервисов
│   └── docker-compose.yml  # Docker Compose файл для всех сервисов и зависимостей
│
├── internal/               # Пакеты, которые не предназначены для использования вне проекта
│   ├── booking/            # Пакеты сервиса бронирования
│   ├── hotelier/           # Пакеты сервиса отельеров
│   └── client/             # Пакеты клиентского сервиса
│       └── repository/     # Логика работы с базой данных для клиентского сервиса
│
├── pkg/                    # Общие пакеты, которые можно переиспользовать между сервисами
│   ├── logger/             # Пакет для структурированного логирования
│   └── tracing/            # Пакет для трассировки
│
├── migrations/             # SQL-файлы для миграции базы данных
│   ├── 001_create_tables.sql
│   └── 002_add_indexes.sql
│
├── scripts/                # Скрипты для автоматизации работы с проектом
│   ├── init_db.sh          # Скрипт инициализации базы данных
│   └── run_tests.sh        # Скрипт для запуска всех тестов
│
├── test/                   # Папка для интеграционных тестов и моков
│   ├── booking_tests/      # Интеграционные тесты сервиса бронирования
│   └── mocks/              # Моки для тестирования
│
├── .env                    # Файл переменных окружения для настройки сервисов
├── .gitignore              # Список файлов и папок, игнорируемых git
└── README.md               # Основная документация по проекту








```

### План реализации hotelier service
#### 1. Контракты и API (api/openapi/hotelier.yaml):
Контракт сервиса содержит эндпоинты для управления отелями и номерамb:
{ 
* POST /hotels: Создать отель
* GET /hotels: Получить список отелей
* POST /rooms: Создать номер
* GET /rooms: Получить номера отеля
}
#### 2. Миграции (migrations/):
   SQL-файлы для создания таблиц Hotels и Rooms, а также их отношений.
#### 3. Основная логика (internal/hotelier/):
   Основные функции, такие как добавление отеля и номера, получение доступных номеров.
#### 4. Тесты (test/hotelier_tests/):
   Тесты для проверки эндпоинтов и бизнес-логики.
#### Конфигурация (configs/hotelier-config.yaml):
 Настройки подключения к базе данных, логирования и переменные окружения.

### Подробнее о реализации:
* API Контракты (api/openapi/hotelier.yaml). Создадим спецификацию API для отелей и номеров.
* Миграции (migrations/001_create_tables.sql). SQL для создания таблиц Hotels и Rooms.
* Основная логика (internal/hotelier/) Создадим основные функции сервиса для обработки отелей и номеров. internal/hotelier/service.go:
* Тесты (test/hotelier_tests/service_test.go)
* Docker и Конфигурация
  1. docker-compose.yml Файл docker-compose.yml задает конфигурацию для базы данных PostgreSQL и сервиса hotelier.
  2. Конфигурация Go сервиса (configs/hotelier-config.yaml) Конфигурационный файл для сервиса hotelier в формате YAML. В нем указываются параметры подключения к базе данных и порты.
  3. Dockerfile для hotelier service Dockerfile для сборки и запуска сервиса hotelier.
* Скрипты миграции:
  Для автоматического применения миграций, добавить инициализацию базы данных в Dockerfile .
```dockerfile
COPY migrations /migrations
RUN ./migrate -path /migrations -database "postgres://hotelier:hotelier_password@db:5432/hotelier_db?sslmode=disable" up
```
* Запуск
  Для запуска hotelier service и базы данных в Docker:
``` bash
docker-compose up --build
```
Этот набор конфигураций поднимет базу данных PostgreSQL, настроит соединение с hotelier service и откроет API на порту 8080.