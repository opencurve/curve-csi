[English version](../clone.md)

# 克隆

- [CurveBS](#curvebs)
  - <a href="#bsverify">确保pvc已绑定</a>
  - <a href="#bscreate">创建新的克隆pvc</a>
  - <a href="#bsnew">查看新的pvc</a>
- [CurveFS](#curvefs)

## CurveBS

### <div id="bsverify">确保pvc已绑定</div>

```bash
$ kubectl get pvc
NAME                     STATUS        VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc         Bound         pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   30Gi       RWO            curvebs          35m
```

### <div id="bscreate">创建新的克隆pvc</div>

```
kubectl create -f ../examples/curvebs/pvc-clone.yaml
```

### <div id="bsnew">查看新的pvc</div>

```bash
$ kubectl get pvc curvebs-pvc-clone
NAME                STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-pvc-clone   Bound    pvc-755cfb47-9b03-41b5-bdf9-0772a1ae41ef   40Gi       RWO            curvebs          3s
```

## CurveFS

