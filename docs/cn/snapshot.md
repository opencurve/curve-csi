[English version](../snapshot.md)

# 快照

- <a href="#req">依赖</a>
- [CurveBS](#curvebs)
  - <a href="#bscreatesc">创建SnapshotClass</a>
  - <a href="#bscreatesnap">创建快照</a> 
  - <a href="#bsrestore">从快照恢复</a>
- [CurveFS](#curvefs)

## <div id="req">依赖</div>

对于快照功能的支持，kubernetes版本需要`>= v1.17`。并且需要部署snapshot控制器，csi控制器需要部署csi-snapshotter sidecar容器。

**Git地址:**  https://github.com/kubernetes-csi/external-snapshotter

**版本支持**

|最新发布	|最小CSI版本	|最大CSI版本	|镜像	|最小K8s版本	|最大K8s版本	|建议K8s版本|
| ---	| --- 	| ---	| ---	| --- |---	|---|
|v4.2.1	|	v1.0.0	|-	|k8s.gcr.io/sig-storage/snapshot-controller:v4.2.1|	v1.20	|-	|v1.22|
| v4.1.1|	v1.0.0	|-	|k8s.gcr.io/sig-storage/snapshot-controller:v4.1.1|	v1.20|	-	|v1.20|
| v3.0.3 (beta)|	v1.0.0|	-|	k8s.gcr.io/sig-storage/snapshot-controller:v3.0.3	|v1.17	|-|	v1.17|


**安装Snapshot Beta CRDs:**

```bash
kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml

kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml

kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml
```

**安装Snapshot控制器:**

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml

kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-3.0/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml
```

## CurveBS

### <div id="bscreatesc">创建SnapshotClass</div>

```bash
kubectl create -f ../examples/curvebs/snapshotclass.yaml
```

### <div id="bscreatesnap">创建快照</div> 

- 确认PVC绑定：

```bash
$ kubectl get pvc
NAME                     STATUS        VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc         Bound         pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   30Gi       RWO            curvebs          35m
```

- 创建此PVC的快照：

```bash
kubectl create -f ../examples/curvebs/snapshot.yaml
```

- 等待快照完成：

```bash
$ kubectl get volumesnapshot curvebs-snapshot-test  -o yaml
apiVersion: snapshot.storage.k8s.io/v1beta1
kind: VolumeSnapshot
metadata:
  creationTimestamp: "2021-12-14T06:53:41Z"
  finalizers:
  - snapshot.storage.kubernetes.io/volumesnapshot-as-source-protection
  - snapshot.storage.kubernetes.io/volumesnapshot-bound-protection
  generation: 1
  name: curvebs-snapshot-test
  namespace: default
  resourceVersion: "319004898"
  selfLink: /apis/snapshot.storage.k8s.io/v1beta1/namespaces/default/volumesnapshots/curvebs-snapshot-test
  uid: 9ed2b88c-e816-438f-996e-8819980f0159
spec:
  source:
    persistentVolumeClaimName: curvebs-test-pvc
  volumeSnapshotClassName: curvebs-snapclass
status:
  boundVolumeSnapshotContentName: snapcontent-9ed2b88c-e816-438f-996e-8819980f0159
  creationTime: "2021-12-14T06:53:41Z"
  readyToUse: true
  restoreSize: 30Gi

$ kubectl get volumesnapshotcontent
NAME                                               READYTOUSE   RESTORESIZE   DELETIONPOLICY   DRIVER                    VOLUMESNAPSHOTCLASS     VOLUMESNAPSHOT          AGE
snapcontent-9ed2b88c-e816-438f-996e-8819980f0159   true         32212254720   Delete           curvebs.csi.netease.com   curvebs-snapclass       curvebs-snapshot-test   5m20s
```

### <div id="bsrestore">从快照恢复</div> 

```bash
kubectl create -f ../examples/curvebs/pvc-restore.yaml
```

查看新的PVC:

```bash
$ kubectl get pvc curvebs-pvc-restore
NAME                  STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-pvc-restore   Bound    pvc-60cdcc88-61d1-48de-b860-0ecc0ff2dd0e   40Gi       RWO            curvebs          4s
```

## CurveFS

