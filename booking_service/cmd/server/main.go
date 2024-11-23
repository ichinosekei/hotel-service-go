package main

import (
	"log"

	"github.com/ichinosekei/hotel-service-go/tree/dfsavffc/internal/app"
	"github.com/ichinosekei/hotel-service-go/tree/dfsavffc/internal/db"
	"github.com/labstack/echo/v4"
)

func main() {
	// Настройка строки подключения к базе данных
	dsn := "host=localhost user=booking_user password=booking_password dbname=booking_db port=5432 sslmode=disable"

	// Инициализация базы данных
	db.InitDatabase(dsn)

	// Создание экземпляра Echo
	e := echo.New()

	// Инициализация серверной реализации
	serverImpl := &app.ServerImplementation{}

	// Регистрация маршрутов из сгенерированного файла
	app.RegisterHandlers(e, serverImpl)

	// Логирование и запуск сервера
	log.Println("Server is running on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
