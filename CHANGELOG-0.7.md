# Changelog since v0.6.1

## Breaking Changes

- The rpc and connection package had duplicate code removed from connection package. [#31](https://github.com/kubernetes-csi/csi-lib-utils/pull/31), [@Madhu-1](https://github.com/Madhu-1))
- The connection package requires the new `metrics.CSIMetricsManager`.


## New Features

- Introduce a CSI Metrics package (github.com/kubernetes-csi/csi-lib-utils/metrics) that can be used to automatically record prometheus metrics for every CSI gRPC call. ([#35](https://github.com/kubernetes-csi/csi-lib-utils/pull/35), [@saad-ali](https://github.com/saad-ali))


## Other Notable Changes

- Switch to vendoring and dependencies managed with "go mod" ([#33](https://github.com/kubernetes-csi/csi-lib-utils/pull/33), [@pohly](https://github.com/pohly))
- Kubernetes dependencies update to v1.17.0

