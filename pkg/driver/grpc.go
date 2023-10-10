package driver

import (
	"context"

	"github.com/valyentdev/ravel/pkg/driver/proto"
)

type RavelDriverGRPCClient struct{ client proto.RavelDriverClient }

func (m *RavelDriverGRPCClient) StartVM(vmId string, vmConfig VMConfig) (*proto.StartVMResponse, error) {
	config := vmConfig.ToProto()
	config.VmId = vmId

	return m.client.StartVM(context.Background(), config)
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
	req *proto.StartVMRequest) (*proto.StartVMResponse, error) {
	return m.Impl.StartVM(req.VmId, VMConfigFromProto(req))
}

func (m *RavelDriverGRPCServer) StopVM(
	ctx context.Context,
	req *proto.StopVMRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.StopVM(req.VmId)
}
