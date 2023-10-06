package driver

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/valyentdev/ravel/pkg/driver/proto"
)

var PluginMap = map[string]plugin.Plugin{
	"driver": &RavelDriverGRPCPlugin{},
}

type RavelDriver interface {
	StartVM(vmId string, vmConfig VMConfig) error
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
