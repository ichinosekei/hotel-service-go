package main

import (
	//"gopkg.in/yaml.v3"
	//"io/ioutil"
	"log"
    "os"
	//"github.com/joho/godotenv"
)
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("Using environment variable %s", key)
		return value
	}
	log.Printf("Environment variable %s not set", key)
	return defaultValue
}

// Загружаем конфигурацию
func loadConfig() Config {
	return Config{
		Server: ServerConfig{
			InternalPort:    getEnv("SERVICE_INTERNAL_PORT", "8080"),
			ExternalPort:    getEnv("SERVICE_EXTERNAL_PORT", "8082"),
			InternalGrpcPort: getEnv("INTERNAL_SERVER_GRPC_PORT", "50051"),
			ExternalGrpcPort: getEnv("EXTERNAL_SERVER_GRPC_PORT", "50051"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "hotelier"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

// Структура конфигурации
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	InternalPort    string
	ExternalPort    string
	InternalGrpcPort string
	ExternalGrpcPort string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
