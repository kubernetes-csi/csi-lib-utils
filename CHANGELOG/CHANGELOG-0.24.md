# Release notes for v0.24.0

# Changelog since v0.23.0

## Changes by Kind

### Bug or Regression

- Fixed crash when parsing cmdline option `--leader-election-labels`. ([#209](https://github.com/kubernetes-csi/csi-lib-utils/pull/209), [@DerekFrank](https://github.com/DerekFrank))

### Other (Cleanup or Flake)

- Removed deprecated OpenTelemetry call and bumped OpenTelemetry to v1.38 to be compatible with Kubernetes 1.35 deps. ([#205](https://github.com/kubernetes-csi/csi-lib-utils/pull/205), [@jsafrane](https://github.com/jsafrane))
- Updated Kubernetes libraries to 1.36.0 ([#213](https://github.com/kubernetes-csi/csi-lib-utils/pull/213), [@jsafrane](https://github.com/jsafrane))
- Updated help text of `--csi-address` command line argument. ([#206](https://github.com/kubernetes-csi/csi-lib-utils/pull/206), [@jsafrane](https://github.com/jsafrane))

### Uncategorized

- Bump google.golang.org/grpc to 1.79.3 a for CVE fix ([#211](https://github.com/kubernetes-csi/csi-lib-utils/pull/211), [@chimanjain](https://github.com/chimanjain))

## Dependencies

### Added
- github.com/cenkalti/backoff/v5: [v5.0.3](https://github.com/cenkalti/backoff/tree/v5.0.3)
- github.com/golang-jwt/jwt/v5: [v5.3.0](https://github.com/golang-jwt/jwt/tree/v5.3.0)
- golang.org/x/tools/go/expect: v0.1.0-deprecated
- golang.org/x/tools/go/packages/packagestest: v0.1.1-deprecated
- gonum.org/v1/gonum: v0.16.0
- k8s.io/streaming: v0.36.0

### Changed
- cel.dev/expr: v0.20.0 → v0.25.1
- cloud.google.com/go/compute/metadata: v0.6.0 → v0.9.0
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.26.0 → v1.30.0](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/compare/detectors/gcp/v1.26.0...detectors/gcp/v1.30.0)
- github.com/alecthomas/units: [b94a6e3 → 0f3dac3](https://github.com/alecthomas/units/compare/b94a6e3...0f3dac3)
- github.com/cncf/xds/go: [2f00578 → ee656c7](https://github.com/cncf/xds/compare/2f00578...ee656c7)
- github.com/emicklei/go-restful/v3: [v3.12.2 → v3.13.0](https://github.com/emicklei/go-restful/compare/v3.12.2...v3.13.0)
- github.com/envoyproxy/go-control-plane/envoy: [v1.32.4 → v1.36.0](https://github.com/envoyproxy/go-control-plane/compare/envoy/v1.32.4...envoy/v1.36.0)
- github.com/envoyproxy/go-control-plane: [v0.13.4 → v0.14.0](https://github.com/envoyproxy/go-control-plane/compare/v0.13.4...v0.14.0)
- github.com/envoyproxy/protoc-gen-validate: [v1.2.1 → v1.3.0](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.2.1...v1.3.0)
- github.com/go-jose/go-jose/v4: [v4.0.4 → v4.1.3](https://github.com/go-jose/go-jose/compare/v4.0.4...v4.1.3)
- github.com/go-logr/logr: [v1.4.2 → v1.4.3](https://github.com/go-logr/logr/compare/v1.4.2...v1.4.3)
- github.com/golang/glog: [v1.2.4 → v1.2.5](https://github.com/golang/glog/compare/v1.2.4...v1.2.5)
- github.com/grpc-ecosystem/grpc-gateway/v2: [v2.26.3 → v2.27.7](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.26.3...v2.27.7)
- github.com/moby/spdystream: [v0.5.0 → v0.5.1](https://github.com/moby/spdystream/compare/v0.5.0...v0.5.1)
- github.com/prometheus/client_golang: [v1.22.0 → v1.23.2](https://github.com/prometheus/client_golang/compare/v1.22.0...v1.23.2)
- github.com/prometheus/client_model: [v0.6.1 → v0.6.2](https://github.com/prometheus/client_model/compare/v0.6.1...v0.6.2)
- github.com/prometheus/common: [v0.62.0 → v0.67.5](https://github.com/prometheus/common/compare/v0.62.0...v0.67.5)
- github.com/prometheus/procfs: [v0.15.1 → v0.19.2](https://github.com/prometheus/procfs/compare/v0.15.1...v0.19.2)
- github.com/rogpeppe/go-internal: [v1.13.1 → v1.14.1](https://github.com/rogpeppe/go-internal/compare/v1.13.1...v1.14.1)
- github.com/spf13/cobra: [v1.9.1 → v1.10.2](https://github.com/spf13/cobra/compare/v1.9.1...v1.10.2)
- github.com/spf13/pflag: [v1.0.6 → v1.0.9](https://github.com/spf13/pflag/compare/v1.0.6...v1.0.9)
- github.com/spiffe/go-spiffe/v2: [v2.5.0 → v2.6.0](https://github.com/spiffe/go-spiffe/compare/v2.5.0...v2.6.0)
- github.com/stretchr/testify: [v1.10.0 → v1.11.1](https://github.com/stretchr/testify/compare/v1.10.0...v1.11.1)
- go.opentelemetry.io/auto/sdk: v1.1.0 → v1.2.1
- go.opentelemetry.io/contrib/detectors/gcp: v1.34.0 → v1.39.0
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.58.0 → v0.65.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.34.0 → v1.40.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.34.0 → v1.40.0
- go.opentelemetry.io/otel/metric: v1.35.0 → v1.41.0
- go.opentelemetry.io/otel/sdk/metric: v1.34.0 → v1.39.0
- go.opentelemetry.io/otel/sdk: v1.34.0 → v1.40.0
- go.opentelemetry.io/otel/trace: v1.35.0 → v1.41.0
- go.opentelemetry.io/otel: v1.35.0 → v1.41.0
- go.opentelemetry.io/proto/otlp: v1.5.0 → v1.9.0
- go.uber.org/zap: v1.27.0 → v1.27.1
- go.yaml.in/yaml/v2: v2.4.2 → v2.4.3
- golang.org/x/crypto: v0.36.0 → v0.47.0
- golang.org/x/mod: v0.20.0 → v0.31.0
- golang.org/x/net: v0.38.0 → v0.49.0
- golang.org/x/oauth2: v0.27.0 → v0.34.0
- golang.org/x/sync: v0.12.0 → v0.19.0
- golang.org/x/sys: v0.31.0 → v0.40.0
- golang.org/x/term: v0.30.0 → v0.39.0
- golang.org/x/text: v0.23.0 → v0.33.0
- golang.org/x/time: v0.9.0 → v0.14.0
- golang.org/x/tools: v0.26.0 → v0.40.0
- google.golang.org/genproto/googleapis/api: a0af3ef → 8636f87
- google.golang.org/genproto/googleapis/rpc: a0af3ef → 8636f87
- google.golang.org/grpc: v1.72.1 → v1.79.3
- google.golang.org/protobuf: v1.36.5 → f2248ac
- gopkg.in/evanphx/json-patch.v4: v4.12.0 → v4.13.0
- k8s.io/api: v0.34.1 → v0.36.0
- k8s.io/apimachinery: v0.34.1 → v0.36.0
- k8s.io/client-go: v0.34.1 → v0.36.0
- k8s.io/component-base: v0.34.1 → v0.36.0
- k8s.io/klog/v2: v2.130.1 → v2.140.0
- k8s.io/kube-openapi: f3f2b99 → 43fb72c
- k8s.io/utils: 4c0f3b2 → b8788ab
- sigs.k8s.io/json: cfa47c3 → 2d32026
- sigs.k8s.io/structured-merge-diff/v6: v6.3.0 → v6.3.2

### Removed
- github.com/armon/go-socks5: [e753329](https://github.com/armon/go-socks5/tree/e753329)
- github.com/cenkalti/backoff/v4: [v4.3.0](https://github.com/cenkalti/backoff/tree/v4.3.0)
- github.com/go-task/slim-sprig/v3: [v3.0.0](https://github.com/go-task/slim-sprig/tree/v3.0.0)
- github.com/gogo/protobuf: [v1.3.2](https://github.com/gogo/protobuf/tree/v1.3.2)
- github.com/google/pprof: [d1b30fe](https://github.com/google/pprof/tree/d1b30fe)
- github.com/gregjones/httpcache: [901d907](https://github.com/gregjones/httpcache/tree/901d907)
- github.com/kisielk/errcheck: [v1.5.0](https://github.com/kisielk/errcheck/tree/v1.5.0)
- github.com/kisielk/gotool: [v1.0.0](https://github.com/kisielk/gotool/tree/v1.0.0)
- github.com/onsi/ginkgo/v2: [v2.21.0](https://github.com/onsi/ginkgo/tree/v2.21.0)
- github.com/onsi/gomega: [v1.35.1](https://github.com/onsi/gomega/tree/v1.35.1)
- github.com/pkg/errors: [v0.9.1](https://github.com/pkg/errors/tree/v0.9.1)
- github.com/yuin/goldmark: [v1.2.1](https://github.com/yuin/goldmark/tree/v1.2.1)
- github.com/zeebo/errs: [v1.4.0](https://github.com/zeebo/errs/tree/v1.4.0)
- golang.org/x/xerrors: 5ec99f8
- gopkg.in/yaml.v2: v2.4.0
