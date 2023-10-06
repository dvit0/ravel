package config

import (
	"os"

	"github.com/charmbracelet/log"
)

const (
	RAVEL_BASE_PATH           = "/var/lib/ravel"
	RAVEL_DRIVES_PATH         = RAVEL_BASE_PATH + "/drives/"
	RAVEL_MACHINES_SOCKS_PATH = RAVEL_BASE_PATH + "/machines"

	RAVEL_TEMP_BASE_PATH = "/tmp/ravel"
	RAVEL_TEMP_ARCHIVES  = RAVEL_TEMP_BASE_PATH + "/archives"
	RAVEL_FOLDERS_PERM   = 0644
)

const errorMessage = "Please check that you run the worker as root."

func InitRavelDirectories() {
	// check if ravel base path exists
	err := os.Mkdir(RAVEL_BASE_PATH, os.FileMode(RAVEL_FOLDERS_PERM))
	if err != nil && !os.IsExist(err) {
		log.Fatal(errorMessage, "error", err)
	}

	err = os.Mkdir(RAVEL_DRIVES_PATH, os.FileMode(RAVEL_FOLDERS_PERM))
	if err != nil && !os.IsExist(err) {
		log.Fatal(errorMessage, "error", err)
	}

	err = os.Mkdir(RAVEL_MACHINES_SOCKS_PATH, os.FileMode(RAVEL_FOLDERS_PERM))
	if err != nil && !os.IsExist(err) {
		log.Fatal(errorMessage, "error", err)
	}

	err = os.Mkdir(RAVEL_TEMP_BASE_PATH, os.FileMode(RAVEL_FOLDERS_PERM))
	if err != nil && !os.IsExist(err) {
		log.Fatal(errorMessage, "error", err)
	}

	err = os.Mkdir(RAVEL_TEMP_ARCHIVES, os.FileMode(RAVEL_FOLDERS_PERM))
	if err != nil && !os.IsExist(err) {
		log.Fatal(errorMessage, "error", err)
	}
}

func GetMachineSocketPath(machineId string) string {
	return RAVEL_MACHINES_SOCKS_PATH + "/" + machineId + ".sock"
}

func GetDriveImagePath(driveId string) string {
	return RAVEL_DRIVES_PATH + driveId + "/drive.img"
}
