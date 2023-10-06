package driver

import (
	"context"

	"github.com/valyentdev/ravel/pkg/driver/proto"
)

type RavelDriverGRPCClient struct{ client proto.RavelDriverClient }

func (m *RavelDriverGRPCClient) StartVM(vmId string, vmConfig VMConfig) error {
	_, err := m.client.StartVM(context.Background(), vmConfig.ToProto())
	return err
}

func (m *RavelDriverGRPCClient) StopVM(vmId string) error {
	_, err := m.client.StopVM(context.Background(), &proto.StopVMRequest{
		VmId: vmId,
	})
	if err != nil {
		return err
	}

	return nil
}

type RavelDriverGRPCServer struct {
	proto.UnimplementedRavelDriverServer
	Impl RavelDriver
}

func (m *RavelDriverGRPCServer) StartVM(
	ctx context.Context,
	req *proto.StartVMRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.StartVM(req.VmId, VMConfigFromProto(req))
}

func (m *RavelDriverGRPCServer) StopVM(
	ctx context.Context,
	req *proto.StopVMRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.StopVM(req.VmId)
}
