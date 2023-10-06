package main

import (
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

	for i := 0; i < 65536; i++ {
	}
	fmt.Println("hey")

}
