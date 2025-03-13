package main

import (
	"fmt"

	"github.com/sangtandoan/subscription_tracker/internal/config"
	"github.com/sangtandoan/subscription_tracker/internal/db"
	"github.com/sangtandoan/subscription_tracker/internal/handler"
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

	service := service.NewService(repo)

	validator := validator.NewAppValidator()

	handler := handler.NewHandler(service, validator)

	router := router.NewRouter(handler)

	srv := server.NewServer(":8080", router.Setup())
	srv.Run()
}
