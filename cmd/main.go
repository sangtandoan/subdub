package main

import (
	"fmt"

	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/chrono"
	"github.com/sangtandoan/subscription_tracker/internal/config"
	"github.com/sangtandoan/subscription_tracker/internal/db"
	"github.com/sangtandoan/subscription_tracker/internal/handler"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/mailer"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/validator"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
	"github.com/sangtandoan/subscription_tracker/internal/router"
	"github.com/sangtandoan/subscription_tracker/internal/server"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := db.NewDB(cfg.Db)
	if err != nil {
		panic(err)
	}

	fmt.Println("")
	fmt.Println("database connected")
	fmt.Println("")

	repo := repo.NewRepo(db)

	authenticator, err := authenticator.NewJWTAuthenticator(cfg.Authenticator)
	if err != nil {
		panic(err)
	}

	service := service.NewService(repo, authenticator)

	validator := validator.NewAppValidator()

	handler := handler.NewHandler(service, validator)

	router := router.NewRouter(handler, authenticator)

	mailer := mailer.NewSMTPMailer(cfg.Mailer)

	chrono := chrono.NewChrono(repo, mailer)
	go chrono.ScheduleDailyTask(11, 29)

	srv := server.NewServer(cfg.Server.Addr, router.Setup())
	srv.Run()
}
