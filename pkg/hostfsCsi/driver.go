package hostfsCsi

type HostfsCsiDriver struct {
    endpoint  string
    csiDriver *csicommon.CSIDriver

    ids *identifiedServer
    cs  *controllerServer
    ns  *nodeServer
}

type controllerServer struct {
    *csicommon.DefaultControllerServer
}

type identifiedServer struct {
    *csicommon.DefaultIdentityServer
}

type nodeServer struct {
    *csicommon.DefaultNodeServer
}

func newCSIDriver(nodeID string) *csicommon.CSIDriver {
    csiDriver := csicommon.NewCSIDriver(driverName, version, nodeID)
    csiDriver.AddControllerServiceCapabilities(
        []csi.ControllerServiceCapability_RPC_Type{
            csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
            csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
            csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
        })
    csiDriver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})
    return csiDriver
}
