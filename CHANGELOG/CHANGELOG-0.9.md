# Release notes for v0.9.0

[Documentation](https://docs.k8s.io/docs/home)
# Changelog since v0.8.1

## Changes by Kind

### Other (Cleanup or Flake)

- Projects using csi-lib-utils should update to klog/v2 or must ensure that klog/v1 and klog/v2 are both configured as described in https://github.com/kubernetes/klog/blob/master/examples/coexist_klog_v1_and_v2/coexist_klog_v1_and_v2.go. ([#60](https://github.com/kubernetes-csi/csi-lib-utils/pull/60), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
- k8s.io/klog: v1.0.0
