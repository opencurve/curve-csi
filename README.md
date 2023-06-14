# Curve CSI Driver

[![Go Report Card](https://goreportcard.com/badge/github.com/opencurve/curve-csi)](https://goreportcard.com/report/github.com/opencurve/curve-csi)

## Overview

The plugin implements the Container Storage Interface(CSI) between
Container Orchestrator(CO) and Curve cluster. It allows dynamically
provisioning curve volumes and attaching them to workloads.

Refer to [https://github.com/opencurve/curve](https://github.com/opencurve/curve) for the Curve details.

## Supported version

The driver is currently developed with csi spec v1.5.0, and supported kubernetes v1.17+.

Other csi-v1.0+ enabled container orchestrator environments may work fine.

### CSI spec and Kubernetes version compatibility

Please refer to the [matrix](https://kubernetes-csi.github.io/docs/#kubernetes-releases)
in the Kubernetes documentation.

## Develop

You can follow the [csi spec](https://github.com/container-storage-interface/spec/blob/master/spec.md)
and [curve interface](docs/curve-interface).

## Setup

Choose a way to deploy the plugin:

1. Using the kubernetes manifests: refer to [deploy doc](https://github.com/opencurve/curve-csi/blob/0ecb1fd4d47819c49acf1f7f92a53ab5ac83c514/docs/README.md?plain=1#L20)
2. Using the helm chart: [helm installation](charts/curve-csi/README.md)

## Test and User Guide

Refer to [doc](docs/README.md), you can get more details and test the driver by CSC tool.
