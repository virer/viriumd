package main

type VolumeRequest struct {
	InitiatorName string `json:"initiator_name"`
	Capacity      int64  `json:"capacity"` // bytes
}
type VolumeResizeRequest struct {
	VolumeID string `json:"volume_id"`
	Capacity int64  `json:"capacity"` // bytes
}

type VolumeResponse struct {
	VolumeID          string `json:"volume_id"`
	TargetPortal      string `json:"targetPortal"`
	Iqn               string `json:"iqn"`
	Lun               string `json:"lun"`
	DiscoveryCHAPAuth string `json:"discoveryCHAPAuth"`
	SessionCHAPAuth   string `json:"sessionCHAPAuth"`
}

type DeleteVolumeRequest struct {
	VolumeID string `json:"volume_id"`
}

type SnapshotRequest struct {
	Name     string `json:"name"`
	VolumeID string `json:"source_volume_id"`
}

type DeleteSnapshotRequest struct {
	VolumeID string `json:"snapshot_id"`
}

type Config struct {
	VGName       string `yaml:"vg_name"`
	Port         string `yaml:"port"`
	Base_iqn     string `yaml:"iqn"`
	TargetPortal string `yaml:"target_portal"`
	Cmd_prefix   string `yaml:"cmd_prefix"`
	API_username string `yaml:"api_username"`
	API_password string `yaml:"api_password"`
}

// Constructor function that sets default values
func NewConfiguration() *Config {
	return &Config{
		VGName:       "vg_data",                      // Default LVM Volume Group is vg_data
		Port:         "8787",                         // Default http port is 8787
		Base_iqn:     "iqn.2025-04.net.virer.virium", // Default iSCSI iqn is iqn.2025-04.net.virer.virium
		TargetPortal: "127.0.0.1:3260",               // Default target portal is localhost:3260
		Cmd_prefix:   "sudo",                         // Default command prefix is sudo
		API_username: "virium_api_username",          // Default API username is virium_api_username
		API_password: "virium_api_password",          // Default API password is virium_api_password
	}
}
