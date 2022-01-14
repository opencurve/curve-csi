[English version](../block-volume.md)

# CurveBS block volume

- <a href="#pvc">创建pvc</a>
- <a href="#pod">创建pod</a>
- <a href="#check">pod内检查</a>

### <div id="pvc">创建pvc</div>

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

### <div id="pod">创建pod</div>

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

### <div id="check">pod内检查</div>

```
## waiting for the pod running
kubectl exec -it csi-curvebs-test-block bash
# mkfs.ext4 /dev/block
# mkdir -p /mnt/data && mount /dev/block /mnt/data
# cd /mnt/data
```
