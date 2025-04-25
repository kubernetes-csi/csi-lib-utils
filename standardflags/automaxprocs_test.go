/*
Copyright 2025 The Kubernetes Authors.

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

package standardflags

import (
	"flag"
	"testing"
)

func TestAutomaxprocsArgument(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		enabled     bool
		expectError bool
	}{
		{
			name:    "with value as true",
			value:   "true",
			enabled: true,
		},
		{
			name:    "with value as false",
			value:   "false",
			enabled: false,
		},
		{
			name:    "without value",
			enabled: true,
		},
		{
			name:        "with invalid value",
			value:       "error",
			expectError: true,
		},
	}

	if flag.Lookup("automaxprocs") == nil {
		AddAutomaxprocs(nil) // pass t.Logf to see the logs
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := flag.Set("automaxprocs", test.value)
			if !test.expectError && err != nil {
				t.Errorf("test %s: failed to set value to %q", test.name, test.value)
			}
		})
	}
}

func TestEnableDisableAutomaxprocs(t *testing.T) {
	// make sure the flags is not enabled yet
	f := flag.Lookup("automaxprocs")
	if f == nil {
		AddAutomaxprocs(nil)
	}
	if automaxprocsIsEnabled() {
		handleAutomaxprocs("false")
	}

	EnableAutomaxprocs()
	if !automaxprocsIsEnabled() {
		t.Errorf("failed to enable automaxprocs")
	}

	// disable again
	handleAutomaxprocs("false")
	if automaxprocsIsEnabled() {
		t.Errorf("failed to disable automaxprocs")
	}
}
