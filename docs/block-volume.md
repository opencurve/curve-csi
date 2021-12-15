# CurveBS block volume

- [Create PVC](#create-pvc)
- [Create Pod](#create-pod)
- [Check in Pod](#check-in-pod)

### Create PVC

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

### Create Pod

```
apiVersion: v1
kind: Pod
metadata:
  annotations:
    # https://kubernetes.io/docs/tutorials/clusters/apparmor/
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

### Check in Pod

```
## waiting for the pod running
kubectl exec -it csi-curvebs-test-block bash
# mkfs.ext4 /dev/block
# mkdir -p /mnt/data && mount /dev/block /mnt/data
# cd /mnt/data
```
