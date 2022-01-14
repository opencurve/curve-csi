[English version](README.md)

# Curve CSI Driver

[![Go Report Card](https://goreportcard.com/badge/github.com/opencurve/curve-csi)](https://goreportcard.com/report/github.com/opencurve/curve-csi)


## 简介

插件在容器编排系统(CO)和Curve服务之间实现了Container Storage Interface(CSI)。它可以动态提供curve块/文件存储卷，并自动挂载到业务负载中。

关于Curve存储系统，可参考[https://github.com/opencurve/curve](https://github.com/opencurve/curve)。

## 版本支持

插件基于csi spec v1.5.0实现，可支持kubernetes v1.17及以上版本。

当然其他容器编排系统，如果支持csi v1.0+，也可能兼容。

### CSI版本和Kubernetes版本的兼容

参考：[matrix](https://kubernetes-csi.github.io/docs/#kubernetes-releases)。

### 插件发布版本

| 分支 | 最新版本 |CSI版本 | Kubernetes 版本 | 功能 |
|--- | ---| --- |--- | ---|
| master/release-2.1 | v2.1.0 | v1.5.0 | v1.17+ | - Support CurveFS |
| release-2.0 | v2.0.0 | v1.5.0 | v1.17+ | - Snapshot<br/> - Clone<br/> - Block mode volume|
| release-csi-1.1 | csi-v1.1.0-rc2 | v1.1.0 | v1.13+ | - Dynamically provision <br/> - Expand volume <br/> - Volume metrics|

不同的发布版本可能有兼容性问题，升级插件时可参考文档：[upgrade-plugin](docs/cn/upgrade-plugin.md)

## 开发

可参考csi接口[csi spec](https://github.com/container-storage-interface/spec/blob/master/spec.md)和curve接口[curve interface](docs/cn/curve-interface)。

## 部署

1. 通过curve集群管理员，在Master机器上部署`curve`工具，在node节点上部署`curve-nbd`工具。
2. 选择一种方式部署插件：

- 使用helm: [helm installation](charts/curve-csi/README.md)。
- 使用kubernetes编排文件: 参考`deploy/`。

## 文档

更多详情参考[doc](docs/cn/README.md)。
