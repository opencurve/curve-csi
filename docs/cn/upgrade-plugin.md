[English version](../upgrade-plugin.md)

# 插件升级

本文档介绍如何升级插件。

首先，插件的镜像不包含curve工具（插件通过调用宿主机的curve工具或http接口实现）。如果你想升级curve客户端，需要参考[curve文档](https://github.com/opencurve/curve/tree/master/docs)，并且需确保接口和[curve接口](./curve-interface)一致。


## 从v2.0.x升级到v2.1.x

插件从v2.1开始支持CurveFS，且默认的块存储驱动名称改为由`curve.csi.netease.com`改为`curvebs.csi.netease.com`。

有2种升级选择：

- 同时保留老的命名为`curve.csi.netease.com`的插件和新的命名为`curvebs.csi.netease.com`的插件，并逐步下线老的pv。
- 使用`--drivername=curve.csi.netease.com`参数升级curve块存储插件，新的curvebs类型卷仍使用老的命名。

## 从csi-v1.1.0-xx升级到v2.0.x

此版本有一个不兼容的改变：在node上stage卷阶段，路径由`req.GetStagingTargetPath()`改为`req.GetStagingTargetPath() + "/" + volumeId`, 所以不支持原地升级。

升级步骤：

- 升级provisioner deployment。
- 更新node-plugin daemonset，且将其更新策略改为`OnDelete`。
- 逐步升级节点：
    - 节点置为不可调度。
    - 驱逐节点上使用curve的pods。
    - 升级该节点的node-plugin pod。

