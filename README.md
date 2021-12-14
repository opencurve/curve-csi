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

#### CSI spec and Kubernetes version compatibility

Please refer to the [matrix](https://kubernetes-csi.github.io/docs/#kubernetes-releases)
in the Kubernetes documentation.

## Develop

You can follow the [csi spec](https://github.com/container-storage-interface/spec/blob/master/spec.md)
and [curve interface](docs/curve-interface).

## Setup

1. Deploy the `curve tool` on the CO "Master" Hosts and `curve-nbd tool` on CO "Node" Hosts by the curve cluster provider.
2. Choose a way to deploy the plugin:

- Using the helm chart: [helm installation](charts/curve-csi/README.md)
- Using the kubernetes manifests: refer to the specific version in `deploy/manifests`

## Test and User Guide

Refer to [doc](docs/README.md), you can get more details and test the driver by CSC tool.
