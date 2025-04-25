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

//	@title			Subscription Tracker API
//	@version		1.0
//	@description	This is an API for subscription tracker.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

// This will set accessToken for Swagger UI
// in header and name of that header field is Authorization
//
//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
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

	service := service.NewService(repo, authenticator, cfg)

	validator := validator.NewAppValidator()

	handler := handler.NewHandler(service, validator)

	router := router.NewRouter(handler, authenticator)

	mailer := mailer.NewSMTPMailer(cfg.Mailer)

	chrono := chrono.NewChrono(repo, mailer)
	go chrono.ScheduleDailyTask(8, 00)

	srv := server.NewServer(cfg.Server.Addr, router.Setup())
	srv.Run()
}
