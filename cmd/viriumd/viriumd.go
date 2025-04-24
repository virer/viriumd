package main

import (
	"flag"
	"net/http"

	klog "k8s.io/klog/v2"
)

var config *Config = NewConfiguration()
var version string = "v0.2.7"

// var once sync.Once // Ensure the config is loaded only once

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Set("v", "1")
	configPath := flag.String("config", "/etc/viriumd/virium.yaml", "Path to configuration file")
	flag.Parse()

	LoadFromFile(*configPath)
	klog.V(5).Infof("Config loaded: %+v\n", config)

	// Create a new mux router
	mux := http.NewServeMux()

	mux.Handle("/api/volumes/create", basicAuthMiddleware(http.HandlerFunc(createVolumeHandler)))
	mux.Handle("/api/volumes/delete", basicAuthMiddleware(http.HandlerFunc(deleteVolumeHandler)))
	mux.Handle("/api/volumes/resize", basicAuthMiddleware(http.HandlerFunc(resizeVolumeHandler)))

	mux.Handle("/api/snapshot/create", basicAuthMiddleware(http.HandlerFunc(createSnapshotHandler)))
	mux.Handle("/api/snapshot/delete", basicAuthMiddleware(http.HandlerFunc(deleteSnapshotHandler)))

	addr := ":" + config.Port
	if config.TLSEnabled {
		klog.V(1).Infof("Starting Viriumd on port %s using SSL with vol:%s (%s)", config.Port, config.VGName, version)
		err := http.ListenAndServeTLS(addr, config.TLSCertFile, config.TLSKeyFile, mux)
		if err != nil {
			klog.Fatalf("HTTPS server failed: %v", err)
		}
	} else {
		klog.V(1).Infof("Starting Viriumd on port %s using vol:%s (%s)", config.Port, config.VGName, version)
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			klog.Fatalf("HTTP server failed: %v", err)
		}
	}
}
