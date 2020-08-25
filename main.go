package main

import (
	"go-jwt/config"
	"go-jwt/database/mysql"
	"go-jwt/internal"
	"go-jwt/logging"
	"go-jwt/logging/slack"
	"go-jwt/routes"
	"log"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// main function to boot up everything
func main() {

	// Set Configs
	cfg := config.Get()

	notifySlack := slack.Get(cfg.Messaging, cfg.General)

	dlog := logging.Get(cfg.General, notifySlack)

	// Connect to Database
	db, err := mysql.Get(cfg.GetDBConnStr(), dlog)
	if err != nil {
		log.Fatalf("Error Connection Database :%s", err.Error())
	}

	ctlr := internal.Get(db, *dlog)

	// Creates a http server
	routes.InitRoutes(ctlr, dlog)
}
