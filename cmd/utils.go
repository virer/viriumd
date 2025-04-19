package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

func LoadConfigFromFile(filename string) error {
	once.Do(func() {
		// Load default config first
		config = NewConfiguration()

		// Open the YaML config file
		file, err := os.Open(filename)
		if err != nil {
			klog.Error("open config file: %w", err)
			return
		}
		defer file.Close()

		// Unmarshal the file content into the config struct
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			klog.Error("decode config: %w", err)
			return
		}
	})
	return nil
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
