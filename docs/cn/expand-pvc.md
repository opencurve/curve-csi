[English version](../expand-pvc.md)

# PVC扩容

- <a href="#req">依赖</a>
- [CurveBS](#curvebs)
  - <a href="#bsfs">扩容curvebs文件存储pvc</a>
  - <a href="#bsblock">扩容curvebs块存储pvc</a>
- [CurveFS](#curvefs)

## <div id="req">依赖</div>

- 对于文件系统扩容的支持，kubernetes版本需要>=v1.15；对于块设备模式的扩容，kubernetes版本需要>=v1.16。并且，`ExpandCSIVolumes`特性门控需要被打开。
- PVC使用的StorageClass的`allowVolumeExpansion`需设为`true`。

## CurveBS

### <div id="bsfs">扩容curvebs文件存储pvc</div>

创建PVC:

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

等待PVC绑定:

```bash
$ kubectl get pvc 
NAME               STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc   Bound    pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   20Gi       RWO            curvebs          38s
```

创建使用这个PVC的Pod, 并校验大小:

```bash
$ kubectl exec -it csi-curvebs-test bash
root@csi-curvebs-test:/# df -h /var/lib/www/html
Filesystem      Size  Used Avail Use% Mounted on
/dev/nbd0        20G   45M   20G   1% /var/lib/www/html
```

编辑PVC增加其容量(pvc.spec.resource.requests.storage)，然后检查：

```bash
$  kubectl get pvc curvebs-test-pvc 
NAME               STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curvebs-test-pvc   Bound    pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   30Gi       RWO            curvebs          7m3s
```

校验容器内的大小：

```bash
$ kubectl exec -it csi-curvebs-test bash
root@csi-curvebs-test:/# df -h /var/lib/www/html
Filesystem      Size  Used Avail Use% Mounted on
/dev/nbd0        30G   44M   30G   1% /var/lib/www/html
```

### <div id="bsblock">扩容curvebs块存储pvc</div>

创建block模式PVC:

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

等待PVC绑定:

```bash
$ kubectl get pvc  curvebs-test-pvc-block
NAME                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc-block   Bound    pvc-6fd55b8f-5a26-422c-b4d9-e9613e5724b5   20Gi       RWX            curvebs          14s
```

创建使用这个PVC的Pod:

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

校验大小:

```bash
$ kubectl exec -it csi-curvebs-test-block bash
root@csi-curvebs-test-block:/# blockdev  --getsize64 /dev/block 
21474836480
```

编辑PVC增加其容量(pvc.spec.resource.requests.storage)，然后检查：

```bash
$ kubectl get pvc curvebs-test-pvc-block
NAME                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc-block   Bound    pvc-6fd55b8f-5a26-422c-b4d9-e9613e5724b5   30Gi       RWX            curvebs          6m45s
```

校验容器内的大小：

```bash
$ kubectl exec -it csi-curvebs-test-block bash
root@csi-curvebs-test-block:/# blockdev  --getsize64 /dev/block 
32212254720
```

## CurveFS
