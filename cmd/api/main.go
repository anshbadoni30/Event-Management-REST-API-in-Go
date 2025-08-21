package main

import (
	"database/sql"
	"log"

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
