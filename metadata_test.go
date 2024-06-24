// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package run

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestNewAuthenticatedRequest(t *testing.T) {
	ctx := context.Background()
	method := http.MethodGet
	url := "https://example.com"
	req, err := NewAuthenticatedRequest(ctx, method, url, nil)
	if err != nil {
		t.Fatal("authenticatedRequest() fails to instantiate request")
	}

	val, ok := req.Header["Authorization"]
	if !ok {
		t.Fatal("authenticatedRequest() contains no 'Authorization' header")
	}

	if len(val) != 1 {
		t.Fatal("authenticatedRequest() contains malformed 'Authorization' header")
	}

	expect := fmt.Sprintf("bearer: %s", ServiceAccountToken())
	if val[0] != expect {
		t.Fatal("authenticatedRequest() contains invalid 'Authorization' header")
	}

	if req.URL.String() != url {
		t.Fatal("authenticatedRequest() constructed bad URL")
	}

	if req.Method != method {
		t.Fatal("authenticatedRequest() constructed bad URL")
	}
}

func TestProjectID(t *testing.T) {
	envVarKey := "GOOGLE_CLOUD_PROJECT"
	envVarVal := "some-valid-project"

	err := os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := ProjectID()
	if result != envVarVal {
		t.Fatalf(`projectID() = %s, want %s`, result, envVarVal)
	}
}

func TestPort(t *testing.T) {
	envVarKey := "PORT"
	envVarVal := "8081"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := Port()
	if result != "8080" {
		t.Fatalf(`port() = %s, want "8080"`, result)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result = Port()
	if result != envVarVal {
		t.Fatalf(`port() = %s, want %s`, result, envVarVal)
	}
}

func TestServiceName(t *testing.T) {
	envVarKey := "K_SERVICE"
	envVarVal := "service-001"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := ServiceName()
	if result != "local" {
		t.Fatalf(`KNativeService() = %s, want "local"`, result)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result = ServiceName()
	if result != envVarVal {
		t.Fatalf(`KNativeService() = %s, want %s`, result, envVarVal)
	}
}

func TestServiceRevision(t *testing.T) {
	envVarKey := "K_REVISION"
	envVarVal := "revision-001"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := ServiceRevision()
	expected := fmt.Sprintf("%s-00001-xxx", Name())
	if result != expected {
		t.Fatalf(`KNativeRevision() = %s, want %s`, result, expected)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result = ServiceRevision()
	if result != envVarVal || err != nil {
		t.Fatalf(`kNativeRevision() = %s, %v, want %s, error`, result, err, envVarVal)
	}
}

func TestJobName(t *testing.T) {
	envVarKey := "CLOUD_RUN_JOB"
	envVarVal := "job-001"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := JobName()
	if result != "local" {
		t.Fatalf(`JobName() = %s want "local"`, result)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result = JobName()
	if result != envVarVal {
		t.Fatalf(`jobName() = %s, want %s`, result, envVarVal)
	}
}

func TestJobExecution(t *testing.T) {
	envVarKey := "CLOUD_RUN_EXECUTION"
	envVarVal := "job-execution-001"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := JobExecution()
	if result != "local" {
		t.Fatalf(`JobExecution() = %s, want "local"`, result)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result = JobExecution()
	if result != envVarVal {
		t.Fatalf(`JobExecution() = %s, want %s`, result, envVarVal)
	}
}

func TestJobTaskIndex(t *testing.T) {
	envVarKey := "CLOUD_RUN_TASK_INDEX"
	envVarVal := 12

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := JobTaskIndex()
	// NOTE: maybe i want -1 here to indicate a problem?
	if result != 0 {
		t.Fatalf(`jobTaskIndex() = %d, want 0`, result)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result = JobTaskIndex()
	if result != envVarVal {
		t.Fatalf(`JobTaskIndex() = %d, want %d`, result, envVarVal)
	}
}

func TestJobTaskAttempt(t *testing.T) {
	envVarKey := "CLOUD_RUN_TASK_ATTEMPT"
	envVarVal := 14

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := JobTaskAttempt()
	// NOTE: maybe i want -1 here to indicate a problem?
	if result != 0 {
		t.Fatalf(`JobTaskAttempt() = %d, want 0`, result)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result = JobTaskAttempt()
	if result != envVarVal {
		t.Fatalf(`JobTaskAttempt() = %d, want %d`, result, envVarVal)
	}
}

func TestJobTaskCount(t *testing.T) {
	envVarKey := "CLOUD_RUN_TASK_COUNT"
	envVarVal := 16

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := JobTaskCount()
	// NOTE: maybe i want -1 here to indicate a problem?
	if result != 0 {
		t.Fatalf(`JobTaskCount() = %d, want 0`, result)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result = JobTaskCount()
	if result != envVarVal {
		t.Fatalf(`JobTaskCount() = %d, want %d`, result, envVarVal)
	}
}
