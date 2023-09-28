package utils

import (
	"os"

	"github.com/charmbracelet/log"
	"golang.org/x/sys/unix"
)

func Fallocate(filePath string, limitBytes int64) error {
	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		log.Error("Error opening file", "err", err)
		return err
	}
	defer file.Close()

	// Perform the actual fallocate
	err = unix.Fallocate(int(file.Fd()), 0, 0, limitBytes)
	if err != nil {
		log.Error("Error during fallocate", "err", err)
		return err
	}

	return nil
}
