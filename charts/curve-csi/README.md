# Curve CSI

The curve-csi chart adds the curve volume support to your k8s cluster.

## Installation

### Prerequisite

* For single node k8s cluster, modify replicas in [values.yaml](https://github.com/opencurve/curve-csi/blob/0ecb1fd4d47819c49acf1f7f92a53ab5ac83c514/charts/curve-csi/values.yaml#L29) to 1.

* Modify the env `MDSADDR` at `csi-deployment.yaml` and `csi-daemonset.yaml` as backend cluster addr.

* Modify the `--snapshot-server` startup parameter at [csi-deployment.yaml](https://github.com/opencurve/curve-csi/blob/0ecb1fd4d47819c49acf1f7f92a53ab5ac83c514/charts/curve-csi/templates/csi-deployment.yaml#L119)
  * Delete it if don't need the snapshot feature.
  * Modify it to the correct backend curvebs snapshotcloneserver address and refer the [docs snapshot](https://github.com/opencurve/curve-csi/blob/master/docs/snapshot.md) to install other components.

### v3.0.0

#### Install

Install the chart to your kubernetes cluster:

```bash
helm install curve-csi --namespace "curve-csi-system" ./curve-csi
```

### v2.0.0

#### Requirements

The curve-csi driver deploys on the hosts contain "Master" hosts and "Node" hosts.

1. Install tools

    Please refer to [deploy doc](https://github.com/opencurve/curve/blob/master/docs/cn/deploy.md) to get how to install `curve` and `curve-nbd` tool.

    * On the "Master" hosts, install `curve tool` which communicates with the curve cluster to manage the volume lifecycle, such as Create/Delete/Expand/Snapshot.
    * On the "Node" hosts, install `curve-nbd tool` which allows attaching/detaching volumes to workloads.

2. Beforehand

    Change to v2.0.0 image(`curvecsi/curvecsi:v2.0.0`) at [provisioner-deploy.yaml](https://github.com/opencurve/curve-csi/blob/0ecb1fd4d47819c49acf1f7f92a53ab5ac83c514/deploy/manifests/provisioner-deploy.yaml#LL103C10-L103C10) and [node-plugin-daemonset.yaml](https://github.com/opencurve/curve-csi/blob/0ecb1fd4d47819c49acf1f7f92a53ab5ac83c514/deploy/manifests/node-plugin-daemonset.yaml#L47)

#### Install

Install the chart to your kubernetes cluster:

```bash
helm install curve-csi --namespace curve-csi-system ./curve-csi
```

### Stauts

After installation succeeds, you can get a status of Chart:

```bash
helm status curve-csi -n curve-csi-system
```

### Delete

If you want to delete your Chart, use this command:

```bash
helm delete --purge curve-csi -n curve-csi-system
```

### Delete namespace

If you want to delete the namespace, use this command:

```bash
kubectl delete namespace curve-csi-system
```
