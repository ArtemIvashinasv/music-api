package main

import (
	"log"
	"os"

	"github.com/artemivashinasv/music-api/pkg/handler"
	"github.com/artemivashinasv/music-api/pkg/repository"
	"github.com/artemivashinasv/music-api/pkg/service"
	"github.com/artemivashinasv/music-api/server"
	"github.com/joho/godotenv"
)

// @title Music API
// @version 1.0
// @description API для управления музыкой (песни, группы, даты релизов и тексты).
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("ERROR: .env файл не найден")
	}
	log.Println("INFO: Переменные окружения загружены")

	// Подключение к базе данных
	db, err := repository.NewPostgresDB(repository.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	})
	if err != nil {
		log.Fatalf("ERROR: Ошибка подключения к базе данных: %s", err.Error())
	}
	log.Println("INFO: Успешное подключение к базе данных")

	// Выполнение миграций
	if err := repository.Migrate(db); err != nil {
		log.Fatalf("ERROR: Ошибка миграции: %s", err.Error())
	}
	log.Println("INFO: Миграции успешно выполнены")

	// Инициализация зависимостей
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запуск сервера
	srv := new(server.Server)
	log.Printf("INFO: Сервер запущен на порту %s", port)
	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		log.Fatalf("ERROR: Ошибка запуска сервера: %s", err.Error())
	}
}
