package main

import (
	"fmt"
	"os/exec"
	"strings"

	klog "k8s.io/klog/v2"
)

func createISCSITarget(volumeID string, volumeName string, volumeInitiator string) ([]byte, error) {
	vgPath := fmt.Sprintf("/dev/%s/%s", config.VGName, volumeName)
	klog.V(2).Info("iSCSI configuration for: ", volumeID)

	iqn := fmt.Sprintf("%s:%s", config.Base_iqn, volumeID)

	commands := []struct {
		Command string
		Message string
	}{
		{fmt.Sprintf("targetcli backstores/block create name=%s dev=%s", volumeName, vgPath), "failed to create backstore"},
		{fmt.Sprintf("targetcli iscsi/ create %s", iqn), "failed to create iSCSI target"},
		{fmt.Sprintf("targetcli iscsi/%s/tpg1/luns/ create /backstores/block/%s", iqn, volumeName), "failed to create LUN"},
		{fmt.Sprintf("targetcli iscsi/%s/tpg1/ set attribute generate_node_acls=1", iqn), "failed to set attributes"},
		{fmt.Sprintf("targetcli iscsi/%s/tpg1/ set attribute demo_mode_write_protect=0", iqn), "failed to set attributes"},
		{fmt.Sprintf("targetcli iscsi/%s/tpg1/acls/ create %s", iqn, volumeInitiator), "failed to create ACL"},
		{"targetcli saveconfig", "failed to save configuration"},
	}
	for _, c := range commands {
		_, err := executeWrapper(c.Command, c.Message)
		if err != nil {
			klog.Error("Error:", err)
			break
		}
	}

	klog.V(2).Info("iSCSI configuration done for: ", volumeID)

	return []byte(iqn), nil
}

func deleteISCSITarget(volumeID, volumeName string) error {
	iqn := fmt.Sprintf("%s:%s", config.Base_iqn, volumeID)

	// Delete iSCSI target — this removes LUNs and backstore link
	iscsidel := exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/ delete %s", iqn))
	out, err := iscsidel.CombinedOutput()
	if err != nil {
		if strings.HasPrefix(string(out), "No such Target in configfs") {
			klog.V(2).Info("iSCSI target already removed, trying to remove backstore...")
		} else {
			return fmt.Errorf("failed to delete iSCSI target: %s %s", err, out)
		}
	}

	// Delete backstore
	iscsidel = exec.Command("sudo", "targetcli", fmt.Sprintf("backstores/block delete name=%s", volumeName))
	out, err = iscsidel.CombinedOutput()
	if err != nil {
		if strings.HasPrefix(string(out), "No storage object named") {
			klog.V(2).Info("iSCSI backstore already removed.")
		} else {
			return fmt.Errorf("failed to delete backstore: %s %s", err, out)
		}
	}

	return nil
}
