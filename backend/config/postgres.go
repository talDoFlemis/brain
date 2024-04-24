package config

import "github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"

type postgresConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
}

func NewPostgresConfig() (*postgres.Config, error) {
	var out postgresConfig
	err := k.Unmarshal("postgres", &out)
	if err != nil {
		return nil, err
	}
	return postgres.NewConfig(out.User, out.Password, out.Host, out.Database, out.Port), nil
}
