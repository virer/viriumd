package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9.:-]+$`)

func isValidInput(s string) bool {
	return validNamePattern.MatchString(s)
}

func LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(config)
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
