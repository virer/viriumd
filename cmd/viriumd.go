package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"sync"

	klog "k8s.io/klog/v2"
)

var config *Config
var version string = "v0.2.0"

var once sync.Once // Ensure the config is loaded only once
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
	config := NewConfiguration() // Loading default validation
	err = LoadConfigFromFile(*configPath)
	if err != nil {
		klog.Fatalf("Failed to load config: %v", err)
	}
	klog.V(5).Infof("Config loaded: %+v\n", config)

	// Create a new mux router
	mux := http.NewServeMux()

	mux.Handle("/api/volumes/create", basicAuthMiddleware(http.HandlerFunc(createVolumeHandler)))
	mux.Handle("/api/volumes/delete", basicAuthMiddleware(http.HandlerFunc(deleteVolumeHandler)))
	mux.Handle("/api/volumes/resize", basicAuthMiddleware(http.HandlerFunc(resizeVolumeHandler)))

	mux.Handle("/api/snapshoft/create", basicAuthMiddleware(http.HandlerFunc(createSnapshotHandler)))
	mux.Handle("/api/snapshoft/delete", basicAuthMiddleware(http.HandlerFunc(deleteSnapshotHandler)))

	klog.V(1).Infof("Starting virium on port %s using vol:%s (%s)", config.Port, config.VGName, version)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), mux); err != nil {
		klog.Fatal("Error starting server:", err)
	}
}
