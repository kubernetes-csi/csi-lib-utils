module github.com/kubernetes-csi/csi-lib-utils

go 1.16

require (
	github.com/container-storage-interface/spec v1.4.0
	github.com/golang/protobuf v1.4.3
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7
	google.golang.org/grpc v1.37.0
	k8s.io/api v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/component-base v0.21.1
	k8s.io/klog/v2 v2.8.0
)

replace k8s.io/api => k8s.io/api v0.21.1

replace k8s.io/apimachinery => k8s.io/apimachinery v0.21.1

replace k8s.io/client-go => k8s.io/client-go v0.21.1

replace k8s.io/component-base => k8s.io/component-base v0.21.1

replace k8s.io/node-api => k8s.io/node-api v0.21.1
