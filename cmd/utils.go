package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

func executeWrapper(commandline string, errmsg string) ([]byte, error) {
	cmdPrefix := config.Cmd_prefix
	parts := strings.Fields(commandline)
	cmdinstance := exec.Command(cmdPrefix, parts[0:]...)
	out, err := cmdinstance.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("%s %s %s", errmsg, err, out)
	}
	return out, err
}
