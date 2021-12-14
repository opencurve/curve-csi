- [Prerequisite](#prerequisite)
- [Expand Curve Filesystem PVC](#expand-curve-filesystem-pvc)
- [Expand Curve Block PVC](#expand-curve-block-pvc)

## Prerequisite

- For filesystem expansion to be supported for your kubernetes cluster, the kubernetes version running in your cluster should be >= v1.15 and for block volume expand support the kubernetes version should be >=1.16. Also, `ExpandCSIVolumes` feature gate has to be enabled for the volume expand functionality to work.

- The controlling StorageClass must have `allowVolumeExpansion` set to `true`.

## Expand Curve Filesystem PVC

Create PVC:

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
      storage: 20Gi
  storageClassName: curve
```

Wait PVC bounded:

```bash
$ kubectl get pvc 
NAME             STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curve-test-pvc   Bound    pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   20Gi       RWO            curve          38s
```

Create the pod using this PVC, check the size:

```bash
$ kubectl exec -it csi-curve-test bash
root@csi-curve-test:/# df -h /var/lib/www/html
Filesystem      Size  Used Avail Use% Mounted on
/dev/nbd0        20G   45M   20G   1% /var/lib/www/html
```

Now expand the PVC by editing the PVC (pvc.spec.resource.requests.storage), then get it:

```bash
$  kubectl get pvc curve-test-pvc 
NAME             STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curve-test-pvc   Bound    pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   30Gi       RWO            curve          7m3s
```

Check the directory size inside the pod where PVC is mounted:

```bash
$ kubectl exec -it csi-curve-test bash
root@csi-curve-test:/# df -h /var/lib/www/html
Filesystem      Size  Used Avail Use% Mounted on
/dev/nbd0        30G   44M   30G   1% /var/lib/www/html
```

## Expand Curve Block PVC

Create block PVC:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: curve-test-pvc-block
spec:
  accessModes:
  - ReadWriteMany
  volumeMode: Block
  storageClassName: curve
  resources:
    requests:
      storage: 20Gi
```

Wait the PVC bounded:

```bash
$ kubectl get pvc  curve-test-pvc-block
NAME                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curve-test-pvc-block   Bound    pvc-6fd55b8f-5a26-422c-b4d9-e9613e5724b5   20Gi       RWX            curve          14s
```

Create the pod using this PVC:

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    container.apparmor.security.beta.kubernetes.io/my-container: unconfined
  name: csi-curve-test-block
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
      claimName: curve-test-pvc-block
```

Check the size:

```bash
$ kubectl exec -it csi-curve-test-block bash
root@csi-curve-test-block:/# blockdev  --getsize64 /dev/block 
21474836480
```

Now expand the PVC by editing the PVC (pvc.spec.resource.requests.storage), then get it:

```bash
$ kubectl get pvc curve-test-pvc-block
NAME                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curve-test-pvc-block   Bound    pvc-6fd55b8f-5a26-422c-b4d9-e9613e5724b5   30Gi       RWX            curve          6m45s
```

Check the block size inside the pod:

```bash
$ kubectl exec -it csi-curve-test-block bash
root@csi-curve-test-block:/# blockdev  --getsize64 /dev/block 
32212254720
```
