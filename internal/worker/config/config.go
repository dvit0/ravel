package config

import "os"

const (
	RAVEL_BASE_PATH           = "/var/lib/ravel"
	RAVEL_DRIVES_PATH         = RAVEL_BASE_PATH + "/drives/"
	RAVEL_MACHINES_SOCKS_PATH = RAVEL_BASE_PATH + "/machines"

	RAVEL_TEMP_BASE_PATH = "/tmp/ravel"
	RAVEL_TEMP_ARCHIVES  = RAVEL_TEMP_BASE_PATH + "/archives"
)

func InitRavelDirectories() {
	os.Mkdir(RAVEL_BASE_PATH, os.FileMode(0644))
	os.Mkdir(RAVEL_DRIVES_PATH, os.FileMode(0644))
	os.Mkdir(RAVEL_MACHINES_SOCKS_PATH, os.FileMode(0644))

	os.Mkdir(RAVEL_TEMP_BASE_PATH, os.FileMode(0644))
	os.Mkdir(RAVEL_TEMP_ARCHIVES, os.FileMode(0644))
}

func GetMachineSocketPath(machineId string) string {
	return RAVEL_MACHINES_SOCKS_PATH + "/" + machineId + ".sock"
}

func GetDriveImagePath(driveId string) string {
	return RAVEL_DRIVES_PATH + driveId + "/drive.img"
}
