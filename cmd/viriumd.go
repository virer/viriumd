package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var config Config
var version string = "v0.1.3.4"
var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9.:-]+$`)

func isValidInput(s string) bool {
	return validNamePattern.MatchString(s)
}

func main() {
	var err error
	configPath := flag.String("config", "/etc/virium/virium.yaml", "Path to configuration file")
	flag.Parse()

	config, err = LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// fmt.Printf("Config loaded: %+v\n", config)

	if config.VGName == "" {
		config.VGName = "vg_data"
	}
	if config.Port == "" {
		config.Port = "8787"
	}
	if config.Base_iqn == "" {
		config.Base_iqn = "iqn.2025-04.net.virer.virium"
	}
	if config.TargetPortal == "" {
		config.Base_iqn = "iqn.2025-04.net.virer.virium"
	}
	log.Printf("Starting virium on port %s using vol:%s (%s)", config.Port, config.VGName, version)

	http.HandleFunc("/api/volumes/create", createVolumeHandler)
	http.HandleFunc("/api/volumes/delete", deleteVolumeHandler)
	http.HandleFunc("/api/volumes/resize", resizeVolumeHandler)

	http.HandleFunc("/api/snapshoft/create", createSnapshotHandler)
	http.HandleFunc("/api/snapshoft/delete", deleteSnapshotHandler)

	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
}
