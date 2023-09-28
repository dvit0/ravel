package init

import (
	"encoding/json"
	"os"
)

type InitConfig struct {
	ImageConfig ImageConfig
	CmdOverride []string
	RootDevice  string
	TTY         bool
	Hostname    string
	ExtraEnv    []string
	Mounts      []Mounts
	EtcResolv   EtcResolv
	EtcHost     []EtcHost
}

type ImageConfig struct {
	Cmd        []string
	Entrypoint []string
	Env        []string
	WorkingDir string
	User       string
}

type Mounts struct {
	MountPath  string
	DevicePath string
}

type EtcResolv struct {
	Nameservers []string
}

type EtcHost struct {
	Host string
	IP   string
	Desc string
}

func DecodeMachine(path string) (InitConfig, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return InitConfig{}, err
	}

	var config InitConfig

	if err := json.Unmarshal(contents, &config); err != nil {
		return InitConfig{}, err
	}

	return config, nil
}

func NewInitConfig(imageConfig ImageConfig) InitConfig {
	return InitConfig{
		ImageConfig: imageConfig,
		RootDevice:  "/dev/vdb",
		TTY:         true,
		ExtraEnv:    []string{},
	}
}
