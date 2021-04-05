package hostfscsi

import (
	"context"
	"ez-cloud/hostpath-provisioner/pkg/csicommon"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/google/uuid"
)

type controllerServer struct {
	csicommon.DefaultControllerServer
	info *driverInfo
}

// csi-provisioner
func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	// 生成一个volumeId.
	volume := &csi.Volume{
		VolumeId:      "hostfs-" + uuid.New().String(),
		CapacityBytes: req.CapacityRange.RequiredBytes,
		VolumeContext: req.GetParameters(),
		ContentSource: req.GetVolumeContentSource(),
	}
	return &csi.CreateVolumeResponse{Volume: volume}, nil
}

func (cs *controllerServer) DeleteVolume(context.Context, *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
					},
				},
			},
		}}, nil
}

// csi-driver.spec.attachRequired = true
// csi-attacher call
func (cs *controllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return &csi.ControllerPublishVolumeResponse{}, nil
}
func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return &csi.ControllerUnpublishVolumeResponse{}, nil
}
