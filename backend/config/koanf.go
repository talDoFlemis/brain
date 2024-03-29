package config

import (
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Koanfson struct {
}

func NewKoanfson() *Koanfson {
	return &Koanfson{}
}

func (kson *Koanfson) LoadFromJSON(path string) error {
	parser := json.Parser()
	if err := k.Load(file.Provider("config.json"), parser); err != nil {
		return err
	}
	return nil
}
