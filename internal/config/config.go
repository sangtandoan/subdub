package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Db            *DBConfig
	Server        *ServerConfig
	Authenticator *AuthenticatorConfig
	GoogleOAuth   *oauth2.Config
	Mailer        *MailerConfig
}

type DBConfig struct {
	ConnString      string
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

type MailerConfig struct {
	From     string
	Host     string
	Username string
	Password string
	Port     int
}

type AuthenticatorConfig struct {
	SecretKey   string
	TokenExpiry string
}

type ServerConfig struct {
	Addr string
}

func LoadConfig() (*Config, error) {
	// This function loads environment variables from a .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbConfig := &DBConfig{
		ConnString:      getEnv("DB_CONN_STRING", ""),
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

	authenticatorConfig := &AuthenticatorConfig{
		SecretKey:   getEnv("JWT_SECRET_KEY", "secret"),
		TokenExpiry: getEnv("TOKEN_EXPIRY", "24h"),
	}

	mailerConfig := &MailerConfig{
		From:     getEnv("MAIL_FROM", "sangvaminh11497@gmai.com"),
		Host:     getEnv("MAIL_HOST", "smtp.gmail.com"),
		Port:     getEnvAsInt("MAIL_PORT", 587),
		Username: getEnv("MAIL_USERNAME", "sangvaminh11497@gmail.com"),
		Password: getEnv("MAIL_PASSWORD", ""),
	}

	googleOAuthConfig := &oauth2.Config{
		ClientID:     getEnv("GOOGLE_CLIENT", ""),
		ClientSecret: getEnv("GOOGLE_SECRET", ""),
		RedirectURL:  "http://localhost:8080/api/v1/oauth2/callback",
		Scopes:       []string{"email"},
		Endpoint:     google.Endpoint,
	}

	srvConfig := &ServerConfig{
		Addr: getEnv("ADDR", ":8080"),
	}

	return &Config{
		Db:            dbConfig,
		Server:        srvConfig,
		Authenticator: authenticatorConfig,
		Mailer:        mailerConfig,
		GoogleOAuth:   googleOAuthConfig,
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
