package hostfsCsi

import (
	"context"
	"ez-cloud/hostpath-provisioner/pkg/csicommon"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog"
	"os"
)

const (
	DefaultMountDir = "/works"
	DefaultDirModel = 0755
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
	info *driverInfo
}

//
// /var/lib/kubelet/plugins/{plugin}/volume/staging ...>>> Staging path
//                           create sub dir:  ..../${volumnId}
// POD 挂 FS volume: /var/lib/kubelet/pods/{podUid}/volume/{pluginName}/{volumeName}
//
// csiPlugin[csi-attacher] -> NodeStageVolume
func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.V(2).Infof("hostfs-csi NodeStageVolume{ StagingTargetPath := %s, VolumeId := %s } ...", req.StagingTargetPath, req.VolumeId)

	path := DefaultMountDir + "/" + req.VolumeId
	err := os.Mkdir(path, DefaultDirModel)
	if err != nil {
		return nil, fmt.Errorf("mkdir hostpath %s error:%s", path, err.Error())
	}
	klog.V(2).Infof("hostfs-csi create volume dir [%s] ... ", path)
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	klog.V(2).Infof("hostfs-csi NodeUnstageVolume{ StagingTargetPath := %s, VolumeId := %s } ...", req.StagingTargetPath, req.VolumeId)
	path := DefaultMountDir + "/" + req.VolumeId
	err := os.Link(path, path+".old")
	if err != nil {
		return nil, fmt.Errorf("mkdir hostpath %s error:%s", path, err.Error())
	}
	klog.V(2).Infof("hostfs-csi rename volume dir [%s] ... ", path)
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.V(2).Infof("hostfs-csi NodePublishVolume{ StagingTargetPath := %s, VolumeId := %s, TargetPath := %s } ...", req.StagingTargetPath, req.VolumeId, req.TargetPath)
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, in *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	return &csi.NodeUnpublishVolumeResponse{}, nil
}
