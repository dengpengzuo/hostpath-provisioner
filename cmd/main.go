package main

import (
	"ez-cloud/hostpath-provisioner/pkg/hostfsCsi"
	"flag"
	"k8s.io/klog/v2"
)

const (
	CSI_NAME = "hostfs.csi.ezcloud.com"
	VERSION  = "1.0.0"
)

var nodeid = flag.String("nodeid", "", "node id")
var csiAddress = flag.String("csi-address", "/csi/csi.sock", "csi socket unix path")
var hostfsType = flag.String("hostfs-type", "", "ns,ids,cs")

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	csiDriver := hostfsCsi.NewHostfsCsiDriver(CSI_NAME, VERSION, *nodeid, *csiAddress)
	csiDriver.Start(*hostfsType)
}
