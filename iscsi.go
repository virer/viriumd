package main

import (
	"fmt"
	"os/exec"
)

func createISCSITarget(volumeName string) error {
	vgPath := fmt.Sprintf("/dev/%s/%s", config.VGName, volumeName)

	iqn := fmt.Sprintf("iqn.2025-04.local.virium:%s", volumeName) // XXX FIXME date

	// Create backstore
	err := exec.Command("sudo", "targetcli", fmt.Sprintf("backstores/block create name=%s dev=%s", volumeName, vgPath))
	if err != nil {
		return fmt.Errorf("failed to create backstore: %s", err)
	}

	// Create target
	err = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/ create %s", iqn))
	if err != nil {
		return fmt.Errorf("failed to create iSCSI target: %s", err)
	}

	// Create LUN
	err = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/%s/tpg1/luns/ create /backstores/block/%s", iqn, volumeName))
	if err != nil {
		return fmt.Errorf("failed to create LUN: %s", err)
	}

	// Enable TPG1 and allow all initiators (simple setup; improve for prod!)
	err = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/%s/tpg1/acls/ create iqn.fake.initiator", iqn))
	if err != nil {
		return fmt.Errorf("failed to create ACL: %s", err)
	}

	return nil
}

func deleteISCSITarget(volumeName string) error {
	iqn := fmt.Sprintf("iqn.2025-04.local.virium:%s", volumeName) // XXX FIXME date

	// Delete iSCSI target — this removes LUNs and backstore link
	err := exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi delete %s", iqn))
	if err != nil {
		return fmt.Errorf("failed to delete iSCSI target: %s", err)
	}

	// Delete backstore
	err = exec.Command("sudo", "targetcli", fmt.Sprintf("backstores/block delete %s", volumeName))
	if err != nil {
		return fmt.Errorf("failed to delete backstore: %s", err)
	}

	return nil
}
