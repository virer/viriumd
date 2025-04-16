package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

type VolumeRequest struct {
	Name     string `json:"name"`
	Capacity int64  `json:"capacity"` // bytes
}

type VolumeResponse struct {
	ID string `json:"id"`
}

type DeleteVolumeRequest struct {
	VolumeID string `json:"volume id"`
}

type Config struct {
	VGName string `json:"vg_name"`
}

var config Config

func createVolumeHandler(w http.ResponseWriter, r *http.Request) {
	var req VolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// LVM: Create logical volume
	lvmsize := req.Capacity / (1024 * 1024)
	// volumeName := req.Name
	volumeName := "virium-vol-" + uuid.New().String()
	volumeGroup := config.VGName

	log.Printf("Creating %d MiB volumeID: %s in volumeGroup %s", lvmsize, volumeName, volumeGroup)

	lvCreateCmd := exec.Command("sudo", "lvcreate", "-L", fmt.Sprintf("%dM", lvmsize), "-n", volumeName, volumeGroup)
	out, err := lvCreateCmd.CombinedOutput()
	if err != nil {
		log.Printf("lvcreate error: %v\n%s", err, out)
		http.Error(w, "LVM create failed", http.StatusInternalServerError)
		return
	}
	log.Println("LVM volume created:", volumeName)

	createISCSITarget(volumeName)

	json.NewEncoder(w).Encode(VolumeResponse{ID: volumeName})
}

func deleteVolumeHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteVolumeRequest
	volumeID := req.VolumeID
	volumeGroup := config.VGName

	log.Printf("Removing volumeID: %s in volumeGroup %s", volumeID, volumeGroup)

	// Remove iSCSI export first
	deleteISCSITarget(volumeID)

	// LVM: Remove logical volume
	lvRemoveCmd := exec.Command("sudo", "lvremove", "-y", fmt.Sprintf("%s/%s", volumeGroup, volumeID))
	out, err := lvRemoveCmd.CombinedOutput()
	if err != nil {
		log.Printf("lvremove error: %v\n%s", err, out)
		http.Error(w, "LVM delete failed", http.StatusInternalServerError)
		return
	}

	log.Println("LVM volume deleted:", volumeID)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	config.VGName = os.Getenv("VG_DATA")
	if config.VGName == "" {
		config.VGName = "vg_data"
	}

	http.HandleFunc("/api/volumes/create", createVolumeHandler)
	http.HandleFunc("/api/volumes/delete", deleteVolumeHandler)

	http.ListenAndServe(":8787", nil)
}
