# Release notes for v0.18.0

# Changelog since v0.17.0

## Changes by Kind

### Feature

- Added support for structured logging (the log messages have been changed due to the activation of structured logging) ([#149](https://github.com/kubernetes-csi/csi-lib-utils/pull/149), [@bells17](https://github.com/bells17))

### Other (Cleanup or Flake)

- Updates Kubernetes dependencies to 1.30. ([#165](https://github.com/kubernetes-csi/csi-lib-utils/pull/165), [@xing-yang](https://github.com/xing-yang))

## Dependencies

### Added
- github.com/fxamacker/cbor/v2: [v2.6.0](https://github.com/fxamacker/cbor/v2/tree/v2.6.0)
- github.com/x448/float16: [v0.8.4](https://github.com/x448/float16/tree/v0.8.4)
- go.uber.org/goleak: v1.3.0
- k8s.io/gengo/v2: 51d4e06

### Changed
- github.com/go-logr/logr: [v1.3.0 → v1.4.1](https://github.com/go-logr/logr/compare/v1.3.0...v1.4.1)
- github.com/go-logr/zapr: [v1.2.3 → v1.3.0](https://github.com/go-logr/zapr/compare/v1.2.3...v1.3.0)
- github.com/golang/protobuf: [v1.5.3 → v1.5.4](https://github.com/golang/protobuf/compare/v1.5.3...v1.5.4)
- github.com/onsi/ginkgo/v2: [v2.13.0 → v2.15.0](https://github.com/onsi/ginkgo/v2/compare/v2.13.0...v2.15.0)
- github.com/onsi/gomega: [v1.29.0 → v1.31.0](https://github.com/onsi/gomega/compare/v1.29.0...v1.31.0)
- go.uber.org/zap: v1.19.0 → v1.26.0
- golang.org/x/crypto: v0.15.0 → v0.21.0
- golang.org/x/mod: v0.8.0 → v0.15.0
- golang.org/x/net: v0.18.0 → v0.23.0
- golang.org/x/sys: v0.14.0 → v0.18.0
- golang.org/x/term: v0.14.0 → v0.18.0
- golang.org/x/tools: v0.12.0 → v0.18.0
- google.golang.org/protobuf: v1.31.0 → v1.33.0
- k8s.io/api: v0.29.0 → v0.30.0
- k8s.io/apimachinery: v0.29.0 → v0.30.0
- k8s.io/client-go: v0.29.0 → v0.30.0
- k8s.io/component-base: v0.29.0 → v0.30.0
- k8s.io/klog/v2: v2.110.1 → v2.120.1
- k8s.io/kube-openapi: 2dd684a → 70dd376

### Removed
- go.uber.org/atomic: v1.10.0
- k8s.io/gengo: 9cce18d
