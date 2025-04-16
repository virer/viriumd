package main

import (
	"fmt"
	"log"
	"os/exec"
)

func createISCSITarget(volumeName string, volumeInitiator string) error {
	vgPath := fmt.Sprintf("/dev/%s/%s", config.VGName, volumeName)
	log.Println("iSCSI configuration for", volumeName)

	iqn := fmt.Sprintf("%s:%s", config.Base_iqn, volumeName)

	// Create backstore
	iscsicreate := exec.Command("sudo", "targetcli", fmt.Sprintf("backstores/block create name=%s dev=%s", volumeName, vgPath))
	out, err := iscsicreate.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create backstore: %s %s", err, out)
	}

	log.Println("iSCSI backstore created: ", volumeName)

	// Create target
	iscsicreate = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/ create %s", iqn))
	out, err = iscsicreate.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create iSCSI target: %s %s", err, out)
	}

	log.Println("iSCSI target created: ", volumeName)

	// Create LUN
	iscsicreate = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/%s/tpg1/luns/ create /backstores/block/%s", iqn, volumeName))
	out, err = iscsicreate.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create LUN: %s %s", err, out)
	}

	log.Println("iSCSI lun created: ", volumeName)

	// Enable TPG1 and allow all initiators ; simple setup improve for prod! XXX FIXME XXX
	iscsicreate = exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/%s/tpg1/acls/ create %s", iqn, volumeInitiator))
	out, err = iscsicreate.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create ACL: %s %s", err, out)
	}

	log.Println("iSCSI configuration done for: ", volumeName)

	return nil
}

func deleteISCSITarget(volumeName string) error {
	iqn := fmt.Sprintf("%s:%s", config.Base_iqn, volumeName)

	// Delete iSCSI target — this removes LUNs and backstore link
	iscsidel := exec.Command("sudo", "targetcli", fmt.Sprintf("iscsi/ delete %s", iqn))
	out, err := iscsidel.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete iSCSI target: %s %s", err, out)
	}

	// Delete backstore
	iscsidel = exec.Command("sudo", "targetcli", fmt.Sprintf("backstores/block delete name=%s", volumeName))
	out, err = iscsidel.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete backstore: %s %s", err, out)
	}

	return nil
}
