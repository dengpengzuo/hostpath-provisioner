package hostfsCsi

import (
	"errors"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

type driverInfo struct {
	name    string
	version string
}

type HostfsCsiDriver struct {
	info    driverInfo
	nodeId  string
	address string

	server *grpc.Server

	ids *identifiedServer
	cs  *controllerServer
	ns  *nodeServer
}

func NewHostfsCsiDriver(name, version, nodeid string, address string) *HostfsCsiDriver {
	return &HostfsCsiDriver{
		nodeId:  nodeid,
		address: address,
		info:    driverInfo{name: name, version: version},
	}
}

func (driver *HostfsCsiDriver) Start(ctrl ...string) error {
	addr := driver.address
	if driver.address[0] != '/' {
		addr = "/" + driver.address
	}
	if e := os.Remove(addr); e != nil && !os.IsNotExist(e) {
		return errors.New(fmt.Sprintf("Failed to remove unix://%s error:%v", addr, e))
	}

	driver.server = grpc.NewServer(grpc.ConnectionTimeout(30 * time.Second))

	for i := 0; i < len(ctrl); i++ {
		v := ctrl[i]
		switch v {
		case "ids":
			driver.ids = &identifiedServer{info: &driver.info}
			csi.RegisterIdentityServer(driver.server, driver.ids)
		case "cs":
			driver.cs = &controllerServer{info: &driver.info}
			csi.RegisterControllerServer(driver.server, driver.cs)
		case "ns":
			driver.ns = &nodeServer{info: &driver.info}
			csi.RegisterNodeServer(driver.server, driver.ns)
		}
	}

	listener, err2 := net.Listen("unix", addr)
	if err2 != nil {
		return err2
	}

	return driver.server.Serve(listener)
}

func (driver *HostfsCsiDriver) Stop() {
	driver.server.GracefulStop()
}
