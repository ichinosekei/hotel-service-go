# Hotel Service GO

## О проекте

Это бэкенд-сервис для управления отелем, написанный на языке Go. Сервис предоставляет API для управления бронированиями, номерами, пользователями и отелями.

## Технологии

* **Go:** Основной язык программирования.
* **Echo:** Высокопроизводительный, расширяемый и минималистичный веб-фреймворк для Go.
* **MongoDB:** NoSQL база данных для хранения информации.
* **JWT (JSON Web Tokens):** Для аутентификации и авторизации пользователей.
* **Bcrypt:** Для хеширования паролей.
* **Docker и Docker Compose:** Для контейнеризации и управления сервисами.

## Структура проекта

* **api/**
    * `admin.go`
    * `api.go`
    * `auth.go`
    * `booking.go`
    * `hotel.go`
    * `room.go`
    * `user.go`
* **cmd/**
    * **api/**
        * `main.go`
* **db/**
    * `db.go`
    * `fixtures.go`
    * `hotel_store.go`
    * `room_store.go`
    * `user_store.go`
* **scripts/**
    * `seed.go`
* **types/**
    * `hotel.go`
    * `user.go`
* `docker-compose.yml`
* `Dockerfile`
* `go.mod`
* `go.sum`
* `README.md`
* `api`: Обработчики HTTP-запросов.
* `cmd/api`: Точка входа в приложение.
* `db`: Взаимодействие с базой данных.
* `scripts`: Скрипты для наполнения базы данных.
* `types`: Определения типов данных.

## Установка и запуск

1.  Клонируйте репозиторий:

    ```bash
    git clone [https://github.com/ichinosekei/hotel-service-go.git](https://github.com/ichinosekei/hotel-service-go.git)
    ```

2.  Перейдите в директорию проекта:

    ```bash
    cd hotel-service-go
    ```

3.  Создайте файл `.env` и укажите в нем следующие переменные:

    ```env
    MONGO_DB_URL=mongodb://localhost:27017
    MONGO_DB_NAME=hotel-reservation
    JWT_SECRET=your_jwt_secret
    ```

4.  Запустите проект с помощью Docker Compose:

    ```bash
    docker-compose up
    ```

5.  Для наполнения базы данных начальными данными, выполните следующую команду в отдельном терминале:

    ```bash
    go run scripts/seed.go
    ```

## API Эндпоинты

Подробное описание всех эндпоинтов можно найти в исходном коде в папке `api`.

**Примеры:**

* **POST /api/user:** Создание нового пользователя.
* **GET /api/user/{id}:** Получение информации о пользователе по ID.
* **POST /api/auth/login:** Аутентификация пользователя.
* **GET /api/hotel:** Получение списка отелей.
* **GET /api/hotel/{id}/rooms:** Получение списка номеров в отеле.

## TODO
* Добавить способ бронирования через телеграмм
* Реализовать более гибкое управление правами доступа.
* Добавить логирование.
* Улучшить валидацию данных.

