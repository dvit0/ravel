package main

import (
	"github.com/valyentdev/ravel/pkg/driver"
)

func main() {
	firecrackerDriver := NewFirecrackerDriver()
	driver.ServeRavelDriver(firecrackerDriver)
}
