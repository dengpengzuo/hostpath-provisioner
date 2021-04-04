package hostfsCsi

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"k8s.io/klog"
	"net"
	"os"
	"time"
)

type driverInfo struct {
	nodeId  string
	name    string
	version string
}

type HostfsCsiDriver struct {
	info    driverInfo
	address string

	server *grpc.Server

	ids *identifiedServer
	cs  *controllerServer
	ns  *nodeServer
}

func NewHostfsCsiDriver(name, version, nodeid string, address string) *HostfsCsiDriver {
	return &HostfsCsiDriver{
		address: address,
		info:    driverInfo{nodeId: nodeid, name: name, version: version},
	}
}

func (driver *HostfsCsiDriver) Start(ctrl ...string) error {
	addr := driver.address
	if driver.address[0] != '/' {
		addr = "/" + driver.address
	}
	if e := cleanupSocketFile(addr); e != nil {
		klog.Errorf("Remove socket: %s with error: %+v", addr, e)
		os.Exit(1)
	}

	listener, err2 := net.Listen("unix", addr)
	if err2 != nil {
		klog.Errorf("failed to listen on socket: %s with error: %+v", addr, err2)
		os.Exit(1)
	}

	klog.Infof("hostfs listen on socket: %s ", addr)
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

	klog.Infof("hostfs start grpc service ...")
	return driver.server.Serve(listener)
}

func (driver *HostfsCsiDriver) Stop() {
	driver.server.GracefulStop()
}

func cleanupSocketFile(socketPath string) error {
	socketExists, err := isSocketExist(socketPath)
	if err != nil {
		return err
	}
	if socketExists {
		if err := os.Remove(socketPath); err != nil {
			return fmt.Errorf("failed to remove stale socket %s with error: %+v", socketPath, err)
		}
	}
	return nil
}

func isSocketExist(socketPath string) (bool, error) {
	fi, err := os.Stat(socketPath)
	if err == nil && (fi.Mode()&os.ModeSocket) != 0 {
		return true, nil
	}
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("failed to stat the socket %s with error: %+v", socketPath, err)
	}
	return false, nil
}
