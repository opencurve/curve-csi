### Verify if PVC is in Bound state

```bash
$ kubectl get pvc
NAME                   STATUS        VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curve-test-pvc         Bound         pvc-b2789e09-9854-4aa1-b556-d9b0e0569f87   30Gi       RWO            curve          35m
```

### Create a new cloned PVC

```
kubectl create -f ../examples/pvc-clone.yaml
```

### Get PVC

```bash
$ kubectl get pvc curve-pvc-clone
NAME              STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
curve-pvc-clone   Bound    pvc-755cfb47-9b03-41b5-bdf9-0772a1ae41ef   40Gi       RWO            curve          3s
```
