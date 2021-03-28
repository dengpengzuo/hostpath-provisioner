package hostfsCsi

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

type driverInfo struct {
	name     string
	nodeId   string
	version  string
	endpoint string
}

type HostfsCsiDriver struct {
	info   driverInfo
	server *grpc.Server
	ids    *identifiedServer
	cs     *controllerServer
	ns     *nodeServer
}

func NewHostfsCsiDriver(name, version, nodeid string, endpoint string) *HostfsCsiDriver {
	return &HostfsCsiDriver{
		info: driverInfo{name: name, version: version, nodeId: nodeid, endpoint: endpoint},
	}
}

func (driver *HostfsCsiDriver) Start(ctrl ...string) error {
	proto, addr, err := parseEndpoint(driver.info.endpoint)
	if err != nil {
		return err
	}
	if proto == "unix" {
		addr = "/" + addr
		if e := os.Remove(addr); e != nil && !os.IsNotExist(e) {
			return fmt.Errorf("Failed to remove %s, error: %s", addr, e.Error())
		}
	}

	driver.server = grpc.NewServer(grpc.ConnectionTimeout(30 * time.Second))

	for i := 0; i < len(ctrl); i++ {
		v := ctrl[i]
		switch v {
		case "ids":
			driver.ids = &identifiedServer{info: &driver.info}
			csi.RegisterIdentityServer(driver.server, driver.ids)
		case "ns":
			driver.ns = &nodeServer{info: &driver.info}
			csi.RegisterNodeServer(driver.server, driver.ns)
		case "cs":
			driver.cs = &controllerServer{info: &driver.info}
			csi.RegisterControllerServer(driver.server, driver.cs)
		}
	}

	listener, err2 := net.Listen(proto, addr)
	if err2 != nil {
		return err2
	}

	return driver.server.Serve(listener)
}

func (driver *HostfsCsiDriver) Stop() {
	driver.server.GracefulStop()
}
