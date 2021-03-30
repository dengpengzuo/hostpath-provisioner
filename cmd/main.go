package main

import (
	"ez-cloud/hostpath-provisioner/pkg/hostfsCsi"
	"ez-cloud/hostpath-provisioner/pkg/pvc"
	"flag"
	"k8s.io/klog/v2"
)

const (
	CSI_NAME = "csi-hostfs"
	VERSION  = "1.0.0"
)

var nodeid = flag.String("nodeid", "", "node id")
var endpoint = flag.String("endpoint", pvc.CsiEndpoint, "csi socket unix path")
var hostfsType = flag.String("hostfs-type", "", "ns,ids,cs")

var hostDir = flag.String("host-dir", pvc.DefaultMountDir, "host dirs")
var provisionerName = flag.String("provisioner-name", pvc.ProvisionerName, "external provisioner name ")

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	csiDriver := hostfsCsi.NewHostfsCsiDriver(CSI_NAME, VERSION, *nodeid, *endpoint)
	csiDriver.Start(*hostfsType)
}
