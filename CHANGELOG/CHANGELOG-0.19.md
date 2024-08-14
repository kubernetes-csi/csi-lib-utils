# Release notes for v0.19.0

# Changelog since v0.18.1

## Changes by Kind

### Other (Cleanup or Flake)

- Update Kubernetes dependencies to 1.31.0 ([#179](https://github.com/kubernetes-csi/csi-lib-utils/pull/179), [@dfajmon](https://github.com/dfajmon))
- Updates Kubernetes dependencies to 1.31.0-beta.0 ([#178](https://github.com/kubernetes-csi/csi-lib-utils/pull/178), [@dfajmon](https://github.com/dfajmon))

## Dependencies

### Added
- cel.dev/expr: v0.15.0
- github.com/go-task/slim-sprig/v3: [v3.0.0](https://github.com/go-task/slim-sprig/tree/v3.0.0)
- gopkg.in/evanphx/json-patch.v4: v4.12.0

### Changed
- cloud.google.com/go/compute/metadata: v0.2.3 → v0.3.0
- github.com/alecthomas/kingpin/v2: [v2.3.2 → v2.4.0](https://github.com/alecthomas/kingpin/compare/v2.3.2...v2.4.0)
- github.com/cenkalti/backoff/v4: [v4.2.1 → v4.3.0](https://github.com/cenkalti/backoff/compare/v4.2.1...v4.3.0)
- github.com/cespare/xxhash/v2: [v2.2.0 → v2.3.0](https://github.com/cespare/xxhash/compare/v2.2.0...v2.3.0)
- github.com/cncf/xds/go: [e9ce688 → 555b57e](https://github.com/cncf/xds/compare/e9ce688...555b57e)
- github.com/davecgh/go-spew: [v1.1.1 → d8f796a](https://github.com/davecgh/go-spew/compare/v1.1.1...d8f796a)
- github.com/envoyproxy/go-control-plane: [v0.11.1 → v0.12.0](https://github.com/envoyproxy/go-control-plane/compare/v0.11.1...v0.12.0)
- github.com/envoyproxy/protoc-gen-validate: [v1.0.2 → v1.0.4](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.0.2...v1.0.4)
- github.com/felixge/httpsnoop: [v1.0.3 → v1.0.4](https://github.com/felixge/httpsnoop/compare/v1.0.3...v1.0.4)
- github.com/fxamacker/cbor/v2: [v2.6.0 → v2.7.0](https://github.com/fxamacker/cbor/compare/v2.6.0...v2.7.0)
- github.com/go-logr/logr: [v1.4.1 → v1.4.2](https://github.com/go-logr/logr/compare/v1.4.1...v1.4.2)
- github.com/go-openapi/swag: [v0.22.3 → v0.22.4](https://github.com/go-openapi/swag/compare/v0.22.3...v0.22.4)
- github.com/golang/glog: [v1.1.2 → v1.2.1](https://github.com/golang/glog/compare/v1.1.2...v1.2.1)
- github.com/google/pprof: [4bb14d4 → 4bfdf5a](https://github.com/google/pprof/compare/4bb14d4...4bfdf5a)
- github.com/google/uuid: [v1.3.1 → v1.6.0](https://github.com/google/uuid/compare/v1.3.1...v1.6.0)
- github.com/grpc-ecosystem/grpc-gateway/v2: [v2.16.0 → v2.20.0](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.16.0...v2.20.0)
- github.com/moby/spdystream: [v0.2.0 → v0.4.0](https://github.com/moby/spdystream/compare/v0.2.0...v0.4.0)
- github.com/moby/term: [1aeaba8 → v0.5.0](https://github.com/moby/term/compare/1aeaba8...v0.5.0)
- github.com/onsi/ginkgo/v2: [v2.15.0 → v2.19.0](https://github.com/onsi/ginkgo/compare/v2.15.0...v2.19.0)
- github.com/onsi/gomega: [v1.31.0 → v1.33.1](https://github.com/onsi/gomega/compare/v1.31.0...v1.33.1)
- github.com/pmezard/go-difflib: [v1.0.0 → 5d4384e](https://github.com/pmezard/go-difflib/compare/v1.0.0...5d4384e)
- github.com/prometheus/client_golang: [v1.16.0 → v1.19.1](https://github.com/prometheus/client_golang/compare/v1.16.0...v1.19.1)
- github.com/prometheus/client_model: [v0.4.0 → v0.6.1](https://github.com/prometheus/client_model/compare/v0.4.0...v0.6.1)
- github.com/prometheus/common: [v0.44.0 → v0.55.0](https://github.com/prometheus/common/compare/v0.44.0...v0.55.0)
- github.com/prometheus/procfs: [v0.10.1 → v0.15.1](https://github.com/prometheus/procfs/compare/v0.10.1...v0.15.1)
- github.com/rogpeppe/go-internal: [v1.10.0 → v1.12.0](https://github.com/rogpeppe/go-internal/compare/v1.10.0...v1.12.0)
- github.com/spf13/cobra: [v1.7.0 → v1.8.1](https://github.com/spf13/cobra/compare/v1.7.0...v1.8.1)
- github.com/stretchr/objx: [v0.5.0 → v0.5.2](https://github.com/stretchr/objx/compare/v0.5.0...v0.5.2)
- github.com/stretchr/testify: [v1.8.4 → v1.9.0](https://github.com/stretchr/testify/compare/v1.8.4...v1.9.0)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.44.0 → v0.53.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.19.0 → v1.27.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.19.0 → v1.28.0
- go.opentelemetry.io/otel/metric: v1.20.0 → v1.28.0
- go.opentelemetry.io/otel/sdk: v1.19.0 → v1.28.0
- go.opentelemetry.io/otel/trace: v1.20.0 → v1.28.0
- go.opentelemetry.io/otel: v1.20.0 → v1.28.0
- go.opentelemetry.io/proto/otlp: v1.0.0 → v1.3.1
- golang.org/x/crypto: v0.21.0 → v0.24.0
- golang.org/x/mod: v0.15.0 → v0.17.0
- golang.org/x/net: v0.23.0 → v0.26.0
- golang.org/x/oauth2: v0.11.0 → v0.21.0
- golang.org/x/sync: v0.3.0 → v0.7.0
- golang.org/x/sys: v0.18.0 → v0.21.0
- golang.org/x/term: v0.18.0 → v0.21.0
- golang.org/x/text: v0.14.0 → v0.16.0
- golang.org/x/tools: v0.18.0 → e35e4cc
- google.golang.org/genproto/googleapis/api: b8732ec → 5315273
- google.golang.org/genproto/googleapis/rpc: bbf56f3 → f6361c8
- google.golang.org/grpc: v1.59.0 → v1.65.0
- google.golang.org/protobuf: v1.33.0 → v1.34.2
- k8s.io/api: v0.30.0 → v0.31.0
- k8s.io/apimachinery: v0.30.0 → v0.31.0
- k8s.io/client-go: v0.30.0 → v0.31.0
- k8s.io/component-base: v0.30.0 → v0.31.0
- k8s.io/klog/v2: v2.120.1 → v2.130.1
- k8s.io/utils: 3b25d92 → 18e509b
- sigs.k8s.io/yaml: v1.3.0 → v1.4.0

### Removed
- github.com/cncf/udpa/go: [c52dc94](https://github.com/cncf/udpa/tree/c52dc94)
- github.com/evanphx/json-patch: [v4.12.0+incompatible](https://github.com/evanphx/json-patch/tree/v4.12.0)
- github.com/go-task/slim-sprig: [52ccab3](https://github.com/go-task/slim-sprig/tree/52ccab3)
- github.com/matttproud/golang_protobuf_extensions: [v1.0.4](https://github.com/matttproud/golang_protobuf_extensions/tree/v1.0.4)
- google.golang.org/genproto: d783a09
v2.130.1
- k8s.io/utils: 3b25d92 → 18e509b
- sigs.k8s.io/yaml: v1.3.0 → v1.4.0

### Removed
- github.com/cncf/udpa/go: [c52dc94](https://github.com/cncf/udpa/tree/c52dc94)
- github.com/evanphx/json-patch: [v4.12.0+incompatible](https://github.com/evanphx/json-patch/tree/v4.12.0)
- github.com/go-task/slim-sprig: [52ccab3](https://github.com/go-task/slim-sprig/tree/52ccab3)
- github.com/matttproud/golang_protobuf_extensions: [v1.0.4](https://github.com/matttproud/golang_protobuf_extensions/tree/v1.0.4)
- google.golang.org/genproto: d783a09
