package pvc

import (
	"context"
	"errors"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"path/filepath"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v6/controller"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"os"
)

type hostpathProvisioner struct {
	provisionerName string
	mountDir        string
	identity        string
	allocInf        DeviceAllocInf
}

func StartHotFsProvisioner(provisionerName, vgName string) {
	// Create an InClusterConfig and use it to create a client for the controller
	// to use to communicate with Kubernetes
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("Failed to create config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Failed to create client: %v", err)
	}

	// The controller needs to know what the server version is because out-of-tree
	// provisioners aren't officially supported until 1.5
	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		klog.Fatalf("Error getting server version: %v", err)
	}

	devInf := NewhostpathDevice(vgName)
	if err = devInf.Init(); err != nil {
		klog.Fatalf("Error init hostpath %s err: %v", vgName, err)
	}

	if _, isOk := devInf.(DeviceAllocInf); !isOk {
		klog.Fatalf("hostpath devInf must implements DeviceAllocInf !")
	}

	// Create the provisioner: it implements the Provisioner interface expected by the controller
	hostpathProvisioner := NewhostpathProvisioner(provisionerName, vgName, devInf.(DeviceAllocInf))

	// Start the provision controller which will dynamically provision hostPath PVs
	pc := controller.NewProvisionController(clientset, provisionerName, hostpathProvisioner, serverVersion.GitVersion)

	// Never stops.
	pc.Run(context.Background())
}

// NewhostpathProvisioner creates a new hostpath provisioner
func NewhostpathProvisioner(provisionerName, mountDir string, allocInf DeviceAllocInf) controller.Provisioner {
	nodeName := os.Getenv("MY_NODE_NAME")
	if nodeName == "" {
		klog.Fatal("env variable MY_NODE_NAME must be set so that this provisioner can identify itself")
	}

	return &hostpathProvisioner{
		provisionerName: provisionerName,
		mountDir:        mountDir,
		identity:        nodeName,
		allocInf:        allocInf,
	}
}

// Provision creates a storage asset and returns a PV object representing it.
func (p *hostpathProvisioner) Provision(ctx context.Context, options controller.ProvisionOptions) (*v1.PersistentVolume, controller.ProvisioningState, error) {
	if p.identity != options.SelectedNode.Name {
		return nil, controller.ProvisioningFinished, &controller.IgnoredError{Reason: "pv identity annotation does not match selectNode"}
	}

	pvsize := options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	fstype := options.StorageClass.Parameters["fsType"]

	if err := p.allocInf.Alloc(options.PVName, pvsize, fstype); err != nil {
		return nil, controller.ProvisioningReschedule, err
	}

	outsidePath := filepath.Join(p.mountDir, options.PVName)
	nodeAffinity, e := generateVolumeNodeAffinity(options.SelectedNode)
	if e != nil {
		return nil, controller.ProvisioningReschedule, e
	}
	var pv = &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
			Annotations: map[string]string{
				ProvisionerId: p.identity,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: *options.StorageClass.ReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				Local: &v1.LocalVolumeSource{
					Path:   outsidePath,
					FSType: &fstype,
				},
			},
			NodeAffinity:     nodeAffinity,
			StorageClassName: options.StorageClass.Name,
			VolumeMode:       options.PVC.Spec.VolumeMode,
			MountOptions:     options.StorageClass.MountOptions,
		},
	}
	return pv, controller.ProvisioningFinished, nil
}

// Delete removes the storage asset that was created by Provision represented
// by the given PV.
func (p *hostpathProvisioner) Delete(ctx context.Context, pv *v1.PersistentVolume) error {
	ann, ok := pv.Annotations[ProvisionerId]
	if !ok {
		return errors.New("identity annotation not found on PV")
	}

	if ann != p.identity {
		return &controller.IgnoredError{Reason: "identity annotation on PV does not match ours"}
	}

	if pv.Spec.Local == nil {
		return &controller.IgnoredError{Reason: "identity annotation on PV does not match ours"}
	}

	fsType := pv.Spec.Local.FSType

	return p.allocInf.Remove(pv.Name, pv.Spec.PersistentVolumeReclaimPolicy == v1.PersistentVolumeReclaimDelete, fsType)
}

func generateVolumeNodeAffinity(node *v1.Node) (*v1.VolumeNodeAffinity, error) {
	if node.Labels == nil {
		return nil, fmt.Errorf("Node does not have labels")
	}
	nodeValue, found := node.Labels[v1.LabelHostname]
	if !found {
		return nil, fmt.Errorf("Node does not have expected label %s", v1.LabelHostname)
	}

	return &v1.VolumeNodeAffinity{
		Required: &v1.NodeSelector{
			NodeSelectorTerms: []v1.NodeSelectorTerm{
				{
					MatchExpressions: []v1.NodeSelectorRequirement{
						{
							Key:      v1.LabelHostname,
							Operator: v1.NodeSelectorOpIn,
							Values:   []string{nodeValue},
						},
					},
				},
			},
		},
	}, nil
}
