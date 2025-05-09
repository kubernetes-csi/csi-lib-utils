# Release notes for v0.20.0

# Changelog since v0.19.0

## Changes by Kind

### Other (Cleanup or Flake)

- The output format of protosanitizer is slightly refined. ([#184](https://github.com/kubernetes-csi/csi-lib-utils/pull/184), [@huww98](https://github.com/huww98))
- Update Kubernetes dependencies to 1.32.0 ([#185](https://github.com/kubernetes-csi/csi-lib-utils/pull/185), [@dfajmon](https://github.com/dfajmon))

## Dependencies

### Added
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.24.2](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/tree/detectors/gcp/v1.24.2)
- github.com/klauspost/compress: [v1.17.11](https://github.com/klauspost/compress/tree/v1.17.11)
- github.com/kylelemons/godebug: [v1.1.0](https://github.com/kylelemons/godebug/tree/v1.1.0)
- github.com/planetscale/vtprotobuf: [0393e58](https://github.com/planetscale/vtprotobuf/tree/0393e58)
- go.opentelemetry.io/auto/sdk: v1.1.0
- go.opentelemetry.io/contrib/detectors/gcp: v1.31.0
- go.opentelemetry.io/otel/sdk/metric: v1.31.0

### Changed
- cel.dev/expr: v0.15.0 → v0.16.2
- cloud.google.com/go/compute/metadata: v0.3.0 → v0.5.2
- github.com/Azure/go-ansiterm: [d185dfc → 306776e](https://github.com/Azure/go-ansiterm/compare/d185dfc...306776e)
- github.com/NYTimes/gziphandler: [56545f4 → v1.1.1](https://github.com/NYTimes/gziphandler/compare/56545f4...v1.1.1)
- github.com/cncf/xds/go: [555b57e → b4127c9](https://github.com/cncf/xds/compare/555b57e...b4127c9)
- github.com/container-storage-interface/spec: [v1.9.0 → v1.11.0](https://github.com/container-storage-interface/spec/compare/v1.9.0...v1.11.0)
- github.com/emicklei/go-restful/v3: [v3.11.0 → v3.12.1](https://github.com/emicklei/go-restful/compare/v3.11.0...v3.12.1)
- github.com/envoyproxy/go-control-plane: [v0.12.0 → v0.13.1](https://github.com/envoyproxy/go-control-plane/compare/v0.12.0...v0.13.1)
- github.com/envoyproxy/protoc-gen-validate: [v1.0.4 → v1.1.0](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.0.4...v1.1.0)
- github.com/go-openapi/jsonpointer: [v0.19.6 → v0.21.0](https://github.com/go-openapi/jsonpointer/compare/v0.19.6...v0.21.0)
- github.com/go-openapi/jsonreference: [v0.20.2 → v0.21.0](https://github.com/go-openapi/jsonreference/compare/v0.20.2...v0.21.0)
- github.com/go-openapi/swag: [v0.22.4 → v0.23.0](https://github.com/go-openapi/swag/compare/v0.22.4...v0.23.0)
- github.com/golang/glog: [v1.2.1 → v1.2.2](https://github.com/golang/glog/compare/v1.2.1...v1.2.2)
- github.com/google/gnostic-models: [v0.6.8 → v0.6.9](https://github.com/google/gnostic-models/compare/v0.6.8...v0.6.9)
- github.com/google/pprof: [4bfdf5a → d1b30fe](https://github.com/google/pprof/compare/4bfdf5a...d1b30fe)
- github.com/gregjones/httpcache: [9cad4c3 → 901d907](https://github.com/gregjones/httpcache/compare/9cad4c3...901d907)
- github.com/mailru/easyjson: [v0.7.7 → v0.9.0](https://github.com/mailru/easyjson/compare/v0.7.7...v0.9.0)
- github.com/moby/spdystream: [v0.4.0 → v0.5.0](https://github.com/moby/spdystream/compare/v0.4.0...v0.5.0)
- github.com/onsi/ginkgo/v2: [v2.19.0 → v2.21.0](https://github.com/onsi/ginkgo/compare/v2.19.0...v2.21.0)
- github.com/onsi/gomega: [v1.33.1 → v1.35.1](https://github.com/onsi/gomega/compare/v1.33.1...v1.35.1)
- github.com/prometheus/client_golang: [v1.19.1 → v1.20.5](https://github.com/prometheus/client_golang/compare/v1.19.1...v1.20.5)
- github.com/prometheus/common: [v0.55.0 → v0.61.0](https://github.com/prometheus/common/compare/v0.55.0...v0.61.0)
- github.com/rogpeppe/go-internal: [v1.12.0 → v1.13.1](https://github.com/rogpeppe/go-internal/compare/v1.12.0...v1.13.1)
- github.com/stretchr/testify: [v1.9.0 → v1.10.0](https://github.com/stretchr/testify/compare/v1.9.0...v1.10.0)
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.46.0 → v0.58.0
- go.opentelemetry.io/otel/metric: v1.28.0 → v1.33.0
- go.opentelemetry.io/otel/sdk: v1.28.0 → v1.31.0
- go.opentelemetry.io/otel/trace: v1.28.0 → v1.33.0
- go.opentelemetry.io/otel: v1.28.0 → v1.33.0
- go.uber.org/zap: v1.26.0 → v1.27.0
- golang.org/x/crypto: v0.24.0 → v0.30.0
- golang.org/x/mod: v0.17.0 → v0.20.0
- golang.org/x/net: v0.26.0 → v0.32.0
- golang.org/x/oauth2: v0.21.0 → v0.24.0
- golang.org/x/sync: v0.7.0 → v0.10.0
- golang.org/x/sys: v0.21.0 → v0.28.0
- golang.org/x/term: v0.21.0 → v0.27.0
- golang.org/x/text: v0.16.0 → v0.21.0
- golang.org/x/time: v0.3.0 → v0.8.0
- golang.org/x/tools: e35e4cc → v0.26.0
- golang.org/x/xerrors: 04be3eb → 5ec99f8
- google.golang.org/genproto/googleapis/api: 5315273 → 796eee8
- google.golang.org/genproto/googleapis/rpc: f6361c8 → 9240e9c
- google.golang.org/grpc: v1.65.0 → v1.69.0
- google.golang.org/protobuf: v1.34.2 → v1.36.0
- k8s.io/api: v0.31.0 → v0.32.0
- k8s.io/apimachinery: v0.31.0 → v0.32.0
- k8s.io/client-go: v0.31.0 → v0.32.0
- k8s.io/component-base: v0.31.0 → v0.32.0
- k8s.io/gengo/v2: 51d4e06 → a7b603a
- k8s.io/kube-openapi: 70dd376 → 2c72e55
- k8s.io/utils: 18e509b → 24370be
- sigs.k8s.io/json: bc3834c → cfa47c3
- sigs.k8s.io/structured-merge-diff/v4: v4.4.1 → v4.5.0

### Removed
- cloud.google.com/go/compute: v1.23.0
- github.com/asaskevich/govalidator: [f61b66f](https://github.com/asaskevich/govalidator/tree/f61b66f)
- github.com/creack/pty: [v1.1.9](https://github.com/creack/pty/tree/v1.1.9)
- github.com/golang/groupcache: [41bb18b](https://github.com/golang/groupcache/tree/41bb18b)
- github.com/imdario/mergo: [v0.3.6](https://github.com/imdario/mergo/tree/v0.3.6)
- github.com/kr/pty: [v1.1.1](https://github.com/kr/pty/tree/v1.1.1)
- google.golang.org/appengine: v1.6.7
