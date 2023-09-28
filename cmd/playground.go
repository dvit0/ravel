package main

import (
	"errors"
	"fmt"
)

const (
	EXT4 Format = "ext4"
)

type Format string

func GetFormat(format Format) {
	fmt.Println(format)
}

func main() {

	err := errors.New("machine not found")

	defer func() {
		if err != nil {
			fmt.Println("Removing file")
		}
	}()

	err = nil

}
