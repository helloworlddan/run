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

func Test_newAuthenticatedRequest(t *testing.T) {
	service := NewService()
	ctx := context.Background()
	method := http.MethodGet
	url := "https://example.com"
	req, err := newAuthenticatedRequest(service, ctx, method, url, nil)
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

	expect := fmt.Sprintf("bearer: %s", service.ServiceAccountToken())
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

func Test_port(t *testing.T) {
	envVarKey := "PORT"
	envVarVal := "8081"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err := port()
	if err == nil {
		t.Fatalf(`port() = %s, %v, want 0, error`, result, err)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err = port()
	if result != envVarVal || err != nil {
		t.Fatalf(`port() = %s, %v, want %s, error`, result, err, envVarVal)
	}
}

func Test_kNativeService(t *testing.T) {
	envVarKey := "K_SERVICE"
	envVarVal := "service-001"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err := kNativeService()
	if err == nil {
		t.Fatalf(`kNativeService() = %s, %v, want 0, error`, result, err)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err = kNativeService()
	if result != envVarVal || err != nil {
		t.Fatalf(`kNativeService() = %s, %v, want %s, error`, result, err, envVarVal)
	}
}

func Test_kNativeRevision(t *testing.T) {
	envVarKey := "K_REVISION"
	envVarVal := "revision-001"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err := kNativeRevision()
	if err == nil {
		t.Fatalf(`kNativeRevision() = %s, %v, want 0, error`, result, err)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err = kNativeRevision()
	if result != envVarVal || err != nil {
		t.Fatalf(`kNativeRevision() = %s, %v, want %s, error`, result, err, envVarVal)
	}
}

func Test_jobName(t *testing.T) {
	envVarKey := "CLOUD_RUN_JOB"
	envVarVal := "job-001"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err := jobName()
	if err == nil {
		t.Fatalf(`jobName() = %s, %v, want 0, error`, result, err)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err = jobName()
	if result != envVarVal || err != nil {
		t.Fatalf(`jobName() = %s, %v, want %s, error`, result, err, envVarVal)
	}
}

func Test_jobExecution(t *testing.T) {
	envVarKey := "CLOUD_RUN_EXECUTION"
	envVarVal := "job-execution-001"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err := jobExecution()
	if err == nil {
		t.Fatalf(`jobExecution() = %s, %v, want 0, error`, result, err)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err = jobExecution()
	if result != envVarVal || err != nil {
		t.Fatalf(`jobExecution() = %s, %v, want %s, error`, result, err, envVarVal)
	}
}

func Test_jobTaskIndex(t *testing.T) {
	envVarKey := "CLOUD_RUN_TASK_INDEX"
	envVarVal := 12

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err := jobTaskIndex()
	if err == nil {
		t.Fatalf(`jobTaskIndex() = %d, %v, want 0, error`, result, err)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err = jobTaskIndex()
	if result != envVarVal || err != nil {
		t.Fatalf(`jobTaskIndex() = %d, %v, want %d, error`, result, err, envVarVal)
	}
}

func Test_jobTaskAttempt(t *testing.T) {
	envVarKey := "CLOUD_RUN_TASK_ATTEMPT"
	envVarVal := 14

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err := jobTaskAttempt()
	if err == nil {
		t.Fatalf(`jobTaskAttempt() = %d, %v, want 0, error`, result, err)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err = jobTaskAttempt()
	if result != envVarVal || err != nil {
		t.Fatalf(`jobTaskAttempt() = %d, %v, want %d, error`, result, err, envVarVal)
	}
}

func Test_jobTaskCount(t *testing.T) {
	envVarKey := "CLOUD_RUN_TASK_COUNT"
	envVarVal := 16

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err := jobTaskCount()
	if err == nil {
		t.Fatalf(`jobTaskCount() = %d, %v, want 0, error`, result, err)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result, err = jobTaskCount()
	if result != envVarVal || err != nil {
		t.Fatalf(`jobTaskCount() = %d, %v, want %d, error`, result, err, envVarVal)
	}
}
