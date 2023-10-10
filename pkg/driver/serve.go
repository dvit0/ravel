package driver

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func ServeRavelDriver(driver RavelDriver) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"RAVEL_DRIVER": &RavelDriverGRPCPlugin{Impl: driver},
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:  "ravel-driver",
			Level: hclog.LevelFromString("DEBUG"),
		}),
	})
}
