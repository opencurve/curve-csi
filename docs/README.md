# CSI Curve Driver

This document provides more detail about curve-csi driver.

- [Deploy](#deploy)
  - [Requirements](#requirements)
  - [Using the helm chart](#using-the-helm-chart)
  - [Using the kubernetes manifests](#using-the-kubernetes-manifests)
- [Debug](#debug)
- [Examples](#examples)
  - [Create StorageClass](#create-storageclass)
  - [Create PersistentVolumeClaim](#create-persistentvolumeclaim)
  - [Create Test Pod](#create-test-pod)
  - [Test block volume](#test-block-volume)
  - [Test volume expanding](#test-volume-expanding)
  - [Test snapshot](#test-snapshot)
  - [Test volume clone](#test-volume-clone)
- [Test Using CSC Tool](#test-using-csc-tool)

## Deploy

#### Requirements

The curve-csi driver deploys on the hosts contain "Master" hosts and "Node" hosts.

- On the "Master" hosts, install [curve tool](https://github.com/opencurve/curve/blob/master/docs/cn/curve%E9%80%9A%E8%BF%87tar%E5%8C%85%E6%89%8B%E5%8A%A8%E9%83%A8%E7%BD%B2sdk.md) which communicates with the curve cluster to manage the volume lifecycle, such as Create/Delete/Expand/Snapshot.
- On the "Node" hosts, install `curve-nbd tool` which allows attaching/detaching volumes to workloads.

Please refer to [deploy doc](https://github.com/opencurve/curve/blob/master/docs/cn/deploy.md) to get how to install `curve` and `curve-nbd` tool.

#### Using the helm chart

```bash
helm install --namespace "csi-system" charts/curve-csi
```

#### Using the kubernetes manifests

Change to the `deploy/manifests/` directory, create the files:

```bash
kubectl apply -f ./*.yaml
```

## Debug

You can dynamically set the log level by enabling the driver parameter `--debug-port`,
and call:

```text
curl -XPUT http://127.0.0.1:<debugPort>/debug/flags/v -d '5'
```

## Examples

#### Create StorageClass

```yaml
allowVolumeExpansion: true
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-curve-sc
parameters:
  user: k8s  # the user of curve volumes which you want to create
provisioner: curve.csi.netease.com
reclaimPolicy: Delete
volumeBindingMode: Immediate
```

#### Create PersistentVolumeClaim

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: curve-test-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi # range of 10Gi~4Ti, and step by 1Gi
  storageClassName: csi-curve-sc
```

#### Create Test Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: csi-curve-test
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
      claimName: curve-test-pvc
```

#### Test block volume

```
## create the pvc and pod:
kubectl create -f ../examples/pvc-block.yaml
kubectl create -f ../examples/pod-block.yaml

## waiting for the pod running
kubectl exec -it csi-curve-test-block bash
# mkfs.ext4 /dev/block
# mkdir -p /mnt/data && mount /dev/block /mnt/data
# cd /mnt/data
```

#### Test volume expanding

- create the normal pvc and pod
- increase pvc.spec.resources.requests.storage
- check the size of mounted path in the container

#### Test snapshot

Prerequisite: [install snapshot-controller](https://kubernetes-csi.github.io/docs/snapshot-controller.html)


Create snapshot:

```
kubectl create -f ../examples/snapshotclass.yaml
kubectl create -f ../examples/snapshot.yaml
```

#### Test volume clone

```
kubectl create -f ../examples/pvc.yaml
kubectl create -f ../examples/pvc-clone.yaml
kubectl create -f ../examples/pvc-restore.yaml
```

## Test Using CSC Tool

#### Get csc tool

Get the csc code from `https://github.com/rexray/gocsi/tree/release/1.1/csc` and run `make`.

#### Start curve csi driver

```bash
curve-csi --nodeid testnode --endpoint tcp://127.0.0.1:10000 -v 4
```

#### Get plugin info

```text
$ csc identity plugin-info --endpoint tcp://127.0.0.1:10000
"curve.csi.netease.com" "csi-v1.1.0-rc1"
```

#### Get supported capabilities

```text
$ csc identity plugin-capabilities --endpoint tcp://127.0.0.1:10000
CONTROLLER_SERVICE
ONLINE
```

#### Get controller implemented capabilities

```text
$ csc controller get-capabilities  --endpoint tcp://127.0.0.1:10000
&{type:CREATE_DELETE_VOLUME }
&{type:CLONE_VOLUME }
&{type:EXPAND_VOLUME }
```

#### Create a volume

```text
$ uuidgen
fa0c04c9-2e93-487e-8986-1e1625fd8c46

$ csc controller create --endpoint tcp://127.0.0.1:10000 \
    --req-bytes 10737418240 \
    --cap 1,mount,ext4 \
    --params user=k8s \
    volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
"0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46"  10737418240     "user"="k8s"
```

If the volume is block type, set: `--cap 5,1`

Check:

```text
$ sudo curve stat --user k8s --filename /k8s/csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
id: 39013
parentid: 39005
filetype: INODE_PAGEFILE
length(GB): 10
createtime: 2020-08-24 10:43:32
user: k8s
filename: csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
fileStatus: Created
root@pubt1-k8s-for-dev2:~#
```

#### NodeStage a volume

```text
$ sudo mkdir -p /mnt/test-csi/volume-globalmount
$ csc node stage --endpoint tcp://127.0.0.1:10000 \
   --cap 1,mount,ext4 \
   --staging-target-path /mnt/test-csi/volume-globalmount \
   0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

If the volume is block type, set: `--cap 5,1`

Check:

```text
$ sudo curve-nbd list-mapped
id    image                                                      device
97297 cbd:k8s//k8s/pvc-ce482926-91d8-11ea-bf6e-fa163e23ce53_k8s_ /dev/nbd0

$ sudo findmnt /mnt/test-csi/volume-globalmount/0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
TARGET                          SOURCE    FSTYPE OPTIONS
/mnt/test-csi/volume-globalmount/0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46 /dev/nbd0 ext4   rw,relatime,data=ordered
```

#### NodePublish a volume

```text
$ csc node publish --endpoint tcp://127.0.0.1:10000 \
    --target-path /mnt/test-csi/test-pod \
    --cap 1,mount,ext4 \
    --staging-target-path /mnt/test-csi/volume-globalmount \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

Check:

```text
$ sudo findmnt /mnt/test-csi/test-pod
TARGET                 SOURCE    FSTYPE OPTIONS
/mnt/test-csi/test-pod /dev/nbd0 ext4   rw,relatime,data=ordered
```

#### NodeGetVolumeStats a volume

```text
$ csc node stats --endpoint tcp://127.0.0.1:10000 \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46:/mnt/test-csi/test-pod
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46    /mnt/test-csi/test-pod       10447220736     10501771264     37773312      BYTES
655349  655360  11      INODES
```

Check:

```
$ sudo df /mnt/test-csi/test-pod
Filesystem     1K-blocks  Used Available Use% Mounted on
/dev/nbd0      10255636 36888  10202364   1% /mnt/test-csi/test-pod

$ sudo df -i /mnt/test-csi/test-pod
Filesystem     Inodes IUsed  IFree IUse% Mounted on
/dev/nbd0      655360    11 655349    1% /mnt/test-csi/test-pod
```

#### Expand a Volume

```text
$ # controllerExpand:
$ csc controller expand --endpoint tcp://127.0.0.1:10000 \
    --req-bytes 21474836480 \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
21474836480

$ # nodeExpand:
$ csc node expand --endpoint tcp://127.0.0.1:10000 \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46 \
    /mnt/test-csi/test-pod
0
```

Check:

```text
$ sudo curve stat --user k8s --filename /k8s/csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46|grep length
length(GB): 20

$ sudo df -h /mnt/test-csi/test-pod
/dev/nbd0        20G   44M   20G   1% /mnt/test-csi/test-pod
```

#### NodeUnpublish a volume

```text
$ csc node unpublish --endpoint tcp://127.0.0.1:10000 \
    --target-path /mnt/test-csi/test-pod \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

Check:

```
$ sudo file /mnt/test-csi/test-pod
/mnt/test-csi/test-pod: cannot open `/mnt/test-csi/test-pod' (No such file or directory)
```

#### NodeUnstage a volume

```text
$ csc node unstage --endpoint tcp://127.0.0.1:10000 \
    --staging-target-path /mnt/test-csi/volume-globalmount \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

Check:

```
$ sudo file /mnt/test-csi/volume-globalmount
/mnt/test-csi/volume-globalmount: cannot open `/mnt/test-csi/volume-globalmount' (No such file or directory)

$ sudo curve-nbd list-mapped
```

#### Delete a volume

```text
$ csc controller delete --endpoint tcp://127.0.0.1:10000 \
  0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

Check:

```
$ sudo curve stat --user k8s --filename /k8s/csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
E 2020-08-24T13:57:55.946636+0800 58360 mds_client.cpp:395] GetFileInfo: filename = /k8s/volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46, owner = k8s, errocde = 6, error msg = kFileNotExists, log id = 1
stat fail, ret = -6
```

#### Snapshot

```
## create a snapshot
$ csc controller create-snapshot --endpoint tcp://127.0.0.1:10000 \
	--source-volume 0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46 \
	snapshot-215d24ff-c04c-4b08-a1fb-692c94627c63
"0024-9ea1a8fc-160d-47ef-b2ef-f0e09677b066-0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46"	23622320128	0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46	seconds:1639297566 nanos:780828000 	true

## delete a snapshot
$ csc controller delete-snapshot --endpoint tcp://127.0.0.1:10000 \
    0024-9ea1a8fc-160d-47ef-b2ef-f0e09677b066-0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0024-9ea1a8fc-160d-47ef-b2ef-f0e09677b066-0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```
