## Dynamic localpv provisioner for Kubernetes

### 1. 发布步骤
```bash
kubectl apply -f deploy/hostpath-provisioner-rbac.yaml
kubectl apply -f deploy/hostpath-provisioner.yaml
kubectl apply -f deploy/storageclass.yaml
kubectl apply -f deploy/example.yaml
```