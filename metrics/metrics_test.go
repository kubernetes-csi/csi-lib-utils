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
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/component-base/metrics/testutil"
)

const (
	SidecarOperationMetric = "csi_sidecar_operations_seconds"
	ProcessStartTimeMetric = "process_start_time_seconds"
)

func TestRecordMetrics(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManager(
		"fake.csi.driver.io" /* driverName */)
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

	if err := testutil.GatherAndCompare(
		cmm.GetRegistry(), strings.NewReader(expectedMetrics), SidecarOperationMetric); err != nil {
		t.Fatal(err)
	}
}

func TestRecordMetrics_NoDriverName(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManager(
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
	cmm := NewCSIMetricsManager(
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

func TestStartMetricsEndPoint_Noop(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManager(
		"fake.csi.driver.io" /* driverName */)
	operationDuration, _ := time.ParseDuration("20s")

	// Act
	cmm.StartMetricsEndpoint(":8080", "/metrics")
	cmm.RecordMetrics(
		"/csi.v1.Controller/ControllerGetCapabilities", /* operationName */
		nil, /* operationErr */
		operationDuration /* operationDuration */)

	// Assert
	request, err := http.NewRequest("GET", "http://localhost:8080/metrics", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Creating request for metrics endpoint failed: %v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Fatalf("Failed to GET metrics. Error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("/metrics response status not 200. Response was: %+v", resp)
	}

	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to parse metrics response.  Response was: %+v Error: %v", resp, err)
	}

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

	actualMetrics := string(contentBytes)
	if err := VerifyMetricsMatch(expectedMetrics, actualMetrics, ProcessStartTimeMetric); err != nil {
		t.Fatalf("Metrics returned by end point do not match expectation: %v", err)
	}
}

func TestProcessStartTimeMetricExist(t *testing.T) {
	// Arrange
	cmm := NewCSIMetricsManager(
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
