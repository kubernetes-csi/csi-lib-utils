# Release notes for v0.16.0

# Changelog since v0.15.0

## Changes by Kind

### Uncategorized

- Updated google.golang.org/grpc to v1.59.0. To keep existing behavior, gRPC connection idling is force disabled, i.e. a connection is *not* closed after 30 minutes of inactivity. ([#153](https://github.com/kubernetes-csi/csi-lib-utils/pull/153), [@jsafrane](https://github.com/jsafrane))

## Dependencies

### Added
_Nothing has changed._

### Changed
- cloud.google.com/go/compute: v1.15.1 → v1.23.0
- github.com/cncf/xds/go: [06c439d → e9ce688](https://github.com/cncf/xds/go/compare/06c439d...e9ce688)
- github.com/envoyproxy/go-control-plane: [v0.10.3 → v0.11.1](https://github.com/envoyproxy/go-control-plane/compare/v0.10.3...v0.11.1)
- github.com/envoyproxy/protoc-gen-validate: [v0.9.1 → v1.0.2](https://github.com/envoyproxy/protoc-gen-validate/compare/v0.9.1...v1.0.2)
- github.com/golang/glog: [v1.0.0 → v1.1.2](https://github.com/golang/glog/compare/v1.0.0...v1.1.2)
- github.com/google/uuid: [v1.3.0 → v1.3.1](https://github.com/google/uuid/compare/v1.3.0...v1.3.1)
- github.com/stretchr/testify: [v1.8.2 → v1.8.4](https://github.com/stretchr/testify/compare/v1.8.2...v1.8.4)
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.41.0 → v0.44.0
- go.opentelemetry.io/otel/metric: v0.38.0 → v1.19.0
- go.opentelemetry.io/otel/trace: v1.15.0 → v1.19.0
- go.opentelemetry.io/otel: v1.15.0 → v1.19.0
- golang.org/x/crypto: v0.11.0 → v0.15.0
- golang.org/x/net: v0.13.0 → v0.18.0
- golang.org/x/oauth2: v0.8.0 → v0.11.0
- golang.org/x/sync: v0.2.0 → v0.3.0
- golang.org/x/sys: v0.10.0 → v0.14.0
- golang.org/x/term: v0.10.0 → v0.14.0
- golang.org/x/text: v0.11.0 → v0.14.0
- google.golang.org/genproto/googleapis/api: dd9d682 → b8732ec
- google.golang.org/genproto/googleapis/rpc: 28d5490 → bbf56f3
- google.golang.org/genproto: 0005af6 → d783a09
- google.golang.org/grpc: v1.54.0 → v1.59.0
- google.golang.org/protobuf: v1.30.0 → v1.31.0

### Removed
_Nothing has changed._
