package driver

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/valyentdev/ravel/pkg/driver/proto"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "RAVEL_DRIVER_PLUGIN",
	MagicCookieValue: "ravel-driver",
}

type RavelDriver interface {
	StartVM(vmId string, vmConfig VMConfig) (*proto.StartVMResponse, error)
	StopVM(vmId string) error
}

type RavelDriverGRPCPlugin struct {
	plugin.Plugin
	Impl RavelDriver
}

func (p *RavelDriverGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRavelDriverServer(s, &RavelDriverGRPCServer{Impl: p.Impl})
	return nil
}

func (p *RavelDriverGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &RavelDriverGRPCClient{client: proto.NewRavelDriverClient(c)}, nil
}
