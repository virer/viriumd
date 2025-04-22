package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	klog "k8s.io/klog/v2"
)

func createSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	var req SnapshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidInput(req.VolumeID) {
		klog.V(5).Info("invalid source volume id", req)
		http.Error(w, "invalid source volume id format", http.StatusBadRequest)
		return
	}

	volumeName := "virium-vol-" + req.VolumeID
	snapshotName := "virium-snap-" + req.Name
	volumeGroup := config.VGName

	klog.V(2).Infof("Creating snapshot ref: %s for vol: %s in volumeGroup %s", snapshotName, volumeName, volumeGroup)
	klog.V(5).Infof("sudo lvcreate -s --size 8M -n %s %s", snapshotName, fmt.Sprintf("/dev/%s/%s", volumeGroup, volumeName))

	// LVM: Create snapshot
	lvCreateCmd := exec.Command("sudo", "lvcreate", "-s", "--size", "8M", "-n", snapshotName, fmt.Sprintf("/dev/%s/%s", volumeGroup, volumeName))
	out, err := lvCreateCmd.CombinedOutput()
	if err != nil {
		klog.Error("lvcreate error:", err, out)
		http.Error(w, "LVM create snapshot failed", http.StatusInternalServerError)
		return
	}
	klog.V(2).Info("LVM snapshot created:", req.VolumeID)

	capacity, _ := GetVolumeSize(fmt.Sprintf("/dev/%s/%s", volumeGroup, volumeName))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		SnapshotResponse{
			SnapshotId:     req.Name,
			SourceVolumeID: volumeName,
			Capacity:       capacity,
		})
}

func deleteSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteSnapshotRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidInput(req.VolumeID) {
		w.WriteHeader(http.StatusNoContent)
		// FIXME http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}
	snapshotName := "virium-snap-" + req.VolumeID
	klog.V(2).Infof("Removing snapshot ID: %s in volumeGroup %s", snapshotName, config.VGName)

	// LVM: Remove logical volume
	lvRemoveCmd := exec.Command("sudo", "lvremove", "-y", fmt.Sprintf("%s/%s", config.VGName, snapshotName))
	out, err := lvRemoveCmd.CombinedOutput()
	if err != nil {
		klog.Error("lvremove error", err, out)
		http.Error(w, "LVM delete failed", http.StatusInternalServerError)
		return
	}

	klog.V(2).Info("LVM snapshot deleted:", req.VolumeID)
	w.WriteHeader(http.StatusNoContent)
}
