# Expand PVC

- [Prerequisite](#prerequisite)
- [CurveBS](#curvebs)
  - [Expand CurveBS Filesystem PVC](#expand-curvebs-filesystem-pvc)
  - [Expand CurveBS Block PVC](#expand-curvebs-block-pvc)
- [CurveFS](#curvefs)

## Prerequisite

- For filesystem expansion to be supported for your kubernetes cluster, the kubernetes version running in your cluster should be >= v1.15 and for block volume expand support the kubernetes version should be >=1.16. Also, `ExpandCSIVolumes` feature gate has to be enabled for the volume expand functionality to work.

- The controlling StorageClass must have `allowVolumeExpansion` set to `true`.

## CurveBS

### Expand CurveBS Filesystem PVC

Create PVC:

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
      storage: 20Gi
  storageClassName: curvebs
```

Wait PVC bounded:

```bash
$ kubectl get pvc 
NAME               STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc   Bound    pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   20Gi       RWO            curvebs          38s
```

Create the pod using this PVC, check the size:

```bash
$ kubectl exec -it csi-curvebs-test bash
root@csi-curvebs-test:/# df -h /var/lib/www/html
Filesystem      Size  Used Avail Use% Mounted on
/dev/nbd0        20G   45M   20G   1% /var/lib/www/html
```

Now expand the PVC by editing the PVC (pvc.spec.resource.requests.storage), then get it:

```bash
$  kubectl get pvc curvebs-test-pvc 
NAME               STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curvebs-test-pvc   Bound    pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   30Gi       RWO            curvebs          7m3s
```

Check the directory size inside the pod where PVC is mounted:

```bash
$ kubectl exec -it csi-curvebs-test bash
root@csi-curvebs-test:/# df -h /var/lib/www/html
Filesystem      Size  Used Avail Use% Mounted on
/dev/nbd0        30G   44M   30G   1% /var/lib/www/html
```

### Expand CurveBS Block PVC

Create block PVC:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: curvebs-test-pvc-block
spec:
  accessModes:
  - ReadWriteMany
  volumeMode: Block
  storageClassName: curvebs
  resources:
    requests:
      storage: 20Gi
```

Wait the PVC bounded:

```bash
$ kubectl get pvc  curvebs-test-pvc-block
NAME                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc-block   Bound    pvc-6fd55b8f-5a26-422c-b4d9-e9613e5724b5   20Gi       RWX            curvebs          14s
```

Create the pod using this PVC:

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    container.apparmor.security.beta.kubernetes.io/my-container: unconfined
  name: csi-curvebs-test-block
spec:
  containers:
  - name: my-container
    image: debian
    command:
    - sleep
    - "3600"
    securityContext:
      capabilities:
        add: ["SYS_ADMIN"]
    volumeDevices:
    - devicePath: /dev/block
      name: my-volume
  volumes:
  - name: my-volume
    persistentVolumeClaim:
      claimName: curvebs-test-pvc-block
```

Check the size:

```bash
$ kubectl exec -it csi-curvebs-test-block bash
root@csi-curvebs-test-block:/# blockdev  --getsize64 /dev/block 
21474836480
```

Now expand the PVC by editing the PVC (pvc.spec.resource.requests.storage), then get it:

```bash
$ kubectl get pvc curvebs-test-pvc-block
NAME                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc-block   Bound    pvc-6fd55b8f-5a26-422c-b4d9-e9613e5724b5   30Gi       RWX            curvebs          6m45s
```

Check the block size inside the pod:

```bash
$ kubectl exec -it csi-curvebs-test-block bash
root@csi-curvebs-test-block:/# blockdev  --getsize64 /dev/block 
32212254720
```

## CurveFS
