package types

type RavelMachineSpec struct {
	Image  string `json:"image"`
	VCpus  int64  `json:"vcpus"`
	Memory int64  `json:"memory"`
}

type RavelMachine struct {
	*RavelMachineSpec
	Id          string             `json:"id"`
	InitDriveId string             `json:"init_drive_id"`
	RootDriveId string             `json:"root_drive_id"`
	Status      RavelMachineStatus `json:"status"`
}

type RavelMachineStatus string

const (
	RavelMachineStatusRunning RavelMachineStatus = "running"
	RavelMachineStatusStopped RavelMachineStatus = "stopped"
	RavelMachineStatusCreated RavelMachineStatus = "created"
	RavelMachineStatusError   RavelMachineStatus = "error"
)
