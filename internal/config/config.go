package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Db            *DBConfig
	Server        *ServerConfig
	Authenticator *AuthenticatorConfig
}

type DBConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	MaxConnLifeTime string
	MaxIdleLifeTime string
	MaxOpenConns    int
	MaxIdleConns    int
}

type AuthenticatorConfig struct {
	SecretKey   string
	TokenExpiry string
}

type ServerConfig struct {
	Addr string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbConfig := &DBConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		Username:        getEnv("DB_USERNAME", "admin"),
		Password:        getEnv("DB_PASSWORD", "secret"),
		DBName:          getEnv("DB_NAME", "subscription"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 20),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 20),
		MaxConnLifeTime: getEnv("DB_MAX_CONN_LIFE_TIME", "30m"),
		MaxIdleLifeTime: getEnv("DB_MAX_IDLE_LIFE_TIME", "10m"),
	}

	srvConfig := &ServerConfig{
		Addr: getEnv("ADDR", ":8080"),
	}

	authenticatorConfig := &AuthenticatorConfig{
		SecretKey:   getEnv("JWT_SECRET_KEY", "secret"),
		TokenExpiry: getEnv("TOKEN_EXPIRY", "30m"),
	}

	return &Config{
		Db:            dbConfig,
		Server:        srvConfig,
		Authenticator: authenticatorConfig,
	}, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return valueAsInt
}

func getEnvAsBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valueAsBool, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return valueAsBool
}
