package main

import (
	"fmt"
	"log"
	"os"

	// Import pg driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	flag "github.com/spf13/pflag"

	// Import migrations.
	_ "github.com/zahar517/brute-force-attack-prevention/migrations"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("migrate: bad args")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	command := args[0]
	err = goose.Run(command, db, "./migrations", args[1:]...)
	if err != nil {
		log.Printf("goose %v: %v", command, err)
	}
}
