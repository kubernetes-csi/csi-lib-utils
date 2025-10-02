package config

import (
	"context"
	"net/http"

	"github.com/kubernetes-csi/csi-lib-utils/features"
	"github.com/kubernetes-csi/csi-lib-utils/leaderelection"
	"github.com/kubernetes-csi/csi-lib-utils/metrics"
	"github.com/kubernetes-csi/csi-lib-utils/standardflags"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func BuildConfig(kubeconfig string, opts standardflags.SidecarConfiguration) (*rest.Config, error) {
	config, err := buildConfig(kubeconfig)
	if err != nil {
		return config, err
	}
	config.QPS = float32(opts.KubeAPIQPS)
	config.Burst = opts.KubeAPIBurst
	return config, nil
}


func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func RunWithLeaderElection(ctx context.Context, 
	config *rest.Config, 
	opts standardflags.SidecarConfiguration, 
	run func(context.Context), 
	driverName string, 
	metricsManager metrics.CSIMetricsManager) {
	
	logger := klog.Background()
	
	// Prepare http endpoint for metrics + leader election healthz
	mux := http.NewServeMux()
	addr := opts.MetricsAddress
	if addr == "" {
		addr = opts.HttpEndpoint
	}

	if addr != "" {
		metricsManager.RegisterToServer(mux, opts.MetricsPath)
		metricsManager.SetDriverName(driverName)
		go func() {
			logger.Info("ServeMux listening", "address", addr, "metricsPath", opts.MetricsPath)
			err := http.ListenAndServe(addr, mux)
			if err != nil {
				logger.Error(err, "Failed to start HTTP server at specified address and metrics path", "address", addr, "metricsPath", opts.MetricsPath)
				klog.FlushAndExit(klog.ExitFlushTimeout, 1)
			}
		}()
	}
	
	if !opts.LeaderElection {
		run(klog.NewContext(context.Background(), logger))
	} else {
		// Create a new clientset for leader election. When the attacher
		// gets busy and its client gets throttled, the leader election
		// can proceed without issues.
		leClientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			logger.Error(err, "Failed to create leaderelection client")
			klog.FlushAndExit(klog.ExitFlushTimeout, 1)
		}

		// Name of config map with leader election lock
		le := leaderelection.NewLeaderElection(leClientset, driverName, run)
		if opts.HttpEndpoint != "" {
			le.PrepareHealthCheck(mux, leaderelection.DefaultHealthCheckTimeout)
		}

		if opts.LeaderElectionNamespace != "" {
			le.WithNamespace(opts.LeaderElectionNamespace)
		}

		// TODO: uncomment once https://github.com/kubernetes-csi/csi-lib-utils/pull/200 is merged
		//if opts.LeaderElectionLabels != nil {
			// le.WithLabels(opts.LeaderElectionLabels)
		//}

		le.WithLeaseDuration(opts.LeaderElectionLeaseDuration)
		le.WithRenewDeadline(opts.LeaderElectionRenewDeadline)
		le.WithRetryPeriod(opts.LeaderElectionRetryPeriod)
		if utilfeature.DefaultFeatureGate.Enabled(features.ReleaseLeaderElectionOnExit) {
			le.WithReleaseOnCancel(true)
			le.WithContext(ctx)
		}

		if err := le.Run(); err != nil {
			logger.Error(err, "Failed to initialize leader election")
			klog.FlushAndExit(klog.ExitFlushTimeout, 1)
		}
	}
}