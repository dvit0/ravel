package main

import (
	"github.com/valyentdev/ravel/pkg/driver"
)

func main() {
	firecrackerDriver := NewFirecrackerDriver()
	defer firecrackerDriver.cleanup()
	driver.ServeRavelDriver("firecracker", firecrackerDriver)
}
