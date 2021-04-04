package hostfsCsi

import (
	"context"
	"ez-cloud/hostpath-provisioner/pkg/csicommon"
	"github.com/container-storage-interface/spec/lib/go/csi"
)

type identifiedServer struct {
	csicommon.DefaultIdentityServer
	info *driverInfo
}

// GetPluginInfo returns plugin information.
func (ids *identifiedServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return &csi.GetPluginInfoResponse{
		Name:          ids.info.name,
		VendorVersion: ids.info.version,
	}, nil
}

// Probe returns empty response.
func (ids *identifiedServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{}, nil
}

// GetPluginCapabilities returns plugin capabilities.
func (ids *identifiedServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}, nil
}
