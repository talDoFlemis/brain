package main

import (
	"log"

	"github.com/taldoflemis/brain.test/config"
	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
	seed "github.com/taldoflemis/brain.test/internal/adapters/driven/postgres/seeds"
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

	cfg, err := config.NewPostgresConfig()
	if err != nil {
		log.Fatal(err)
	}

	connStr := postgres.GenerateConnectionString(cfg)

	err = seed.Seed(connStr, "")
	if err != nil {
		log.Fatal(err)
	}
}
