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
	VolumeID string `json:"snapshot_id"`
}

type DeleteSnapshotRequest struct {
	VolumeID string `json:"snapshot_id"`
}

type Config struct {
	VGName       string `yaml:"vg_name"`
	Port         string `yaml:"port"`
	Base_iqn     string `yaml:"iqn"`
	TargetPortal string `yaml:"target_portal"`
}
