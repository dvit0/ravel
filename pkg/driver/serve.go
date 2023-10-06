package driver

import "github.com/hashicorp/go-plugin"

func ServeRavelDriver(name string, driver RavelDriver) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{ // will be versioned later
			MagicCookieKey:   "RAVEL_DRIVER_PLUGIN",
			MagicCookieValue: "ravel-driver",
		},
		Plugins: map[string]plugin.Plugin{
			name: &RavelDriverGRPCPlugin{Impl: driver},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
