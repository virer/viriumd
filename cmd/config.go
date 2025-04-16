package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfigFromFile(filename string) (Config, error) {
	var cfg Config

	file, err := os.Open(filename)
	if err != nil {
		return cfg, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("decode config: %w", err)
	}

	return cfg, nil
}
