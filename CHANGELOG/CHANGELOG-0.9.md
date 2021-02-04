# Release notes for v0.9.1

# Changelog since v0.9.0

## Changes by Kind

### Feature

- Add "Migrated" metrics option to CSI MetricsManager, ([#77](https://github.com/kubernetes-csi/csi-lib-utils/pull/77), [@Jiawei0227](https://github.com/Jiawei0227))

### Other (Cleanup or Flake)

- Default log level of connection message reduced to '5'. ([#74](https://github.com/kubernetes-csi/csi-lib-utils/pull/74), [@steved](https://github.com/steved))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Release notes for v0.9.0

# Changelog since v0.8.1

## Changes by Kind

### Urgent Upgrade Notes

- HTTP serving logic for the metrics manager has been refactored.
  - Sidecars should create an HTTP mux (e.g. `http.ServeMux`) and pass it into `RegisterToServer()`.
  - Sidecars are responsible for starting a server with this mux. ([#70](https://github.com/kubernetes-csi/csi-lib-utils/pull/70), [@verult](https://github.com/verult))

### Feature

- Added leader election health check.
    - Sidecars should create an HTTP mux (e.g. `http.ServeMux`) and pass it into `RegisterHealthCheck()`.
    - Sidecars are responsible for starting a server with this mux.
    - A liveness probe has to be added to the pod spec for the sidecar container. ([#70](https://github.com/kubernetes-csi/csi-lib-utils/pull/70), [@verult](https://github.com/verult))

### Bug or Regression

- Workaround issue of process_start_time metric not showing up. ([#68](https://github.com/kubernetes-csi/csi-lib-utils/pull/68), [@Jiawei0227](https://github.com/Jiawei0227))
- Process_start_time should be unique in a process and therefore can now be disabled in the metrics manager registry if not needed or when it conflicts with other collectors ([#67](https://github.com/kubernetes-csi/csi-lib-utils/pull/67), [@pohly](https://github.com/pohly))

### Other (Cleanup or Flake)

- Projects using csi-lib-utils should update to klog/v2 or must ensure that klog/v1 and klog/v2 are both configured as described in https://github.com/kubernetes/klog/blob/master/examples/coexist_klog_v1_and_v2/coexist_klog_v1_and_v2.go. ([#60](https://github.com/kubernetes-csi/csi-lib-utils/pull/60), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
- k8s.io/klog: v1.0.0
