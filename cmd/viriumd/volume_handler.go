package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	klog "k8s.io/klog/v2"
)

func createVolumeHandler(w http.ResponseWriter, r *http.Request) {
	var req VolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidInput(req.InitiatorName) {
		http.Error(w, "invalid initiator name format", http.StatusBadRequest)
		return
	}

	lvmsize := req.Capacity / (1024 * 1024)
	volumeID := uuid.New().String()
	volumeName := "virium-vol-" + volumeID

	klog.V(1).Infof("Creating %d MiB volumeID: %s in volumeGroup %s", lvmsize, volumeID, config.VGName)

	// LVM: Create logical volume
	lvCreateCmd := exec.Command("sudo", "lvcreate", "-T", "-L", fmt.Sprintf("%dM", lvmsize), "-n", volumeName, config.VGName)
	out, err := lvCreateCmd.CombinedOutput()
	if err != nil {
		klog.V(2).Infof("lvcreate error: %v\n%s", err, out)
		http.Error(w, "LVM create failed", http.StatusInternalServerError)
		return
	}
	klog.V(2).Info("LVM volume created:", volumeName)

	iqn, err := createISCSITarget(volumeID, volumeName, req.InitiatorName)
	if err != nil {
		klog.Error("iSCSI error:", err)
		http.Error(w, "iSCSI error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(VolumeResponse{
		VolumeID:          volumeID,
		TargetPortal:      config.TargetPortal,
		Iqn:               string(iqn),
		Lun:               "0",
		DiscoveryCHAPAuth: "true",
		SessionCHAPAuth:   "false",
	})
}

func deleteVolumeHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteVolumeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidInput(req.VolumeID) {
		http.Error(w, "invalid ID name format", http.StatusBadRequest)
		return
	}
	volumeName := "virium-vol-" + req.VolumeID
	klog.V(2).Infof("Removing volumeID: %s in volumeGroup %s", volumeName, config.VGName)

	// Remove iSCSI export first
	err := deleteISCSITarget(req.VolumeID, volumeName)
	if err != nil {
		klog.Error("iSCSI error:", err)
		http.Error(w, "iSCSI error", http.StatusInternalServerError)
		return
	}

	// LVM: Remove logical volume
	lvRemoveCmd := exec.Command("sudo", "lvremove", "-y", fmt.Sprintf("%s/%s", config.VGName, volumeName))
	out, err := lvRemoveCmd.CombinedOutput()
	if err != nil {
		if strings.HasPrefix(string(out), "  Failed to find logical volume") {
			klog.V(2).Info("logical volume already removed!")
		} else {
			klog.Error("lvremove error", err, out)
			http.Error(w, "LVM delete failed", http.StatusInternalServerError)
			return
		}
	} else {
		klog.V(2).Info("LVM volume deleted:", req.VolumeID)
	}

	w.WriteHeader(http.StatusNoContent)
}

func resizeVolumeHandler(w http.ResponseWriter, r *http.Request) {
	var req VolumeResizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidInput(req.VolumeID) {
		http.Error(w, "invalid volume id format", http.StatusBadRequest)
		return
	}

	lvmsize := req.Capacity / (1024 * 1024)
	volumeID := req.VolumeID
	volumeName := "virium-vol-" + volumeID

	klog.V(1).Infof("Extending %d MiB volumeID: %s in volumeGroup %s", lvmsize, volumeName, config.VGName)

	// LVM: Create logical volume
	lvCreateCmd := exec.Command("sudo", "lvextend", "-L", fmt.Sprintf("%dM", lvmsize), "-n", fmt.Sprintf("%s/%s", config.VGName, volumeName))
	out, err := lvCreateCmd.CombinedOutput()
	if err != nil {
		klog.Error("lvextend error:", err, out)
		http.Error(w, "LVM extend failed", http.StatusInternalServerError)
		return
	}
	klog.V(2).Info("LVM volume extended:", volumeName)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(VolumeResponse{VolumeID: volumeID})
}
