package hostfscsi

import (
	"context"
	"ez-cloud/hostpath-provisioner/pkg/csicommon"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog"
	"os"
	"path/filepath"
)

const (
	DefaultDirModel = 0755
	DefaultWorkDir  = "/work"
)

type nodeServer struct {
	csicommon.DefaultNodeServer
	info *driverInfo
}

// 必须实现.
// kubelet.pluginManager(CSIPlugin, csi.RegistrationHandler)
// 内部调用了 Node.GetInfo.
//
func (ns *nodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: ns.info.nodeId,
	}, nil
}

// csiPlugin[csi-attacher] -> NodeStageVolume
func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.Infof("hostfs-csi NodeStageVolume{ StagingTargetPath := %s, VolumeId := %s } ...", req.StagingTargetPath, req.VolumeId)
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	klog.Infof("hostfs-csi NodeUnstageVolume{ StagingTargetPath := %s, VolumeId := %s } ...", req.StagingTargetPath, req.VolumeId)
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.Infof("hostfs-csi NodePublishVolume{ StagingTargetPath := %s, VolumeId := %s, TargetPath := %s } ...", req.StagingTargetPath, req.VolumeId, req.TargetPath)

	path := filepath.Join(DefaultWorkDir, req.VolumeId)
	err := os.Mkdir(path, DefaultDirModel)
	if err != nil {
		return nil, fmt.Errorf("hostfs-csi NodePublishVolume dir [%s] error:%v", path, err)
	}

	opt := []string{"bind"}
	err = UtilMount(path, req.TargetPath, "", opt[:])
	if err != nil {
		klog.Error("hostfs-csi NodePublishVolume dir [%s], error: %v ", req.VolumeId, err)
		return nil, err
	}
	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.Infof("hostfs-csi NodeUnpublishVolume{ VolumeId := %s, TargetPath := %s } ...", req.VolumeId, req.TargetPath)

	err := UtilUnMount(req.TargetPath)
	if err != nil {
		klog.Error("hostfs-csi NodeUnpublishVolume dir [%s], error: %v ", req.VolumeId, err)
		return nil, err
	}

	path := filepath.Join(DefaultWorkDir, req.VolumeId)
	err = os.Link(path, "."+path)
	if err != nil {
		return nil, fmt.Errorf("hostfs-csi NodeUnpublishVolume dir [%s], error: %v ", path, err)
	}
	return &csi.NodeUnpublishVolumeResponse{}, nil
}