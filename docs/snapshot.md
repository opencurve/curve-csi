- [Prerequisite](#prerequisite)
- [Create SnapshotClass](#create-snapshotclass)
- [Create Snapshot](#create-snapshot)
- [Restore Snapshot to a new PVC](#restore-snapshot-to-a-new-pvc)

## Prerequisite

For snapshot functionality to be supported for your Kubernetes cluster, the Kubernetes version running in your cluster should be `>= v1.17`. We also need the snapshot controller deployed in your Kubernetes cluster along with csi-snapshotter sidecar container.

**Git Repository:**  https://github.com/kubernetes-csi/external-snapshotter

**Supported Versions**

|Latest stable release	|Min CSI Version	|Max CSI Version	|Container Image	|Min K8s Version	|Max K8s Version	|Recommended K8s Version|
| ---	| --- 	| ---	| ---	| --- |---	|---|
|v4.2.1	|	v1.0.0	|-	|registry.k8s.io/sig-storage/snapshot-controller:v4.2.1|	v1.20	|-	|v1.22|
| v4.1.1|	v1.0.0	|-	|registry.k8s.io/sig-storage/snapshot-controller:v4.1.1|	v1.20|	-	|v1.20|
| v3.0.3 (beta)|	v1.0.0|	-|	registry.k8s.io/sig-storage/snapshot-controller:v3.0.3	|v1.17	|-|	v1.17|


**Install Snapshot Beta CRDs:**

```bash
kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml

kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml

kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml
```

**Install Snapshot Controller:**

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml

kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml
```

## Create SnapshotClass

```bash
kubectl create -f ../examples/snapshotclass.yaml
```

## Create Snapshot

- Verify if PVC is in Bound state

```bash
$ kubectl get pvc
NAME                   STATUS        VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curve-test-pvc         Bound         pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   30Gi       RWO            curve          35m
```

- Create snapshot of the bound PVC

```bash
kubectl create -f ../examples/snapshot.yaml
```

- Wait the snapshot ready

```bash
$ kubectl get volumesnapshot curve-snapshot-test  -o yaml
apiVersion: snapshot.storage.k8s.io/v1beta1
kind: VolumeSnapshot
metadata:
  creationTimestamp: "2021-12-14T06:53:41Z"
  finalizers:
  - snapshot.storage.kubernetes.io/volumesnapshot-as-source-protection
  - snapshot.storage.kubernetes.io/volumesnapshot-bound-protection
  generation: 1
  name: curve-snapshot-test
  namespace: default
  resourceVersion: "319004898"
  selfLink: /apis/snapshot.storage.k8s.io/v1beta1/namespaces/default/volumesnapshots/curve-snapshot-test
  uid: 9ed2b88c-e816-438f-996e-8819980f0159
spec:
  source:
    persistentVolumeClaimName: curve-test-pvc
  volumeSnapshotClassName: curve-snapclass
status:
  boundVolumeSnapshotContentName: snapcontent-9ed2b88c-e816-438f-996e-8819980f0159
  creationTime: "2021-12-14T06:53:41Z"
  readyToUse: true
  restoreSize: 30Gi

$ kubectl get volumesnapshotcontent
NAME                                               READYTOUSE   RESTORESIZE   DELETIONPOLICY   DRIVER                  VOLUMESNAPSHOTCLASS   VOLUMESNAPSHOT        AGE
snapcontent-9ed2b88c-e816-438f-996e-8819980f0159   true         32212254720   Delete           curve.csi.netease.com   curve-snapclass       curve-snapshot-test   5m20s
```

## Restore Snapshot to a new PVC

```bash
kubectl create -f ../examples/pvc-restore.yaml
```

Get the pvc:

```bash
$ kubectl get pvc curve-pvc-restore
NAME                STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curve-pvc-restore   Bound    pvc-60cdcc88-61d1-48de-b860-0ecc0ff2dd0e   40Gi       RWO            curve          4s
```
