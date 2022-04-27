package main

import (
	"log"

	// Import pg driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	flag "github.com/spf13/pflag"

	// Import migrations.
	"github.com/zahar517/brute-force-attack-prevention/internal/config"
	_ "github.com/zahar517/brute-force-attack-prevention/migrations"
)

var configFile string

func init() {
	flag.StringVarP(&configFile, "config", "c", "", "Path to configuration file")
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("migrate: bad args")
	}

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := goose.OpenDBWithDriver("pgx", config.Database.Dsn)
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
