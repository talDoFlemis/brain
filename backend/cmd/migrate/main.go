package main

import (
	"log"

	"github.com/taldoflemis/brain.test/config"
	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
	migrate "github.com/taldoflemis/brain.test/internal/adapters/driven/postgres/migrations"
)

func main() {
	koanf := config.NewKoanfson()
	err := koanf.LoadFromJSON("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = koanf.LoadFromEnv("BRAIN_")
	if err != nil {
		log.Fatal(err)
	}
	pgCfg, err := config.NewPostgresConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(pgCfg.Database)

	connStr := postgres.GenerateConnectionString(pgCfg)
	log.Print(connStr)
	migrate.Migrate(connStr, "internal/adapters/driven/postgres/migrations/")
}
