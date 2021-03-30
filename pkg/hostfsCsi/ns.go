package hostfsCsi

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
)

type nodeServer struct {
	info *driverInfo
	csi.NodeServer
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, in *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	panic("implement me")
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, in *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	panic("implement me")
}
