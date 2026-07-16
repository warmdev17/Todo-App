package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sqlx.DB

func initDB() {
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}

	var err error
	DB, err = sqlx.Open("pgx", dsn)

	if err != nil {
		log.Fatal("Failed to open DB connection: ", err)
	}

	err = DB.Ping()

	if err != nil {
		log.Fatal("Failed to ping DB")
	}

	log.Println("PostgreSQL (via pgx) connected successfully!")
}
