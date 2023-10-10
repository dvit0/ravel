package driversmanager

import (
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/valyentdev/ravel/internal/worker/store"
	"github.com/valyentdev/ravel/pkg/driver"
)

type DriversManager struct {
	clients map[string]*plugin.Client
	store   *store.Store
}

func NewDriversManager(store *store.Store) *DriversManager {
	return &DriversManager{
		clients: map[string]*plugin.Client{},
		store:   store,
	}
}

func (driversManager *DriversManager) GetDriver(name string) (driver.RavelDriver, error) {
	log.Info("Getting driver", "name", name)
	client, ok := driversManager.clients[name]

	if ok && !client.Exited() {
		log.Info("Driver found, reusing", "name", name)
		clientProto, err := client.Client()
		if err != nil {
			return nil, err
		}

		return driverFromClientProtocol(clientProto)
	}

	log.Info("Driver not found", "name", name)
	log.Info("Checking for a reattach config driver", "name", name)
	reattachConfig, err := driversManager.store.GetDriverReattachConfig(name)
	if err != nil {
		return nil, err
	}
	if reattachConfig != nil {
		log.Info("Reattach config found for", "driver", name)
	}

	log.Info("Creating new driver client from reattach config", "driver", name)
	client = newDriverClient(name, reattachConfig)

	clientProto, err := client.Client()
	if err != nil {
		log.Info("Can't creating new driver client from reattach config", "error", err)
		log.Info("Creating new driver client from scratch", "driver", name)
		client = newDriverClient(name, nil)
		clientProto, err = client.Client()
		if err != nil {
			return nil, err
		}
	}

	driver, err := driverFromClientProtocol(clientProto)
	if err != nil {
		return nil, err
	}

	driversManager.clients[name] = client

	err = driversManager.store.SetDriverReattachConfig(name, client.ReattachConfig())
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func driverFromClientProtocol(client plugin.ClientProtocol) (driver.RavelDriver, error) {

	raw, err := client.Dispense("RAVEL_DRIVER")
	if err != nil {
		log.Error("Failed to request driver instance", "error", err)
		return nil, err
	}

	return raw.(driver.RavelDriver), nil

}

func newDriverClient(driverName string, reattachConfig *plugin.ReattachConfig) *plugin.Client {

	clientConfig := &plugin.ClientConfig{
		HandshakeConfig: driver.Handshake,
		Plugins: plugin.PluginSet{
			"RAVEL_DRIVER": &driver.RavelDriverGRPCPlugin{},
		},
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:  "ravel-driver",
			Level: hclog.LevelFromString("DEBUG"),
		}),
	}

	if reattachConfig != nil {
		clientConfig.Reattach = reattachConfig
	} else {
		clientConfig.Cmd = exec.Command("sh", "-c", "./bin/drivers/"+driverName)
	}
	client := plugin.NewClient(clientConfig)

	return client
}
