package main

type VolumeRequest struct {
	InitiatorName string `json:"initiator_name"`
	Capacity      int64  `json:"capacity"` // bytes
}
type VolumeResizeRequest struct {
	VolumeID string `json:"id"`
	Capacity int64  `json:"capacity"` // bytes
}

type VolumeResponse struct {
	VolumeID     string `json:"id"`
	TargetPortal string `json:"targetPortal"`
	Iqn          string `json:"iqn"`
	Lun          string `json:"lun"`
}

type DeleteVolumeRequest struct {
	VolumeID string `json:"id"`
}

type SnapshotRequest struct {
	VolumeID string `json:"id"`
}

type DeleteSnapshotRequest struct {
	VolumeID string `json:"id"`
}

type Config struct {
	VGName       string `yaml:"vg_name"`
	Port         string `yaml:"port"`
	Base_iqn     string `yaml:"iqn"`
	TargetPortal string `yaml:"target_portal"`
}
