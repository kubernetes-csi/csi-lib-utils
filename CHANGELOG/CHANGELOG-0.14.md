# Release notes for v0.14.0

# Changelog since v0.13.0

## Changes by Kind

### Feature

- Add the `GetGroupControllerCapabilities` func under the `rpc` directory. ([#133](https://github.com/kubernetes-csi/csi-lib-utils/pull/133), [@carlory](https://github.com/carlory))

### Bug

- Add a timeout logic to the existing func `connect.Connect()` to aviod block infinitely for a caller. When it tries to connect for 30 seconds, this func will return an error if no connection has been established at that point. ([#131](https://github.com/kubernetes-csi/csi-lib-utils/pull/131), [@ConnorJC3](https://github.com/ConnorJC3))

### Other (Cleanup or Flake)

- Update dependency go modules for k8s v1.27.0. ([#132](https://github.com/kubernetes-csi/csi-lib-utils/pull/132), [@carlory](https://github.com/carlory))

## Dependencies

### Added

- github.com/google/uuid: [v1.3.0](https://github.com/google/uuid/tree/v1.3.0)

### Changed

- github.com/container-storage-interface/spec: [v1.7.0 → v1.8.0](https://github.com/container-storage-interface/spec/compare/v1.8.0...v1.7.0)
- github.com/go-openapi/jsonpointer: [v0.19.5 → v0.19.6](https://github.com/go-openapi/jsonpointer/compare/v0.19.6...v0.19.5)
- github.com/go-openapi/jsonreference: [v0.20.0 → v0.20.1](https://github.com/go-openapi/jsonreference/compare/v0.20.1...v0.20.0)
- github.com/go-openapi/swag: [v0.19.14 → v0.22.3](https://github.com/go-openapi/swag/compare/v0.22.3...v0.19.14)
- github.com/golang/protobuf: [v1.5.2 → v1.5.3](https://github.com/golang/protobuf/compare/v1.5.3...v1.5.2)
- github.com/mailru/easyjson: [v0.7.6 → v0.7.7](https://github.com/mailru/easyjson/compare/v0.7.7...v0.7.7)
- github.com/stretchr/testify: [v1.8.1 → v1.8.2](https://github.com/stretchr/testify/compare/v1.8.2...v1.8.1)
- golang.org/x/net: [v0.4.0 → v0.8.0](https://golang.org/x/net/compare/v0.8.0...v0.4.0)
- golang.org/x/sys: [v0.3.0 → v0.6.0](https://golang.org/x/sys/compare/v0.6.0...v0.3.0)
- golang.org/x/term: [v0.3.0 → v0.6.0](https://golang.org/x/term/compare/v0.6.0...v0.3.0)
- golang.org/x/text: [v0.5.0 → v0.8.0](https://golang.org/x/text/compare/v0.8.0...v0.5.0)
- google.golang.org/grpc: [v1.49.0 → v1.51.0](https://google.golang.org/grpc/compare/v1.51.0...v1.49.0)
- k8s.io/api: [v0.26.0 → v0.27.0](https://github.com/kubernetes/api/compare/v0.27.0...v0.26.0)
- k8s.io/apimachinery: [v0.26.0 → v0.27.0](https://github.com/kubernetes/apimachinery/compare/v0.27.0...v0.26.0)
- k8s.io/client-go: [v0.26.0 → v0.27.0](https://github.com/kubernetes/client-go/compare/v0.27.0...v0.26.0)
- k8s.io/klog/v2: [v2.80.1 → v2.90.1](https://github.com/kubernetes/klog/compare/v2.90.1...v2.80.1)
- k8s.io/kube-openapi: v0.0.0-20221012153701-172d655c2280 → v0.0.0-20230308215209-15aac26d736a
- k8s.io/utils: v0.0.0-20221107191617-1a15be271d1d → v0.0.0-20230209194617-a36077c30491
- sigs.k8s.io/json: v0.0.0-20220713155537-f223a00ba0e2 → v0.0.0-20221116044647-bc3834ca7abd
### Removed

_Nothing has changed._
