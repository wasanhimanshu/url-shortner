package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT     string
	DBName   string
	DBHost   string
	DBPort   string
	DBPasswd string
	DBUser   string
}

var Envs = InitConfig()

func InitConfig() Config {
	godotenv.Load()
	return Config{
		PORT:     getEnv("PORT", "8080"),
		DBName:   getEnv("DB_NAME", "urlshortner"),
		DBHost:   getEnv("DB_HOST", "127.0.0.1"),
		DBPort:   getEnv("DB_PORT", "3306"),
		DBPasswd: getEnv("DB_PASSWORD", ""),
		DBUser:   getEnv("DB_USER", "root"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
