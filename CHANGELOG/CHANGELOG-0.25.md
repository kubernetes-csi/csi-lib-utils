# Release notes for v0.25.0

# Changelog since v0.24.0

## Changes by Kind

### Feature

- Updated dependencies to Kubernetes 1.37. ([#217](https://github.com/kubernetes-csi/csi-lib-utils/pull/217), [@jsafrane](https://github.com/jsafrane))

## Dependencies

### Added
- cloud.google.com/go/auth: v0.18.2
- github.com/go-openapi/swag/cmdutils: [v0.29.1](https://github.com/go-openapi/swag/tree/cmdutils/v0.29.1)
- github.com/go-openapi/swag/conv: [v0.29.1](https://github.com/go-openapi/swag/tree/conv/v0.29.1)
- github.com/go-openapi/swag/fileutils: [v0.29.1](https://github.com/go-openapi/swag/tree/fileutils/v0.29.1)
- github.com/go-openapi/swag/jsonutils/fixtures_test: [v0.29.1](https://github.com/go-openapi/swag/tree/jsonutils/fixtures_test/v0.29.1)
- github.com/go-openapi/swag/jsonutils: [v0.29.1](https://github.com/go-openapi/swag/tree/jsonutils/v0.29.1)
- github.com/go-openapi/swag/loading: [v0.29.1](https://github.com/go-openapi/swag/tree/loading/v0.29.1)
- github.com/go-openapi/swag/mangling: [v0.29.1](https://github.com/go-openapi/swag/tree/mangling/v0.29.1)
- github.com/go-openapi/swag/netutils: [v0.29.1](https://github.com/go-openapi/swag/tree/netutils/v0.29.1)
- github.com/go-openapi/swag/pools: [v0.29.1](https://github.com/go-openapi/swag/tree/pools/v0.29.1)
- github.com/go-openapi/swag/stringutils: [v0.29.1](https://github.com/go-openapi/swag/tree/stringutils/v0.29.1)
- github.com/go-openapi/swag/typeutils: [v0.29.1](https://github.com/go-openapi/swag/tree/typeutils/v0.29.1)
- github.com/go-openapi/swag/yamlutils: [v0.29.1](https://github.com/go-openapi/swag/tree/yamlutils/v0.29.1)
- github.com/go-openapi/testify/enable/yaml/v2: [v2.6.1](https://github.com/go-openapi/testify/tree/enable/yaml/v2/v2.6.1)
- github.com/go-openapi/testify/v2: [v2.6.1](https://github.com/go-openapi/testify/tree/v2.6.1)
- github.com/google/s2a-go: [v0.1.9](https://github.com/google/s2a-go/tree/v0.1.9)
- github.com/googleapis/enterprise-certificate-proxy: [v0.3.11](https://github.com/googleapis/enterprise-certificate-proxy/tree/v0.3.11)
- github.com/googleapis/gax-go/v2: [v2.17.0](https://github.com/googleapis/gax-go/tree/v2.17.0)

### Changed
- cel.dev/expr: v0.25.1 → v0.25.2
- github.com/Azure/go-ansiterm: [306776e → faa5f7b](https://github.com/Azure/go-ansiterm/compare/306776e...faa5f7b)
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.30.0 → v1.33.0](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/compare/detectors/gcp/v1.30.0...detectors/gcp/v1.33.0)
- github.com/cncf/xds/go: [ee656c7 → dba9d58](https://github.com/cncf/xds/compare/ee656c7...dba9d58)
- github.com/container-storage-interface/spec: [v1.11.0 → v1.13.0](https://github.com/container-storage-interface/spec/compare/v1.11.0...v1.13.0)
- github.com/envoyproxy/go-control-plane/envoy: [v1.36.0 → v1.37.0](https://github.com/envoyproxy/go-control-plane/compare/envoy/v1.36.0...envoy/v1.37.0)
- github.com/envoyproxy/protoc-gen-validate: [v1.3.0 → v1.3.3](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.3.0...v1.3.3)
- github.com/fxamacker/cbor/v2: [v2.9.0 → v2.9.3](https://github.com/fxamacker/cbor/compare/v2.9.0...v2.9.3)
- github.com/go-jose/go-jose/v4: [v4.1.3 → v4.1.4](https://github.com/go-jose/go-jose/compare/v4.1.3...v4.1.4)
- github.com/go-logr/logr: [v1.4.3 → v1.4.4](https://github.com/go-logr/logr/compare/v1.4.3...v1.4.4)
- github.com/go-openapi/jsonpointer: [v0.21.0 → v1.0.0](https://github.com/go-openapi/jsonpointer/compare/v0.21.0...v1.0.0)
- github.com/go-openapi/jsonreference: [v0.21.0 → v1.0.1](https://github.com/go-openapi/jsonreference/compare/v0.21.0...v1.0.1)
- github.com/go-openapi/swag: [v0.23.0 → v0.29.1](https://github.com/go-openapi/swag/compare/v0.23.0...v0.29.1)
- github.com/golang-jwt/jwt/v5: [v5.3.0 → v5.3.1](https://github.com/golang-jwt/jwt/compare/v5.3.0...v5.3.1)
- github.com/google/gnostic-models: [v0.7.0 → v0.7.1](https://github.com/google/gnostic-models/compare/v0.7.0...v0.7.1)
- github.com/grpc-ecosystem/grpc-gateway/v2: [v2.27.7 → v2.29.0](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.27.7...v2.29.0)
- github.com/klauspost/compress: [v1.18.0 → v1.19.1](https://github.com/klauspost/compress/compare/v1.18.0...v1.19.1)
- github.com/moby/term: [v0.5.0 → v0.5.2](https://github.com/moby/term/compare/v0.5.0...v0.5.2)
- github.com/prometheus/client_golang: [v1.23.2 → v1.24.1](https://github.com/prometheus/client_golang/compare/v1.23.2...v1.24.1)
- github.com/prometheus/common: [v0.67.5 → v0.70.1](https://github.com/prometheus/common/compare/v0.67.5...v0.70.1)
- github.com/prometheus/procfs: [v0.19.2 → v0.21.1](https://github.com/prometheus/procfs/compare/v0.19.2...v0.21.1)
- github.com/spf13/pflag: [v1.0.9 → v1.0.10](https://github.com/spf13/pflag/compare/v1.0.9...v1.0.10)
- github.com/spiffe/go-spiffe/v2: [v2.6.0 → v2.7.0](https://github.com/spiffe/go-spiffe/compare/v2.6.0...v2.7.0)
- github.com/stretchr/objx: [v0.5.2 → v0.5.3](https://github.com/stretchr/objx/compare/v0.5.2...v0.5.3)
- github.com/stretchr/testify: [v1.11.1 → v1.12.1](https://github.com/stretchr/testify/compare/v1.11.1...v1.12.1)
- go.opentelemetry.io/contrib/detectors/gcp: v1.39.0 → v1.44.0
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.58.0 → v0.71.0
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.65.0 → v0.69.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.40.0 → v1.44.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.40.0 → v1.44.0
- go.opentelemetry.io/otel/metric: v1.41.0 → v1.46.0
- go.opentelemetry.io/otel/sdk/metric: v1.39.0 → v1.46.0
- go.opentelemetry.io/otel/sdk: v1.40.0 → v1.46.0
- go.opentelemetry.io/otel/trace: v1.41.0 → v1.46.0
- go.opentelemetry.io/otel: v1.41.0 → v1.46.0
- go.opentelemetry.io/proto/otlp: v1.9.0 → v1.10.0
- go.yaml.in/yaml/v2: v2.4.3 → v2.4.4
- go.yaml.in/yaml/v3: v3.0.4 → v3.0.5
- golang.org/x/crypto: v0.47.0 → v0.55.0
- golang.org/x/mod: v0.31.0 → v0.38.0
- golang.org/x/net: v0.49.0 → v0.58.0
- golang.org/x/oauth2: v0.34.0 → v0.36.0
- golang.org/x/sync: v0.19.0 → v0.22.0
- golang.org/x/sys: v0.40.0 → v0.47.0
- golang.org/x/term: v0.39.0 → v0.45.0
- golang.org/x/text: v0.33.0 → v0.41.0
- golang.org/x/time: v0.14.0 → v0.15.0
- golang.org/x/tools: v0.40.0 → v0.48.0
- gonum.org/v1/gonum: v0.16.0 → v0.17.0
- google.golang.org/genproto/googleapis/api: 8636f87 → 3dc84a4
- google.golang.org/genproto/googleapis/rpc: 8636f87 → da73d73
- google.golang.org/grpc: v1.79.3 → v1.83.2
- google.golang.org/protobuf: f2248ac → v1.36.12
- k8s.io/api: v0.36.0 → v0.37.0
- k8s.io/apimachinery: v0.36.0 → v0.37.0
- k8s.io/client-go: v0.36.0 → v0.37.0
- k8s.io/component-base: v0.36.0 → v0.37.0
- k8s.io/gengo/v2: 85fd79d → ec3ebc5
- k8s.io/kube-openapi: 43fb72c → be32def
- k8s.io/streaming: v0.36.0 → v0.37.0
- k8s.io/utils: b8788ab → cf1189d
- sigs.k8s.io/structured-merge-diff/v6: v6.3.2 → v6.4.2

### Removed
- github.com/josharian/intern: [v1.0.0](https://github.com/josharian/intern/tree/v1.0.0)
- github.com/mailru/easyjson: [v0.9.0](https://github.com/mailru/easyjson/tree/v0.9.0)
