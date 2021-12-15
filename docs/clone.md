# CSI Clone doc

- [CurveBS](#curvebs)
  - [Verify if PVC is in Bound state](#verify-if-pvc-is-in-bound-state)
  - [Create a new cloned PVC](#create-a-new-cloned-pvc)
  - [Get PVC](#get-pvc)
- [CurveFS](#curvefs)

## CurveBS

### Verify if PVC is in Bound state

```bash
$ kubectl get pvc
NAME                     STATUS        VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-test-pvc         Bound         pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   30Gi       RWO            curvebs          35m
```

### Create a new cloned PVC

```
kubectl create -f ../examples/curvebs/pvc-clone.yaml
```

### Get PVC

```bash
$ kubectl get pvc curvebs-pvc-clone
NAME                STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS     AGE
curvebs-pvc-clone   Bound    pvc-755cfb47-9b03-41b5-bdf9-0772a1ae41ef   40Gi       RWO            curvebs          3s
```

## CurveFS

