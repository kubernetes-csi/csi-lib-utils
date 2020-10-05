# Release notes for v0.7.1

# Changelog since v0.7.0

## Changes by Kind

### Bug or Regression

- Add `process_start_time_seconds` into the csi metric lib. ([#56](https://github.com/kubernetes-csi/csi-lib-utils/pull/56), [@Jiawei0227](https://github.com/Jiawei0227))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Release notes for v0.7.0

# Changelog since v0.6.1

## Breaking Changes

- The rpc and connection package had duplicate code removed from connection package. [#31](https://github.com/kubernetes-csi/csi-lib-utils/pull/31), [@Madhu-1](https://github.com/Madhu-1))
- The connection package requires the new `metrics.CSIMetricsManager`.


## New Features

- Introduce a CSI Metrics package (github.com/kubernetes-csi/csi-lib-utils/metrics) that can be used to automatically record prometheus metrics for every CSI gRPC call. ([#35](https://github.com/kubernetes-csi/csi-lib-utils/pull/35), [@saad-ali](https://github.com/saad-ali))


## Other Notable Changes

- Switch to vendoring and dependencies managed with "go mod" ([#33](https://github.com/kubernetes-csi/csi-lib-utils/pull/33), [@pohly](https://github.com/pohly))
- Kubernetes dependencies update to v1.17.0

