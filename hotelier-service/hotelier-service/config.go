package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// Загружаем конфигурацию
func loadConfig() Config {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return config
}

// Структура конфигурации
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}
