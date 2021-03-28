package hostfsCsi

import "github.com/container-storage-interface/spec/lib/go/csi"

type controllerServer struct {
	info *driverInfo
	csi.ControllerServer
}
