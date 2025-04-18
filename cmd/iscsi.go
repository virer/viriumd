package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func createISCSITarget(volumeID string, volumeName string, volumeInitiator string) ([]byte, error) {
	vgPath := fmt.Sprintf("/dev/%s/%s", config.VGName, volumeName)
	log.Println("iSCSI configuration for", volumeID)

	iqn := fmt.Sprintf("%s:%s", config.Base_iqn, volumeID)

	// Create backstore
	iscsicreate := exec.Command("sudo", "targetcli", fmt.Sprintf("backstores/block create name=%s dev=%s", volumeName, vgPath))
	out, err := iscsicreate.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create backstore: %s %s", err, out)
	}

	// log.Println("iSCSI backstore created: ", volumeName)

	// Create target
	iscsicreate = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/ create %s", iqn))
	out, err = iscsicreate.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create iSCSI target: %s %s", err, out)
	}

	// log.Println("iSCSI target created: ", volumeName)

	// Create LUN
	iscsicreate = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/%s/tpg1/luns/ create /backstores/block/%s", iqn, volumeName))
	out, err = iscsicreate.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create LUN: %s %s", err, out)
	}

	// log.Println("iSCSI lun created: ", volumeName)
	// Attribute
	iscsicreate = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/%s/tpg1/", iqn), "set", "attribute", "generate_node_acls=1")
	out, err = iscsicreate.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to set attributes: %s %s", err, out)
	}
	// Write protection off
	iscsicreate = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/%s/tpg1/", iqn), "set", "attribute", "demo_mode_write_protect=0")
	out, err = iscsicreate.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to set attributes: %s %s", err, out)
	}

	// Enable TPG1 and ACL to allow initiator to R/W access (default)
	iscsicreate = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/%s/tpg1/acls/ create %s", iqn, volumeInitiator))
	out, err = iscsicreate.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create ACL: %s %s", err, out)
	}
	log.Println("iSCSI configuration done for:", volumeID)

	return []byte(iqn), nil
}

func deleteISCSITarget(volumeID, volumeName string) error {
	iqn := fmt.Sprintf("%s:%s", config.Base_iqn, volumeID)

	// Delete iSCSI target — this removes LUNs and backstore link
	iscsidel := exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/ delete %s", iqn))
	out, err := iscsidel.CombinedOutput()
	if err != nil {
		if strings.HasPrefix(string(out), "No such Target in configfs") {
			fmt.Println("iSCSI target already removed, trying to remove backstore...")
		} else {
			return fmt.Errorf("failed to delete iSCSI target: %s %s", err, out)
		}
	}

	// Delete backstore
	iscsidel = exec.Command("sudo", "targetcli", fmt.Sprintf("backstores/block delete name=%s", volumeName))
	out, err = iscsidel.CombinedOutput()
	if err != nil {
		if strings.HasPrefix(string(out), "No storage object named") {
			fmt.Println("iSCSI backstore already removed.")
		} else {
			return fmt.Errorf("failed to delete backstore: %s %s", err, out)
		}
	}

	return nil
}
