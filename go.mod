module github.com/kubernetes-csi/csi-lib-utils

go 1.16

require (
	github.com/container-storage-interface/spec v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023
	google.golang.org/grpc v1.38.0
	k8s.io/api v0.22.0
	k8s.io/client-go v0.22.0
	k8s.io/component-base v0.22.0
	k8s.io/klog/v2 v2.9.0
)

replace (
	k8s.io/api => k8s.io/api v0.22.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.0
	k8s.io/client-go => k8s.io/client-go v0.22.0
	k8s.io/component-base => k8s.io/component-base v0.22.0
	k8s.io/node-api => k8s.io/node-api v0.22.0
)
