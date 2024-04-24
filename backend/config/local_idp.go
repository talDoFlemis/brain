package config

import "github.com/taldoflemis/brain.test/internal/adapters/driven/auth"

type localIdpConfig struct {
	Ed25519Seed         string `koanf:"ed25519_seed"`
	Issuer              string `koanf:"issuer"`
	Audience            string `koanf:"audience"`
	AccessTimeInMinutes int    `koanf:"access_time_in_minutes"`
	RefreshtimeInHours  int    `koanf:"refresh_time_in_hours"`
}

func NewLocalIDPConfig() (*auth.LocalIdpConfig, error) {
	var out localIdpConfig
	err := k.Unmarshal("auth", &out)
	if err != nil {
		return nil, err
	}
	cfg := auth.NewLocalIdpConfig(
		out.Ed25519Seed,
		out.Issuer,
		out.Audience,
		out.AccessTimeInMinutes,
		out.RefreshtimeInHours,
	)
	return cfg, nil
}
