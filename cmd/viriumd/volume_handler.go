package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
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

	volumeSource := ""
	lvmsize := req.Capacity / (1024 * 1024)
	if req.ContentSource != nil {
		klog.V(5).Infof("ContentSource %+v", req.ContentSource.Type)
		if req.ContentSource.Type.Snapshot != nil {
			volumeSource = "virium-snap-" + req.ContentSource.Type.Snapshot.SnapshotID
		} else if req.ContentSource.Type.Volume != nil {
			volumeSource = "virium-vol-" + req.ContentSource.Type.Volume.VolumeID
		}
		tmplvmsize, err := GetVolumeSize(fmt.Sprintf("/dev/%s/%s", config.VGName, volumeSource))
		if err != nil || tmplvmsize <= 0 {
			klog.V(2).Infof("ContentSource lvm size error: %d: %s", tmplvmsize, err)
			http.Error(w, "LVM ContentSource create failed", http.StatusInternalServerError)
		}
	}

	volumeID := uuid.New().String()
	volumeName := "virium-vol-" + volumeID

	klog.V(1).Infof("Creating %d MiB volumeID: %s in volumeGroup %s", lvmsize, volumeID, config.VGName)

	// LVM: Create logical volume
	lvCreateCmd := exec.Command("sudo", "lvcreate", "-y", "-L", fmt.Sprintf("%dM", lvmsize), "-n", volumeName, config.VGName)
	out, err := lvCreateCmd.CombinedOutput()
	if err != nil {
		klog.V(2).Infof("lvcreate error: %v\n%s", err, out)
		http.Error(w, "LVM create failed", http.StatusInternalServerError)
		return
	}
	klog.V(2).Info("LVM volume created:", volumeName)

	if req.ContentSource != nil {
		// LVM: Create logical volume from specified source
		klog.V(1).Infof("Creating a new volume volumeID: %s from source: %+v", volumeID, volumeSource)
		copyvol := exec.Command("sudo", "dd", fmt.Sprintf("if=/dev/%s/%s", config.VGName, volumeSource), fmt.Sprintf("of=/dev/%s/%s", config.VGName, volumeName), "bs=1M")
		out, err := copyvol.CombinedOutput()
		if err != nil {
			klog.V(2).Infof("copy vol error: %v\n%s", err, out)
			http.Error(w, "Copy volume failed", http.StatusInternalServerError)
			return
		}
	}

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

func GetVolumeSize(volumePath string) (size int64, error error) {
	out, err := exec.Command("lvs", "--units", "b", "--nosuffix", "--noheadings", "-o", "lv_size", volumePath).Output()
	if err != nil {
		klog.Error("LV Size Error:", err)
		return
	}

	// Trim and clean the output
	output := strings.TrimSpace(string(out))

	// Optionally: print the raw output
	klog.V(5).Info("LV Size Raw output:", output)

	// Remove extra spaces (if needed)
	fields := strings.Fields(output)

	if len(fields) > 0 {
		sizeStr := fields[0] // assuming the number is the first field
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			klog.Error("LV Parse error:", err)
			return 0, err
		}
		klog.V(5).Info("LV Size in bytes:", size)
		return size, nil
	}
	return 0, err
}
