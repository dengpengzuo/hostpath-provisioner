package hostfsCsi

import (
	"context"
	"ez-cloud/hostpath-provisioner/pkg/csicommon"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type identifiedServer struct {
	*csicommon.DefaultIdentityServer
	info *driverInfo
}

// GetPluginInfo returns plugin information.
func (ids *identifiedServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	if ids.info.name == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured")
	}

	if ids.info.version == "" {
		return nil, status.Error(codes.Unavailable, "Driver is missing version")
	}

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
