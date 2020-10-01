# Release notes for v0.8.1

# Changelog since v0.8.0

## Changes by Kind

### Feature

- Add `process_start_time_seconds` into the csi metric lib. ([#54](https://github.com/kubernetes-csi/csi-lib-utils/pull/54), [@Jiawei0227](https://github.com/Jiawei0227))

### Bug or Regression

- Fix k8s.io/client-go version in go.mod ([#53](https://github.com/kubernetes-csi/csi-lib-utils/pull/53), [@avorima](https://github.com/avorima))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Release notes for v0.8.0

# Changelog since v0.7.0

## Urgent Upgrade Notes 

### (No, really, you MUST read this before you upgrade)

 - Client-go gets updated to Kubernetes 1.18.0 which changes the API of several functions (context added, create/delete/update options). ([#39](https://github.com/kubernetes-csi/csi-lib-utils/pull/39), [@humblec](https://github.com/humblec))
 
## Changes by Kind

### Feature
 - Leaderelection add WithContext func ([#44](https://github.com/kubernetes-csi/csi-lib-utils/pull/44), [@aimuz](https://github.com/aimuz))
 - Metrics support can also be used by CSI drivers ([#49](https://github.com/kubernetes-csi/csi-lib-utils/pull/49), [@pohly](https://github.com/pohly))

### Uncategorized
 - Build with Go 1.15 ([#48](https://github.com/kubernetes-csi/csi-lib-utils/pull/48), [@pohly](https://github.com/pohly))

## Dependencies

### Added
- github.com/cncf/udpa/go: [269d4d4](https://github.com/cncf/udpa/go/tree/269d4d4)
- github.com/docopt/docopt-go: [ee0de3b](https://github.com/docopt/docopt-go/tree/ee0de3b)
- sigs.k8s.io/structured-merge-diff/v3: v3.0.0

### Changed
- github.com/container-storage-interface/spec: [v1.1.0 → v1.2.0](https://github.com/container-storage-interface/spec/compare/v1.1.0...v1.2.0)
- github.com/elazarl/goproxy: [c4fc265 → 947c36d](https://github.com/elazarl/goproxy/compare/c4fc265...947c36d)
- github.com/envoyproxy/go-control-plane: [5f8ba28 → v0.9.4](https://github.com/envoyproxy/go-control-plane/compare/5f8ba28...v0.9.4)
- github.com/gogo/protobuf: [65acae2 → v1.3.1](https://github.com/gogo/protobuf/compare/65acae2...v1.3.1)
- github.com/golang/groupcache: [5b532d6 → 8c9f03a](https://github.com/golang/groupcache/compare/5b532d6...8c9f03a)
- github.com/golang/protobuf: [v1.3.2 → v1.3.5](https://github.com/golang/protobuf/compare/v1.3.2...v1.3.5)
- github.com/google/gofuzz: [v1.0.0 → v1.1.0](https://github.com/google/gofuzz/compare/v1.0.0...v1.1.0)
- github.com/googleapis/gnostic: [v0.2.0 → v0.3.1](https://github.com/googleapis/gnostic/compare/v0.2.0...v0.3.1)
- github.com/onsi/ginkgo: [v1.10.2 → v1.12.0](https://github.com/onsi/ginkgo/compare/v1.10.2...v1.12.0)
- github.com/onsi/gomega: [v1.7.0 → v1.7.1](https://github.com/onsi/gomega/compare/v1.7.0...v1.7.1)
- github.com/pkg/errors: [v0.8.1 → v0.9.1](https://github.com/pkg/errors/compare/v0.8.1...v0.9.1)
- github.com/prometheus/client_model: [14fe0d1 → v0.2.0](https://github.com/prometheus/client_model/compare/14fe0d1...v0.2.0)
- github.com/stretchr/testify: [v1.4.0 → v1.5.1](https://github.com/stretchr/testify/compare/v1.4.0...v1.5.1)
- golang.org/x/crypto: 60c769a → bac4c82
- golang.org/x/net: c0dbc17 → d3edc99
- golang.org/x/sys: 0732a99 → e3b113b
- google.golang.org/genproto: 5c49e3e → 33397c5
- google.golang.org/grpc: v1.26.0 → v1.28.0
- gopkg.in/yaml.v2: v2.2.4 → v2.2.8
- k8s.io/api: v0.17.0 → v0.18.0
- k8s.io/apimachinery: v0.17.1-beta.0 → v0.18.0
- k8s.io/client-go: v0.17.0 → v0.18.0
- k8s.io/component-base: v0.17.0 → v0.18.0
- k8s.io/kube-openapi: 30be4d1 → bf4fb3b
- k8s.io/utils: e782cd3 → a9aa75a
- sigs.k8s.io/yaml: v1.1.0 → v1.2.0

### Removed
- sigs.k8s.io/structured-merge-diff: 15d366b
