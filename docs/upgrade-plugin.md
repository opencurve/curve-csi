[中文版](cn/upgrade-plugin.md)

# Upgrade plugin

This page explains how to upgrade a plugin.

First of all, the plugin image does not contain any curve-tools(It implements by invoking the http interface or tools in the host). If you want to update the curve client, refer to the [curve doc](https://github.com/opencurve/curve/tree/master/docs), and ensure the interface consistent with [curve-interface](./curve-interface).

## Upgrading from v2.0.x to v2.1.x

Plugin v2.1 adds CurveFS support, and the default name of block storage is changed to `curvebs.csi.netease.com`.

There is 2 choices to support v2.1:

- Keep the old plugin named `curve.csi.netease.com` and new plugin named `curvebs.csi.netease.com` existing, offline the old persistent volumes gradually.
- Deploy the curvebs plugin with specific flag `--drivername=curve.csi.netease.com` to continue to use the old driver name.

## Upgrading from csi-v1.1.0-xx to v2.0.x

The breaking change: when node stages volume, stagingTargetPath is `req.GetStagingTargetPath() + "/" + volumeId` in v2.0.x and `req.GetStagingTargetPath()` in csi-v1.1.0-xx. So it does not support upgrading in place.

The upgrade workflow is the following:

- Apply provisioner deployment.
- Update node-plugin daemonset and set its updateStrategy to `OnDelete`.
- Cordon node and evict pods using curve volumes, then update curve node-plugin pod.
