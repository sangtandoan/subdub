package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sangtandoan/subscription_tracker/internal/config"
	"github.com/sangtandoan/subscription_tracker/internal/models"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/enums"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
)

var names []string = []string{
	"Spotify Premium",
	"Netflix Premium",
	"YouTube Premium",
	"Apple One",
	"Notion Plus",
	"Figma Professional",
	"GitHub Pro",
	"Adobe Creative Cloud",
	"Zoom Pro",
	"Microsoft 365 Personal",
}

var db *sql.DB

func main() {
	connectDB()
	defer db.Close()

	for range 1000 {
		FillSubscriptions()
	}
}

func connectDB() {
	cfg := &config.DBConfig{
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

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, _ = sql.Open("postgres", connStr)

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}

func FillSubscriptions() {
	startDate := randomStartDate()
	duration := randomDuration()
	endDate := calculateEndDate(startDate, duration)
	name := randomName()
	isCancelled := randomBool()

	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err)
		return
	}

	query := `
		INSERT INTO 
		subscriptions (id, user_id, name, start_date, end_date, duration, is_cancelled) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, name, start_date, end_date, duration, is_cancelled
	`
	userID := "676ead13-20d9-11f0-95a9-902e1685779a"
	userIDUUID, _ := uuid.Parse(userID)

	row := db.QueryRowContext(
		context.Background(),
		query,
		id,
		userIDUUID,
		name,
		time.Time(startDate),
		endDate,
		duration.String(),
		isCancelled,
	)

	var subcription repo.SubscriptionRow
	err = row.Scan(
		&subcription.ID,
		&subcription.UserID,
		&subcription.Name,
		&subcription.StartDate,
		&subcription.EndDate,
		&subcription.Duration,
		&subcription.IsCancelled,
	)
	if err != nil {
		fmt.Printf("Error inserting subscription: %v\n", err)
		return
	}
}

func randomStartDate() models.SubscriptionTime {
	startDate := time.Now().AddDate(0, 0, -rand.Intn(30))
	return models.SubscriptionTime(time.Date(
		startDate.Year(),
		startDate.Month(),
		startDate.Day(),
		0,
		0,
		0,
		0,
		startDate.Location(),
	))
}

func randomDuration() enums.Duration {
	d, _ := enums.ParseString2Duration(enums.AllDurations[rand.Intn(len(enums.AllDurations))])
	return d
}

func randomName() string {
	return names[rand.Intn(len(names))]
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

func calculateEndDate(
	startDate models.SubscriptionTime,
	duration enums.Duration,
) time.Time {
	return duration.AddDurationToTime(time.Time(startDate))
}

func randomBool() bool {
	return rand.Intn(2) == 0
}
