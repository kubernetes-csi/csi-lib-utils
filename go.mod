module github.com/kubernetes-csi/csi-lib-utils

go 1.12

require (
	github.com/container-storage-interface/spec v1.2.0
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.3.5
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/sys v0.0.0-20200406155108-e3b113bbe6a4 // indirect
	google.golang.org/genproto v0.0.0-20200406120821-33397c535dc2 // indirect
	google.golang.org/grpc v1.28.0
	k8s.io/api v0.18.0
	k8s.io/client-go v1.5.1
	k8s.io/component-base v0.18.0
	k8s.io/klog v1.0.0
)

replace k8s.io/api => k8s.io/api v0.18.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.0

replace k8s.io/apimachinery => k8s.io/apimachinery v0.18.0

replace k8s.io/apiserver => k8s.io/apiserver v0.18.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.0

replace k8s.io/client-go => k8s.io/client-go v0.18.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.0

replace k8s.io/code-generator => k8s.io/code-generator v0.18.0

replace k8s.io/component-base => k8s.io/component-base v0.18.0

replace k8s.io/cri-api => k8s.io/cri-api v0.18.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.0

replace k8s.io/kubectl => k8s.io/kubectl v0.18.0

replace k8s.io/kubelet => k8s.io/kubelet v0.18.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.0

replace k8s.io/metrics => k8s.io/metrics v0.18.0

replace k8s.io/node-api => k8s.io/node-api v0.18.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.0

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.18.0

replace k8s.io/sample-controller => k8s.io/sample-controller v0.18.0
