package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/google/uuid"
)

func createSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	var req SnapshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidInput(req.VolumeID) {
		http.Error(w, "invalid initiator name format", http.StatusBadRequest)
		return
	}

	snapshotID := uuid.New().String()
	volumeName := "virium-vol-" + req.VolumeID
	snapshotName := "virium-snap-" + snapshotID
	volumeGroup := config.VGName

	log.Printf("Creating snapshot ref: %s for vol: %s in volumeGroup %s", snapshotID, volumeName, volumeGroup)

	// LVM: Create snapshot
	lvCreateCmd := exec.Command("sudo", "lvcreate", "-s", "--size", "8M", "-n", snapshotName, fmt.Sprintf("/dev/%s/%s", volumeGroup, volumeName))
	out, err := lvCreateCmd.CombinedOutput()
	if err != nil {
		log.Printf("lvcreate error: %v\n%s", err, out)
		http.Error(w, "LVM create snapshot failed", http.StatusInternalServerError)
		return
	}
	log.Println("LVM snapshot created:", req.VolumeID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SnapshotRequest{VolumeID: snapshotID})
}

func deleteSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteSnapshotRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidInput(req.VolumeID) {
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}
	snapshotName := "virium-snap-" + req.VolumeID
	log.Printf("Removing snapshot ID: %s in volumeGroup %s", snapshotName, config.VGName)

	// LVM: Remove logical volume
	lvRemoveCmd := exec.Command("sudo", "lvremove", "-y", fmt.Sprintf("%s/%s", config.VGName, snapshotName))
	out, err := lvRemoveCmd.CombinedOutput()
	if err != nil {
		log.Printf("lvremove error: %v\n%s", err, out)
		http.Error(w, "LVM delete failed", http.StatusInternalServerError)
		return
	}

	log.Println("LVM snapshot deleted:", req.VolumeID)
	w.WriteHeader(http.StatusNoContent)
}
