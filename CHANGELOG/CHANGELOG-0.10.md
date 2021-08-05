# Release notes for v0.10.0

# Changelog since v0.9.1

## Changes by Kind

### Feature

- Adds mappings for PV access modes to new CSI access modes: `SINGLE_NODE_SINGLE_WRITER` and `SINGLE_NODE_MULTI_WRITER`. ([#89](https://github.com/kubernetes-csi/csi-lib-utils/pull/89), [@chrishenzie](https://github.com/chrishenzie))
- Updates Kubernetes dependencies to v1.22.0 ([#98](https://github.com/kubernetes-csi/csi-lib-utils/pull/98), [@chrishenzie](https://github.com/chrishenzie))
- Updates container-storage-interface dependency to v1.5.0 ([#90](https://github.com/kubernetes-csi/csi-lib-utils/pull/90), [@chrishenzie](https://github.com/chrishenzie))

### Other (Cleanup or Flake)

- Fixes issue where sidecars exit without output about goroutines when losing the connection to the driver ([#81](https://github.com/kubernetes-csi/csi-lib-utils/pull/81), [@pohly](https://github.com/pohly))

### Uncategorized

- Updated runtime (Go 1.16) and dependencies ([#82](https://github.com/kubernetes-csi/csi-lib-utils/pull/82), [@pohly](https://github.com/pohly))

## Dependencies

### Added
- cloud.google.com/go/firestore: v1.1.0
- github.com/Azure/go-autorest: [v14.2.0+incompatible](https://github.com/Azure/go-autorest/tree/v14.2.0)
- github.com/OneOfOne/xxhash: [v1.2.2](https://github.com/OneOfOne/xxhash/tree/v1.2.2)
- github.com/antihax/optional: [v1.0.0](https://github.com/antihax/optional/tree/v1.0.0)
- github.com/armon/circbuf: [bbbad09](https://github.com/armon/circbuf/tree/bbbad09)
- github.com/armon/go-metrics: [f0300d1](https://github.com/armon/go-metrics/tree/f0300d1)
- github.com/armon/go-radix: [7fddfc3](https://github.com/armon/go-radix/tree/7fddfc3)
- github.com/asaskevich/govalidator: [f61b66f](https://github.com/asaskevich/govalidator/tree/f61b66f)
- github.com/benbjohnson/clock: [v1.0.3](https://github.com/benbjohnson/clock/tree/v1.0.3)
- github.com/bgentry/speakeasy: [v0.1.0](https://github.com/bgentry/speakeasy/tree/v0.1.0)
- github.com/bketelsen/crypt: [5cbc8cc](https://github.com/bketelsen/crypt/tree/5cbc8cc)
- github.com/cespare/xxhash: [v1.1.0](https://github.com/cespare/xxhash/tree/v1.1.0)
- github.com/coreos/bbolt: [v1.3.2](https://github.com/coreos/bbolt/tree/v1.3.2)
- github.com/coreos/etcd: [v3.3.13+incompatible](https://github.com/coreos/etcd/tree/v3.3.13)
- github.com/coreos/go-semver: [v0.3.0](https://github.com/coreos/go-semver/tree/v0.3.0)
- github.com/coreos/go-systemd: [95778df](https://github.com/coreos/go-systemd/tree/95778df)
- github.com/coreos/pkg: [399ea9e](https://github.com/coreos/pkg/tree/399ea9e)
- github.com/cpuguy83/go-md2man/v2: [v2.0.0](https://github.com/cpuguy83/go-md2man/v2/tree/v2.0.0)
- github.com/creack/pty: [v1.1.11](https://github.com/creack/pty/tree/v1.1.11)
- github.com/dgryski/go-sip13: [e10d5fe](https://github.com/dgryski/go-sip13/tree/e10d5fe)
- github.com/fatih/color: [v1.7.0](https://github.com/fatih/color/tree/v1.7.0)
- github.com/felixge/httpsnoop: [v1.0.1](https://github.com/felixge/httpsnoop/tree/v1.0.1)
- github.com/form3tech-oss/jwt-go: [v3.2.3+incompatible](https://github.com/form3tech-oss/jwt-go/tree/v3.2.3)
- github.com/go-gl/glfw: [e6da0ac](https://github.com/go-gl/glfw/tree/e6da0ac)
- github.com/go-kit/log: [v0.1.0](https://github.com/go-kit/log/tree/v0.1.0)
- github.com/gopherjs/gopherjs: [0766667](https://github.com/gopherjs/gopherjs/tree/0766667)
- github.com/gorilla/websocket: [v1.4.2](https://github.com/gorilla/websocket/tree/v1.4.2)
- github.com/grpc-ecosystem/go-grpc-middleware: [v1.0.0](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/v1.0.0)
- github.com/grpc-ecosystem/go-grpc-prometheus: [v1.2.0](https://github.com/grpc-ecosystem/go-grpc-prometheus/tree/v1.2.0)
- github.com/grpc-ecosystem/grpc-gateway: [v1.16.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.16.0)
- github.com/hashicorp/consul/api: [v1.1.0](https://github.com/hashicorp/consul/api/tree/v1.1.0)
- github.com/hashicorp/consul/sdk: [v0.1.1](https://github.com/hashicorp/consul/sdk/tree/v0.1.1)
- github.com/hashicorp/errwrap: [v1.0.0](https://github.com/hashicorp/errwrap/tree/v1.0.0)
- github.com/hashicorp/go-cleanhttp: [v0.5.1](https://github.com/hashicorp/go-cleanhttp/tree/v0.5.1)
- github.com/hashicorp/go-immutable-radix: [v1.0.0](https://github.com/hashicorp/go-immutable-radix/tree/v1.0.0)
- github.com/hashicorp/go-msgpack: [v0.5.3](https://github.com/hashicorp/go-msgpack/tree/v0.5.3)
- github.com/hashicorp/go-multierror: [v1.0.0](https://github.com/hashicorp/go-multierror/tree/v1.0.0)
- github.com/hashicorp/go-rootcerts: [v1.0.0](https://github.com/hashicorp/go-rootcerts/tree/v1.0.0)
- github.com/hashicorp/go-sockaddr: [v1.0.0](https://github.com/hashicorp/go-sockaddr/tree/v1.0.0)
- github.com/hashicorp/go-syslog: [v1.0.0](https://github.com/hashicorp/go-syslog/tree/v1.0.0)
- github.com/hashicorp/go-uuid: [v1.0.1](https://github.com/hashicorp/go-uuid/tree/v1.0.1)
- github.com/hashicorp/go.net: [v0.0.1](https://github.com/hashicorp/go.net/tree/v0.0.1)
- github.com/hashicorp/hcl: [v1.0.0](https://github.com/hashicorp/hcl/tree/v1.0.0)
- github.com/hashicorp/logutils: [v1.0.0](https://github.com/hashicorp/logutils/tree/v1.0.0)
- github.com/hashicorp/mdns: [v1.0.0](https://github.com/hashicorp/mdns/tree/v1.0.0)
- github.com/hashicorp/memberlist: [v0.1.3](https://github.com/hashicorp/memberlist/tree/v0.1.3)
- github.com/hashicorp/serf: [v0.8.2](https://github.com/hashicorp/serf/tree/v0.8.2)
- github.com/inconshreveable/mousetrap: [v1.0.0](https://github.com/inconshreveable/mousetrap/tree/v1.0.0)
- github.com/jonboulle/clockwork: [v0.1.0](https://github.com/jonboulle/clockwork/tree/v0.1.0)
- github.com/jpillora/backoff: [v1.0.0](https://github.com/jpillora/backoff/tree/v1.0.0)
- github.com/jtolds/gls: [v4.20.0+incompatible](https://github.com/jtolds/gls/tree/v4.20.0)
- github.com/magiconair/properties: [v1.8.1](https://github.com/magiconair/properties/tree/v1.8.1)
- github.com/mattn/go-colorable: [v0.0.9](https://github.com/mattn/go-colorable/tree/v0.0.9)
- github.com/mattn/go-isatty: [v0.0.3](https://github.com/mattn/go-isatty/tree/v0.0.3)
- github.com/miekg/dns: [v1.0.14](https://github.com/miekg/dns/tree/v1.0.14)
- github.com/mitchellh/cli: [v1.0.0](https://github.com/mitchellh/cli/tree/v1.0.0)
- github.com/mitchellh/go-homedir: [v1.1.0](https://github.com/mitchellh/go-homedir/tree/v1.1.0)
- github.com/mitchellh/go-testing-interface: [v1.0.0](https://github.com/mitchellh/go-testing-interface/tree/v1.0.0)
- github.com/mitchellh/gox: [v0.4.0](https://github.com/mitchellh/gox/tree/v0.4.0)
- github.com/mitchellh/iochan: [v1.0.0](https://github.com/mitchellh/iochan/tree/v1.0.0)
- github.com/mitchellh/mapstructure: [v1.1.2](https://github.com/mitchellh/mapstructure/tree/v1.1.2)
- github.com/moby/spdystream: [v0.2.0](https://github.com/moby/spdystream/tree/v0.2.0)
- github.com/niemeyer/pretty: [a10e7ca](https://github.com/niemeyer/pretty/tree/a10e7ca)
- github.com/nxadm/tail: [v1.4.4](https://github.com/nxadm/tail/tree/v1.4.4)
- github.com/oklog/ulid: [v1.3.1](https://github.com/oklog/ulid/tree/v1.3.1)
- github.com/pascaldekloe/goe: [57f6aae](https://github.com/pascaldekloe/goe/tree/57f6aae)
- github.com/pelletier/go-toml: [v1.2.0](https://github.com/pelletier/go-toml/tree/v1.2.0)
- github.com/posener/complete: [v1.1.1](https://github.com/posener/complete/tree/v1.1.1)
- github.com/prometheus/tsdb: [v0.7.1](https://github.com/prometheus/tsdb/tree/v0.7.1)
- github.com/rogpeppe/fastuuid: [v1.2.0](https://github.com/rogpeppe/fastuuid/tree/v1.2.0)
- github.com/russross/blackfriday/v2: [v2.0.1](https://github.com/russross/blackfriday/v2/tree/v2.0.1)
- github.com/ryanuber/columnize: [9b3edd6](https://github.com/ryanuber/columnize/tree/9b3edd6)
- github.com/sean-/seed: [e2103e2](https://github.com/sean-/seed/tree/e2103e2)
- github.com/shurcooL/sanitized_anchor_name: [v1.0.0](https://github.com/shurcooL/sanitized_anchor_name/tree/v1.0.0)
- github.com/smartystreets/assertions: [b2de0cb](https://github.com/smartystreets/assertions/tree/b2de0cb)
- github.com/smartystreets/goconvey: [v1.6.4](https://github.com/smartystreets/goconvey/tree/v1.6.4)
- github.com/soheilhy/cmux: [v0.1.4](https://github.com/soheilhy/cmux/tree/v0.1.4)
- github.com/spaolacci/murmur3: [f09979e](https://github.com/spaolacci/murmur3/tree/f09979e)
- github.com/spf13/cast: [v1.3.0](https://github.com/spf13/cast/tree/v1.3.0)
- github.com/spf13/cobra: [v1.1.3](https://github.com/spf13/cobra/tree/v1.1.3)
- github.com/spf13/jwalterweatherman: [v1.0.0](https://github.com/spf13/jwalterweatherman/tree/v1.0.0)
- github.com/spf13/viper: [v1.7.0](https://github.com/spf13/viper/tree/v1.7.0)
- github.com/stoewer/go-strcase: [v1.2.0](https://github.com/stoewer/go-strcase/tree/v1.2.0)
- github.com/subosito/gotenv: [v1.2.0](https://github.com/subosito/gotenv/tree/v1.2.0)
- github.com/tmc/grpc-websocket-proxy: [0ad062e](https://github.com/tmc/grpc-websocket-proxy/tree/0ad062e)
- github.com/xiang90/probing: [43a291a](https://github.com/xiang90/probing/tree/43a291a)
- github.com/yuin/goldmark: [v1.3.5](https://github.com/yuin/goldmark/tree/v1.3.5)
- go.etcd.io/bbolt: v1.3.2
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.20.0
- go.opentelemetry.io/contrib: v0.20.0
- go.opentelemetry.io/otel/exporters/otlp: v0.20.0
- go.opentelemetry.io/otel/metric: v0.20.0
- go.opentelemetry.io/otel/oteltest: v0.20.0
- go.opentelemetry.io/otel/sdk/export/metric: v0.20.0
- go.opentelemetry.io/otel/sdk/metric: v0.20.0
- go.opentelemetry.io/otel/sdk: v0.20.0
- go.opentelemetry.io/otel/trace: v0.20.0
- go.opentelemetry.io/otel: v0.20.0
- go.opentelemetry.io/proto/otlp: v0.7.0
- golang.org/x/term: 6a3ed07
- gopkg.in/ini.v1: v1.51.0
- gopkg.in/resty.v1: v1.12.0
- gopkg.in/yaml.v3: 496545a
- rsc.io/quote/v3: v3.1.0
- rsc.io/sampler: v1.3.0

### Changed
- cloud.google.com/go/bigquery: v1.0.1 → v1.4.0
- cloud.google.com/go/datastore: v1.0.0 → v1.1.0
- cloud.google.com/go/pubsub: v1.0.1 → v1.2.0
- cloud.google.com/go/storage: v1.0.0 → v1.6.0
- cloud.google.com/go: v0.51.0 → v0.54.0
- github.com/Azure/go-ansiterm: [d6e3b33 → d185dfc](https://github.com/Azure/go-ansiterm/compare/d6e3b33...d185dfc)
- github.com/Azure/go-autorest/autorest/adal: [v0.8.2 → v0.9.13](https://github.com/Azure/go-autorest/autorest/adal/compare/v0.8.2...v0.9.13)
- github.com/Azure/go-autorest/autorest/date: [v0.2.0 → v0.3.0](https://github.com/Azure/go-autorest/autorest/date/compare/v0.2.0...v0.3.0)
- github.com/Azure/go-autorest/autorest/mocks: [v0.3.0 → v0.4.1](https://github.com/Azure/go-autorest/autorest/mocks/compare/v0.3.0...v0.4.1)
- github.com/Azure/go-autorest/autorest: [v0.9.6 → v0.11.18](https://github.com/Azure/go-autorest/autorest/compare/v0.9.6...v0.11.18)
- github.com/Azure/go-autorest/logger: [v0.1.0 → v0.2.1](https://github.com/Azure/go-autorest/logger/compare/v0.1.0...v0.2.1)
- github.com/Azure/go-autorest/tracing: [v0.5.0 → v0.6.0](https://github.com/Azure/go-autorest/tracing/compare/v0.5.0...v0.6.0)
- github.com/PuerkitoBio/purell: [v1.0.0 → v1.1.1](https://github.com/PuerkitoBio/purell/compare/v1.0.0...v1.1.1)
- github.com/PuerkitoBio/urlesc: [5bd2802 → de5bf2a](https://github.com/PuerkitoBio/urlesc/compare/5bd2802...de5bf2a)
- github.com/alecthomas/units: [c3de453 → f65c72e](https://github.com/alecthomas/units/compare/c3de453...f65c72e)
- github.com/blang/semver: [v3.5.0+incompatible → v3.5.1+incompatible](https://github.com/blang/semver/compare/v3.5.0...v3.5.1)
- github.com/cncf/udpa/go: [269d4d4 → 5459f2c](https://github.com/cncf/udpa/go/compare/269d4d4...5459f2c)
- github.com/container-storage-interface/spec: [v1.2.0 → v1.5.0](https://github.com/container-storage-interface/spec/compare/v1.2.0...v1.5.0)
- github.com/envoyproxy/go-control-plane: [v0.9.4 → 668b12f](https://github.com/envoyproxy/go-control-plane/compare/v0.9.4...668b12f)
- github.com/evanphx/json-patch: [v4.9.0+incompatible → v4.11.0+incompatible](https://github.com/evanphx/json-patch/compare/v4.9.0...v4.11.0)
- github.com/ghodss/yaml: [73d445a → v1.0.0](https://github.com/ghodss/yaml/compare/73d445a...v1.0.0)
- github.com/go-gl/glfw/v3.3/glfw: [12ad95a → 6f7a984](https://github.com/go-gl/glfw/v3.3/glfw/compare/12ad95a...6f7a984)
- github.com/go-logfmt/logfmt: [v0.4.0 → v0.5.0](https://github.com/go-logfmt/logfmt/compare/v0.4.0...v0.5.0)
- github.com/go-logr/logr: [v0.2.0 → v0.4.0](https://github.com/go-logr/logr/compare/v0.2.0...v0.4.0)
- github.com/go-openapi/jsonpointer: [46af16f → v0.19.3](https://github.com/go-openapi/jsonpointer/compare/46af16f...v0.19.3)
- github.com/go-openapi/jsonreference: [13c6e35 → v0.19.3](https://github.com/go-openapi/jsonreference/compare/13c6e35...v0.19.3)
- github.com/go-openapi/swag: [1d0bd11 → v0.19.5](https://github.com/go-openapi/swag/compare/1d0bd11...v0.19.5)
- github.com/gogo/protobuf: [v1.3.1 → v1.3.2](https://github.com/gogo/protobuf/compare/v1.3.1...v1.3.2)
- github.com/golang/groupcache: [215e871 → 41bb18b](https://github.com/golang/groupcache/compare/215e871...41bb18b)
- github.com/golang/mock: [v1.3.1 → v1.4.1](https://github.com/golang/mock/compare/v1.3.1...v1.4.1)
- github.com/golang/protobuf: [v1.4.2 → v1.5.2](https://github.com/golang/protobuf/compare/v1.4.2...v1.5.2)
- github.com/google/btree: [v1.0.0 → v1.0.1](https://github.com/google/btree/compare/v1.0.0...v1.0.1)
- github.com/google/go-cmp: [v0.4.0 → v0.5.5](https://github.com/google/go-cmp/compare/v0.4.0...v0.5.5)
- github.com/google/pprof: [d4f498a → 1ebb73c](https://github.com/google/pprof/compare/d4f498a...1ebb73c)
- github.com/google/uuid: [v1.1.1 → v1.1.2](https://github.com/google/uuid/compare/v1.1.1...v1.1.2)
- github.com/googleapis/gnostic: [v0.4.1 → v0.5.5](https://github.com/googleapis/gnostic/compare/v0.4.1...v0.5.5)
- github.com/json-iterator/go: [v1.1.10 → v1.1.11](https://github.com/json-iterator/go/compare/v1.1.10...v1.1.11)
- github.com/julienschmidt/httprouter: [v1.2.0 → v1.3.0](https://github.com/julienschmidt/httprouter/compare/v1.2.0...v1.3.0)
- github.com/kisielk/errcheck: [v1.2.0 → v1.5.0](https://github.com/kisielk/errcheck/compare/v1.2.0...v1.5.0)
- github.com/kr/text: [v0.1.0 → v0.2.0](https://github.com/kr/text/compare/v0.1.0...v0.2.0)
- github.com/mailru/easyjson: [d5b7844 → b2ccc51](https://github.com/mailru/easyjson/compare/d5b7844...b2ccc51)
- github.com/moby/term: [672ec06 → 9d4ed18](https://github.com/moby/term/compare/672ec06...9d4ed18)
- github.com/mwitkow/go-conntrack: [cc309e4 → 2f06839](https://github.com/mwitkow/go-conntrack/compare/cc309e4...2f06839)
- github.com/onsi/ginkgo: [v1.11.0 → v1.14.0](https://github.com/onsi/ginkgo/compare/v1.11.0...v1.14.0)
- github.com/onsi/gomega: [v1.7.0 → v1.10.1](https://github.com/onsi/gomega/compare/v1.7.0...v1.10.1)
- github.com/prometheus/client_golang: [v1.7.1 → v1.11.0](https://github.com/prometheus/client_golang/compare/v1.7.1...v1.11.0)
- github.com/prometheus/common: [v0.10.0 → v0.26.0](https://github.com/prometheus/common/compare/v0.10.0...v0.26.0)
- github.com/prometheus/procfs: [v0.1.3 → v0.6.0](https://github.com/prometheus/procfs/compare/v0.1.3...v0.6.0)
- github.com/stretchr/testify: [v1.5.1 → v1.7.0](https://github.com/stretchr/testify/compare/v1.5.1...v1.7.0)
- go.opencensus.io: v0.22.2 → v0.22.3
- go.uber.org/atomic: v1.4.0 → v1.7.0
- go.uber.org/multierr: v1.1.0 → v1.6.0
- go.uber.org/zap: v1.10.0 → v1.17.0
- golang.org/x/crypto: 75b2880 → 5ea612d
- golang.org/x/exp: da58074 → 6cc2880
- golang.org/x/lint: fdd1cda → 6edffad
- golang.org/x/mod: c90efee → v0.4.2
- golang.org/x/net: ab34263 → 37e1c6a
- golang.org/x/oauth2: 858c2ad → bf48bf1
- golang.org/x/sync: cd5d95a → 036812b
- golang.org/x/sys: ed371f2 → 59db8d7
- golang.org/x/text: v0.3.3 → v0.3.6
- golang.org/x/time: 555d28b → 1f47c86
- golang.org/x/tools: 7b8e75d → v0.1.2
- golang.org/x/xerrors: 9bdfabe → 5ec99f8
- google.golang.org/api: v0.15.0 → v0.20.0
- google.golang.org/genproto: cb27e3a → f16073e
- google.golang.org/grpc: v1.29.0 → v1.38.0
- google.golang.org/protobuf: v1.24.0 → v1.26.0
- gopkg.in/check.v1: 41f04d3 → 8fa4692
- gopkg.in/yaml.v2: v2.2.8 → v2.4.0
- gotest.tools/v3: v3.0.2 → v3.0.3
- honnef.co/go/tools: v0.0.1-2019.2.3 → v0.0.1-2020.1.3
- k8s.io/api: v0.19.0 → v0.22.0
- k8s.io/apimachinery: v0.19.0 → v0.22.0
- k8s.io/client-go: v0.19.0 → v0.22.0
- k8s.io/component-base: v0.19.0 → v0.22.0
- k8s.io/klog/v2: v2.2.0 → v2.9.0
- k8s.io/kube-openapi: 6aeccd4 → 9528897
- k8s.io/utils: d5654de → 4b05e18
- sigs.k8s.io/structured-merge-diff/v4: v4.0.1 → v4.1.2

### Removed
- github.com/docker/spdystream: [449fdfc](https://github.com/docker/spdystream/tree/449fdfc)
- github.com/go-openapi/spec: [6aced65](https://github.com/go-openapi/spec/tree/6aced65)
- gotest.tools: v2.2.0+incompatible
