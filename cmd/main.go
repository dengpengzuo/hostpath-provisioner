package main

import (
	"ez-cloud/hostpath-provisioner/pkg/pvc"
	"flag"
	"k8s.io/klog/v2"
)

var hostDir = flag.String("host-dir", pvc.DefaultMountDir, "localpv host dirs")
var provisionerName = flag.String("provisioner-name", pvc.ProvisionerName, "external provisioner name ")

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	pvc.StartHotFsProvisioner(*provisionerName, *hostDir)
}
