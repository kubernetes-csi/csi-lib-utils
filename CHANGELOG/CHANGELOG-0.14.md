# Release notes for v0.14.1

[Documentation](https://kubernetes-csi.github.io)

## Changes by Kind

### Bug or Regression

- Updated google.golang.org/grpc to v1.59.0. To keep existing behavior, gRPC connection idling is force disabled, i.e. a connection is *not* closed after 30 minutes of inactivity. ([#159](https://github.com/kubernetes-csi/csi-lib-utils/pull/159), [@msau42](https://github.com/msau42))

## Dependencies

### Added
- cloud.google.com/go/compute/metadata: v0.2.3
- cloud.google.com/go/compute: v1.23.0
- google.golang.org/genproto/googleapis/api: b8732ec
- google.golang.org/genproto/googleapis/rpc: b8732ec

### Changed
- github.com/census-instrumentation/opencensus-proto: [v0.2.1 → v0.4.1](https://github.com/census-instrumentation/opencensus-proto/compare/v0.2.1...v0.4.1)
- github.com/cespare/xxhash/v2: [v2.1.2 → v2.2.0](https://github.com/cespare/xxhash/v2/compare/v2.1.2...v2.2.0)
- github.com/cncf/udpa/go: [04548b0 → c52dc94](https://github.com/cncf/udpa/go/compare/04548b0...c52dc94)
- github.com/cncf/xds/go: [cb28da3 → e9ce688](https://github.com/cncf/xds/go/compare/cb28da3...e9ce688)
- github.com/envoyproxy/go-control-plane: [49ff273 → v0.11.1](https://github.com/envoyproxy/go-control-plane/compare/49ff273...v0.11.1)
- github.com/envoyproxy/protoc-gen-validate: [v0.1.0 → v1.0.2](https://github.com/envoyproxy/protoc-gen-validate/compare/v0.1.0...v1.0.2)
- github.com/golang/glog: [23def4e → v1.1.2](https://github.com/golang/glog/compare/23def4e...v1.1.2)
- github.com/google/uuid: [v1.3.0 → v1.3.1](https://github.com/google/uuid/compare/v1.3.0...v1.3.1)
- golang.org/x/crypto: 75b2880 → v0.12.0
- golang.org/x/net: v0.8.0 → v0.14.0
- golang.org/x/oauth2: ee48083 → v0.11.0
- golang.org/x/sync: 0de741c → v0.3.0
- golang.org/x/sys: v0.6.0 → v0.11.0
- golang.org/x/term: v0.6.0 → v0.11.0
- golang.org/x/text: v0.8.0 → v0.12.0
- google.golang.org/genproto: c8bf987 → b8732ec
- google.golang.org/grpc: v1.51.0 → v1.59.0
- google.golang.org/protobuf: v1.28.1 → v1.31.0

### Removed
- github.com/antihax/optional: [v1.0.0](https://github.com/antihax/optional/tree/v1.0.0)
- github.com/ghodss/yaml: [v1.0.0](https://github.com/ghodss/yaml/tree/v1.0.0)
- github.com/grpc-ecosystem/grpc-gateway: [v1.16.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.16.0)
- github.com/rogpeppe/fastuuid: [v1.2.0](https://github.com/rogpeppe/fastuuid/tree/v1.2.0)

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
