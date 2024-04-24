package main

import (
	"log"

	"github.com/taldoflemis/brain.test/config"
	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
)

func main() {
	koanf := config.NewKoanfson()
	err := koanf.LoadFromJSON("config.json")
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
	postgres.Migrate(connStr, "internal/adapters/driven/postgres/migrations/")
}
