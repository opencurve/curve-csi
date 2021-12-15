# Curve CSI

The curve-csi-curvebs chart adds the curvebs(Curve Block Storage) support to your k8s cluster.

# Installation

Install the chart to your kubernetes cluster:

```bash
helm install --namespace "curve-csi-system" ../curve-csi-curvebs
```

After installation succeeds, you can get a status of Chart:

```bash
helm status "curve-csi-curvebs"
```

If you want to delete your Chart, use this command:

```bash
helm delete --purge "curve-csi-curvebs"
```

If you want to delete the namespace, use this command:

```bash
kubectl delete namespace curve-csi-system
```
