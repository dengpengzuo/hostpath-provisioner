### How to write a CSI Driver?
    The kubernetes-csi site details how to develop, deploy, and test a CSI driver on Kubernetes. In general, CSI Drivers should be deployed on Kubernetes along with the following sidecar (helper) containers:

#### external-attacher
Watches Kubernetes VolumeAttachment objects and triggers ControllerPublish and ControllerUnpublish operations against a CSI endpoint.

#### external-provisioner
Watches Kubernetes PersistentVolumeClaim objects and triggers CreateVolume and DeleteVolume operations against a CSI endpoint.

#### node-driver-registrar
Registers the CSI driver with kubelet using the Kubelet device plugin mechanism.

#### cluster-driver-registrar (Alpha)
Registers a CSI Driver with the Kubernetes cluster by creating a CSIDriver object which enables the driver to customize how Kubernetes interacts with it.

#### external-snapshotter (Alpha)
Watches Kubernetes VolumeSnapshot CRD objects and triggers CreateSnapshot and DeleteSnapshot operations against a CSI endpoint.

#### livenessprobe
May be included in a CSI plugin pod to enable the Kubernetes Liveness Probe mechanism.
Storage vendors can build Kubernetes deployments for their plugins using these components, while leaving their CSI driver completely unaware of Kubernetes.

