[English version](../README.md)

# 概述

此文档介绍curve-csi部署/设计/测试/开发等细节。

- <a href="#deploybs">部署CurveBS插件</a>
  - <a href="#bsreq">依赖</a>
  - <a href="#bshelm">使用helm部署</a>
  - <a href="#bsk8s">使用kubernetes编排文件</a>
- <a href="#deployfs">部署CurveFS插件</a>
- [Debug](#debug)
- <a href="#eg">例子</a>
  - [CurveBS](#curvebs)
    - <a href="#bssc">创建StorageClass</a>
    - <a href="#bspvc">创建PersistentVolumeClaim</a>
    - <a href="#bspod">创建测试Pod</a>
  - [CurveFS](#curvefs)
  - <a href="#block">Block模式卷</a>
  - <a href="#expand">卷扩容</a>
  - <a href="#snap">快照</a>
  - <a href="#clone">克隆</a>
- <a href="#test">测试</a>
  - <a href="#csctest">CSC工具测试</a>
  - [E2E](#e2e)


## <div id="deploybs">部署CurveBS插件</div>

#### <div id="bsreq">依赖</div>

curve-csi插件是以容器部署在Master和Node机器上。

- 在Master机器上，需安装`curve`工具，来负责管理卷的生命周期，比如创建/删除/扩容/快照。
- 在Node机器上，需安装`curve-nbd`工具，负责挂载/卸载卷。

请参考文档[deploy doc](https://github.com/opencurve/curve/blob/master/docs/cn/deploy.md) 来安装`curve`和`curve-nbd`工具。

#### <div id="bshelm">使用helm部署</div>

```bash
helm install --namespace "csi-system" charts/curve-csi-curvebs
```

#### <div id="bsk8s">使用kubernetes编排文件</div>

进到`deploy/curve-csi-curvebs`目录，创建：

```bash
kubectl apply -f ./*.yaml
```

## <div id="deployfs">部署CurveFS插件</div>

## Debug

排障时可以动态调整日志级别：

```bash
curl -XPUT http://127.0.0.1:<debugPort>/debug/flags/v -d '5'
```

## <div id="eg">例子</div>

### CurveBS

#### <div id="bssc">创建StorageClass</div>

```yaml
allowVolumeExpansion: true
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: curvebs
parameters:
  user: k8s  # 卷的用户名，通过此字段进行租户隔离
  cloneLazy: "true"
provisioner: curvebs.csi.netease.com
reclaimPolicy: Delete
volumeBindingMode: Immediate
```

#### <div id="bspvc">创建PersistentVolumeClaim</div>

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
      storage: 10Gi # 范围在10Gi～4Ti，且步长为1Gi
  storageClassName: curvebs
```

#### <div id="bspod">创建测试Pod</div>

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

### <div id="block">Block模式卷</div>

参考[block-volume](./block-volume.md)

### <div id="expand">卷扩容</div>

参考[expand-pvc](./expand-pvc.md)

### <div id="snap">快照</div>

参考[snapshot](./snapshot.md)

### <div id="clone">克隆</div>

参考[clone](./clone.md)

## <div id="test">测试</div>

### <div id="csctest">CSC工具测试</div>

参考[csc-test](./csc-test.md)

### E2E

参考[e2e](./e2e.md)

