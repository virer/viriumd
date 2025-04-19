package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"

	klog "k8s.io/klog/v2"
)

var config Config
var version string = "v0.2.0"
var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9.:-]+$`)

func isValidInput(s string) bool {
	return validNamePattern.MatchString(s)
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Set("v", "1")
	configPath := flag.String("config", "/etc/virium/virium.yaml", "Path to configuration file")
	flag.Parse()

	var err error
	config, err = LoadConfigFromFile(*configPath)
	if err != nil {
		klog.Fatalf("Failed to load config: %v", err)
	}
	klog.V(2).Infof("Config loaded: %+v\n", config)

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
		config.TargetPortal = "127.0.0.1:3260"
	}
	if config.Cmd_prefix == "" {
		config.Cmd_prefix = "sudo"
	}
	klog.V(1).Infof("Starting virium on port %s using vol:%s (%s)", config.Port, config.VGName, version)

	http.HandleFunc("/api/volumes/create", createVolumeHandler)
	http.HandleFunc("/api/volumes/delete", deleteVolumeHandler)
	http.HandleFunc("/api/volumes/resize", resizeVolumeHandler)

	http.HandleFunc("/api/snapshoft/create", createSnapshotHandler)
	http.HandleFunc("/api/snapshoft/delete", deleteSnapshotHandler)

	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
}
