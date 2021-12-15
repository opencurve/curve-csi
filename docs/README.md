# CSI Curve Driver

This document provides more detail about curve-csi driver.

- [Deploy CurveBS CSI](#deploy-curvebs-csi)
  - [Requirements](#requirements)
  - [Using the helm chart](#using-the-helm-chart)
  - [Using the kubernetes manifests](#using-the-kubernetes-manifests)
- [Deploy CurveFS CSI](#deploy-curvefs-csi)
- [Debug](#debug)
- [Examples](#examples)
  - [CurveBS](#curvebs)
    - [Create StorageClass](#create-storageclass)
    - [Create PersistentVolumeClaim](#create-persistentvolumeclaim)
    - [Create Test Pod](#create-test-pod)
  - [CurveFS](#curvefs)
  - [Block volume](#block-volume)
  - [Volume expanding](#volume-expanding)
  - [Snapshot](#snapshot)
  - [Volume clone](#volume-clone)
- [Test](#test)
  - [Test Using CSC Tool](#test-using-csc-tool)
  - [E2E](#e2e)


## Deploy CurveBS CSI

#### Requirements

The curve-csi driver deploys on the hosts contain "Master" hosts and "Node" hosts.

- On the "Master" hosts, install `curve tool` which communicates with the curve cluster to manage the volume lifecycle, such as Create/Delete/Expand/Snapshot.
- On the "Node" hosts, install `curve-nbd tool` which allows attaching/detaching volumes to workloads.

Please refer to [deploy doc](https://github.com/opencurve/curve/blob/master/docs/cn/deploy.md) to get how to install `curve` and `curve-nbd` tool.

#### Using the helm chart

```bash
helm install --namespace "csi-system" charts/curve-csi-curvebs
```

#### Using the kubernetes manifests

Change to the `deploy/curve-csi-curvebs/` directory, create the files:

```bash
kubectl apply -f ./*.yaml
```

## Deploy CurveFS CSI



## Debug

You can dynamically set the log level by enabling the driver parameter `--debug-port`,
and call:

```text
curl -XPUT http://127.0.0.1:<debugPort>/debug/flags/v -d '5'
```

## Examples

### CurveBS

#### Create StorageClass

```yaml
allowVolumeExpansion: true
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: curvebs
parameters:
  user: k8s  # the user of curve volumes which you want to create
  cloneLazy: "true"
provisioner: curvebs.csi.netease.com
reclaimPolicy: Delete
volumeBindingMode: Immediate
```

#### Create PersistentVolumeClaim

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: curvebs-test-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi # range of 10Gi~4Ti, and step by 1Gi
  storageClassName: curvebs
```

#### Create Test Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: csi-curvebs-test
spec:
  containers:
  - name: web-server
    image: nginx
    volumeMounts:
    - name: mypvc
      mountPath: /var/lib/www/html
  volumes:
  - name: mypvc
    persistentVolumeClaim:
      claimName: curvebs-test-pvc
```

### CurveFS


### Block volume

Refer to [block-volume](./block-volume.md)

### Volume expanding

Refer to [expand-pvc](./expand-pvc.md)

### Snapshot

Refer to [snapshot](./snapshot.md)

### Volume clone

Refer to [clone](./clone.md)

## Test

### Test Using CSC Tool

Refer to [csc-test](./csc-test.md)

### E2E 

Refer to [e2e](./e2e.md)
