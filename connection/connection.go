package connection

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"google.golang.org/grpc"
)

const (
	// Interval of logging connection errors
	connectionLoggingInterval = 10 * time.Second
)

// Connect opens insecure gRPC connection to a CSI driver. Address must be either absolute path to a socket file
// or have format '<protocol>://', following gRPC name resolution mechanism at https://github.com/grpc/grpc/blob/master/doc/naming.md.
// The function tries to connect indefinitely every second until it connects. The function automatically adds
// interceptor for gRPC message logging.
func Connect(address string, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOptions = append(dialOptions,
		grpc.WithInsecure(),                   // Don't use TLS, it's usually local Unix domain socket in a container.
		grpc.WithBackoffMaxDelay(time.Second), // Retry every second after failure.
		grpc.WithBlock(),                      // Block until connection succeeds.
		grpc.WithUnaryInterceptor(LogGRPC),    // Log all messages.
	)
	if strings.HasPrefix(address, "/") {
		// It looks like filesystem path.
		dialOptions = append(dialOptions, grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
	}
	glog.Infof("Connecting to %s", address)

	// Connect in background.
	var conn *grpc.ClientConn
	var err error
	ready := make(chan bool)
	go func() {
		conn, err = grpc.Dial(address, dialOptions...)
		close(ready)
	}()

	// Log error every connectionLoggingInterval
	ticker := time.NewTicker(connectionLoggingInterval)
	defer ticker.Stop()

	// Wait until Dial() succeeds.
	for {
		select {
		case <-ticker.C:
			glog.Warningf("Still connecting to %s", address)

		case <-ready:
			return conn, err
		}
	}
}

// LogGRPC is gPRC unary interceptor for logging of CSI messages at level 5. It removes any secrets from the message.
func LogGRPC(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	glog.V(5).Infof("GRPC call: %s", method)
	glog.V(5).Infof("GRPC request: %s", protosanitizer.StripSecrets(req))
	err := invoker(ctx, method, req, reply, cc, opts...)
	glog.V(5).Infof("GRPC response: %s", protosanitizer.StripSecrets(reply))
	glog.V(5).Infof("GRPC error: %v", err)
	return err
}
