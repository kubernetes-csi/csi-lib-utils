/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/testutil"
)

const (
	SidecarOperationMetric = "csi_sidecar_operations_seconds"
	ProcessStartTimeMetric = "process_start_time_seconds"
)

func TestRecordMetrics(t *testing.T) {
	testcases := map[string]struct {
		subsystem      string
		stabilityLevel metrics.StabilityLevel
	}{
		"default": {},
		"sidecar": {subsystem: SubsystemSidecar},
		"driver":  {subsystem: SubsystemPlugin},
		"other":   {subsystem: "other"},
		"stable":  {stabilityLevel: metrics.STABLE},
		"alpha":   {stabilityLevel: metrics.ALPHA},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			testRecordMetrics(t, tc.subsystem, tc.stabilityLevel)
		})
	}
}

func testRecordMetrics(t *testing.T, subsystem string, stabilityLevel metrics.StabilityLevel) {
	// Arrange
	var cmm CSIMetricsManager
	driverName := "fake.csi.driver.io"
	if stabilityLevel == "" {
		// Cover the two dedicated calls.
		switch subsystem {
		case SubsystemSidecar:
			cmm = NewCSIMetricsManagerForSidecar(driverName)
		case SubsystemPlugin:
			cmm = NewCSIMetricsManagerForPlugin(driverName)
		}
	}
	if cmm == nil {
		// The flexible construction is the fallback.
		var options []MetricsManagerOption
		if subsystem != "" {
			options = append(options, WithSubsystem(subsystem))
		}
		if stabilityLevel != "" {
			options = append(options, WithStabilityLevel(stabilityLevel))
		}
		cmm = NewCSIMetricsManagerWithOptions(driverName, options...)
	}

	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmm.RecordMetrics(
		"/csi.v1.Controller/ControllerGetCapabilities", /* operationName */
		nil, /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="1"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="10"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="15"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="25"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="50"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="120"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="300"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="600"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities"} 20
		csi_sidecar_operations_seconds_count{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities"} 1
		`
	metricName := SidecarOperationMetric
	if subsystem != "" {
		expectedMetrics = strings.Replace(expectedMetrics, "csi_sidecar", subsystem, -1)
		metricName = strings.Replace(metricName, "csi_sidecar", subsystem, -1)
	}
	if stabilityLevel != "" {
		expectedMetrics = strings.Replace(expectedMetrics, "ALPHA", string(stabilityLevel), -1)
	}

	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), metricName); err != nil {
		t.Fatal(err)
	}
}

func TestFixedLabels(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerWithOptions(
		"", /* driverName */
		WithLabels(map[string]string{"a": "111", "b": "222"}),
	)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmm.RecordMetrics(
		"myOperation", /* operationName */
		nil,           /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="10"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="15"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="25"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="50"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="120"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="300"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="600"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 20
		csi_sidecar_operations_seconds_count{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 1
	`

	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), SidecarOperationMetric); err != nil {
		t.Fatal(err)
	}
}

func TestVaryingLabels(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerWithOptions(
		"", /* driverName */
		WithLabelNames("a", "b"),
	)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmmv, err := cmm.WithLabelValues(map[string]string{"a": "111"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmmv, err = cmmv.WithLabelValues(map[string]string{"b": "222"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmmv.RecordMetrics(
		"myOperation", /* operationName */
		nil,           /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="10"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="15"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="25"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="50"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="120"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="300"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="600"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 20
		csi_sidecar_operations_seconds_count{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 1
	`

	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), SidecarOperationMetric); err != nil {
		t.Fatal(err)
	}
}

func TestTwoVaryingLabels(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerWithOptions(
		"", /* driverName */
		WithLabelNames("a", "b"),
	)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmmv, err := cmm.WithLabelValues(map[string]string{"a": "111"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmmv2, err := cmmv.WithLabelValues(map[string]string{"b": "222"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmmv2.RecordMetrics(
		"myOperation", /* operationName */
		nil,           /* operationErr */
		operationDuration /* operationDuration */)
	cmmv3, err := cmmv.WithLabelValues(map[string]string{"b": "xxx"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmmv3.RecordMetrics(
		"myOtherOperation", /* operationName */
		nil,                /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="10"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="15"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="25"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="50"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="120"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="300"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="600"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 20
		csi_sidecar_operations_seconds_count{a="111",b="222",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="10"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="15"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="25"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="50"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="120"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="300"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="600"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation"} 20
		csi_sidecar_operations_seconds_count{a="111",b="xxx",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOtherOperation"} 1
	`

	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), SidecarOperationMetric); err != nil {
		t.Fatal(err)
	}
}

func TestVaryingLabelsBackfill(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerWithOptions(
		"", /* driverName */
		WithLabelNames("a", "b"),
	)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmmv, err := cmm.WithLabelValues(map[string]string{"a": "111"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmmv.RecordMetrics(
		"myOperation", /* operationName */
		nil,           /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="10"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="15"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="25"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="50"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="120"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="300"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="600"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 20
		csi_sidecar_operations_seconds_count{a="111",b="",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 1
	`

	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), SidecarOperationMetric); err != nil {
		t.Fatal(err)
	}
}

func TestVaryingLabels_NameError(t *testing.T) {
	cmm := NewCSIMetricsManagerWithOptions(
		"", /* driverName */
		WithLabelNames("a", "b"),
	)

	_, err := cmm.WithLabelValues(map[string]string{"c": "111"})
	if err == nil {
		t.Fatal("unexpected success")
	}
}

func TestVaryingLabels_OverwriteError(t *testing.T) {
	cmm := NewCSIMetricsManagerWithOptions(
		"", /* driverName */
		WithLabelNames("a", "b"),
	)

	cmmv, err := cmm.WithLabelValues(map[string]string{"a": "111", "b": "222"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = cmmv.WithLabelValues(map[string]string{"a": "111"})
	if err == nil {
		t.Fatal("unexpected success")
	}
}

func TestCombinedLabels(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerWithOptions(
		"", /* driverName */
		WithLabelNames("a", "b"),
		WithLabels(map[string]string{"c": "333", "d": "444"}),
	)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmmv, err := cmm.WithLabelValues(map[string]string{"a": "111", "b": "222"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmmv.RecordMetrics(
		"myOperation", /* operationName */
		nil,           /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="1"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="5"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="10"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="15"} 0
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="25"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="50"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="120"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="300"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="600"} 1
		csi_sidecar_operations_seconds_bucket{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 20
		csi_sidecar_operations_seconds_count{a="111",b="222",c="333",d="444",driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 1
	`

	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), SidecarOperationMetric); err != nil {
		t.Fatal(err)
	}
}

func TestRecordMetrics_NoDriverName(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerForSidecar(
		"" /* driverName */)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmm.RecordMetrics(
		"myOperation", /* operationName */
		nil,           /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="1"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="10"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="15"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="25"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="50"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="120"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="300"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="600"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 20
		csi_sidecar_operations_seconds_count{driver_name="unknown-driver",grpc_status_code="OK",method_name="myOperation"} 1
	`

	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), SidecarOperationMetric); err != nil {
		t.Fatal(err)
	}
}

func TestRecordMetrics_Negative(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerForSidecar(
		"fake.csi.driver.io" /* driverName */)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmm.RecordMetrics(
		"myOperation", /* operationName */
		status.Error(codes.InvalidArgument, "invalid input"), /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="1"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="10"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="15"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="25"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="50"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="120"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="300"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="600"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation"} 20
		csi_sidecar_operations_seconds_count{driver_name="fake.csi.driver.io",grpc_status_code="InvalidArgument",method_name="myOperation"} 1
		`
	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), SidecarOperationMetric); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterToServer_Noop(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerForSidecar(
		"fake.csi.driver.io" /* driverName */)
	operationDuration, _ := time.ParseDuration("20s")
	mux := http.NewServeMux()

	// Act
	cmm.RegisterToServer(mux, "/metrics")
	cmm.RecordMetrics(
		"/csi.v1.Controller/ControllerGetCapabilities", /* operationName */
		nil, /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	request := httptest.NewRequest("GET", "/metrics", strings.NewReader(""))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, request)
	resp := rec.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("/metrics response status not 200. Response was: %+v", resp)
	}

	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to parse metrics response.  Response was: %+v Error: %v", resp, err)
	}
	actualMetrics := string(contentBytes)

	expectedMetrics := `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
		# TYPE csi_sidecar_operations_seconds histogram
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="0.1"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="0.25"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="0.5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="1"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="2.5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="5"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="10"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="15"} 0
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="25"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="50"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="120"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="300"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="600"} 1
		csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities",le="+Inf"} 1
		csi_sidecar_operations_seconds_sum{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities"} 20
		csi_sidecar_operations_seconds_count{driver_name="fake.csi.driver.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerGetCapabilities"} 1
	`

	if err := VerifyMetricsMatch(expectedMetrics, actualMetrics, ProcessStartTimeMetric); err != nil {
		t.Fatalf("Metrics returned by end point do not match expectation: %v", err)
	}
}

func TestRegisterPprofToServer_AllEndpointsAvailable(t *testing.T) {
	endpoints := []string{
		"/debug/pprof/",
		"/debug/pprof/cmdline",
		"/debug/pprof/profile",
		"/debug/pprof/symbol",
		"/debug/pprof/trace",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			testRegisterPprofToServer_AllEndpointsAvailable(t, endpoint)
		})
	}
}

func testRegisterPprofToServer_AllEndpointsAvailable(t *testing.T, endpoint string) {
	// Arrange
	cmm := NewCSIMetricsManagerForSidecar(
		"fake.csi.driver.io" /* driverName */)
	mux := http.NewServeMux()

	// Act
	cmm.RegisterPprofToServer(mux)

	// Assert
	request := httptest.NewRequest("GET", endpoint, strings.NewReader(""))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, request)
	resp := rec.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("%s response status not 200. Response was: %+v", endpoint, resp)
	}

	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to parse pprof index response.  Response was: %+v Error: %v", resp, err)
	}

	// Other endpoints return binary data
	if endpoint == "/debug/pprof/" {
		actualPprofIndex := string(contentBytes)

		// This is the exepcted index html page if pprof is running
		expectedPprofIndexSubstr := `<body>
/debug/pprof/
<br>
<p>Set debug=1 as a query parameter to export in legacy text format</p>
<br>
Types of profiles available:
<table>`

		if ok := strings.Contains(actualPprofIndex, expectedPprofIndexSubstr); !ok {
			t.Fatalf("Pprof index returned by end point do not match expectation. Expected: %s \nGot: %s", expectedPprofIndexSubstr, actualPprofIndex)
		}
	}
}

func TestProcessStartTimeMetricExist(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManagerForSidecar(
		"fake.csi.driver.io" /* driverName */)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmm.RecordMetrics(
		"/csi.v1.Controller/ControllerGetCapabilities", /* operationName */
		nil, /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	metricsFamilies, err := cmm.GetRegistry().Gather()
	if err != nil {
		t.Fatalf("Error fetching metrics: %v", err)
	}

	// check process_start_time_seconds exist
	for _, metricsFamily := range metricsFamilies {
		if metricsFamily.GetName() == ProcessStartTimeMetric {
			return
		}
	}

	t.Fatalf("Metrics does not contain %v. Scraped content: %v", ProcessStartTimeMetric, metricsFamilies)
}
