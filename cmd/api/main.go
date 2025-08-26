package main

import (
	"database/sql"
	"log"
	_ "github.com/anshbadoni30/event-management-app/docs"
	"github.com/anshbadoni30/event-management-app/internal/database"
	"github.com/anshbadoni30/event-management-app/internal/env"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

// @title Event Management System API
// @version 1.0
// @description A RestAPI in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}
	models := database.NewModels(db)
	app := application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    models,
	}
	er := app.serve()
	if er != nil {
		log.Fatal(er)
	}

}
