package config

import "github.com/taldoflemis/brain.test/internal/adapters/drivers/web"

type fiberConfig struct {
	Prefix           string `koanf:"prefix"`
	ListenIP         string `koanf:"listen_ip"`
	Port             int    `koanf:"port"`
	CORSAllowOrigins string `koanf:"cors.allow_origins"`
	CORSAllowHeaders string `koanf:"cors.allow_headers"`
	CORSAllowMethods string `koanf:"cors.allow_methods"`
}

func NewFiberConfig() (*web.Config, error) {
	var out fiberConfig
	err := k.Unmarshal("fiber", &out)
	if err != nil {
		return nil, err
	}
	return &web.Config{
		Prefix:           k.String("fiber.prefix"),
		ListenIP:         k.String("fiber.listen_ip"),
		Port:             k.Int("fiber.port"),
		CORSAllowOrigins: k.String("fiber.cors.allow_origins"),
		CORSAllowHeaders: k.String("fiber.cors.allow_headers"),
		CORSAllowMethods: k.String("fiber.cors.allow_methods"),
	}, nil
}
