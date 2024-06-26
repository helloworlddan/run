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
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestAddOAuth2Header(t *testing.T) {
	resetInstance()
	method := http.MethodGet
	url := "https://example.com"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal("AddOAuth2Header() fails to instantiate request")
	}

	req = AddOAuth2Header(req)

	val, ok := req.Header["Authorization"]
	if !ok {
		t.Fatal("AddOAuth2Header() contains no 'Authorization' header")
	}

	if len(val) != 1 {
		t.Fatal("AddOAuth2Header() contains malformed 'Authorization' header")
	}

	expect := fmt.Sprintf("bearer: %s", ServiceAccountAccessToken())
	if val[0] != expect {
		t.Fatal("AddOAuth2Header() contains invalid 'Authorization' header")
	}
}

func TestProjectID(t *testing.T) {
	resetInstance()
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

func TestServicePort(t *testing.T) {
	resetInstance()
	envVarKey := "PORT"
	envVarVal := "8081"

	err := os.Setenv(envVarKey, "")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := ServicePort()
	if result != "8080" {
		t.Fatalf(`port() = %s, want "8080"`, result)
	}

	err = os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	resetInstance()
	result = ServicePort()
	if result != envVarVal {
		t.Fatalf(`port() = %s, want %s`, result, envVarVal)
	}
}

func TestServiceName(t *testing.T) {
	resetInstance()
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

	resetInstance()
	result = ServiceName()
	if result != envVarVal {
		t.Fatalf(`KNativeService() = %s, want %s`, result, envVarVal)
	}
}

func TestServiceRevision(t *testing.T) {
	resetInstance()
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

	resetInstance()
	result = ServiceRevision()
	if result != envVarVal || err != nil {
		t.Fatalf(`kNativeRevision() = %s, %v, want %s, error`, result, err, envVarVal)
	}
}

func TestJobName(t *testing.T) {
	resetInstance()
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

	resetInstance()
	result = JobName()
	if result != envVarVal {
		t.Fatalf(`jobName() = %s, want %s`, result, envVarVal)
	}
}

func TestJobExecution(t *testing.T) {
	resetInstance()
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

	resetInstance()
	result = JobExecution()
	if result != envVarVal {
		t.Fatalf(`JobExecution() = %s, want %s`, result, envVarVal)
	}
}

func TestJobTaskIndex(t *testing.T) {
	resetInstance()
	envVarKey := "CLOUD_RUN_TASK_INDEX"
	envVarVal := 12

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := JobTaskIndex()
	if result != -1 {
		t.Fatalf(`jobTaskIndex() = %d, want -1`, result)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	resetInstance()
	result = JobTaskIndex()
	if result != envVarVal {
		t.Fatalf(`JobTaskIndex() = %d, want %d`, result, envVarVal)
	}
}

func TestJobTaskAttempt(t *testing.T) {
	resetInstance()
	envVarKey := "CLOUD_RUN_TASK_ATTEMPT"
	envVarVal := 14

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := JobTaskAttempt()
	if result != -1 {
		t.Fatalf(`JobTaskAttempt() = %d, want -1`, result)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	resetInstance()
	result = JobTaskAttempt()
	if result != envVarVal {
		t.Fatalf(`JobTaskAttempt() = %d, want %d`, result, envVarVal)
	}
}

func TestJobTaskCount(t *testing.T) {
	resetInstance()
	envVarKey := "CLOUD_RUN_TASK_COUNT"
	envVarVal := 16

	err := os.Setenv(envVarKey, "wrong value")
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	result := JobTaskCount()
	if result != -1 {
		t.Fatalf(`JobTaskCount() = %d, want -1`, result)
	}

	err = os.Setenv(envVarKey, fmt.Sprintf("%d", envVarVal))
	if err != nil {
		t.Fatalf("unable to test: %v", err)
	}

	resetInstance()
	result = JobTaskCount()
	if result != envVarVal {
		t.Fatalf(`JobTaskCount() = %d, want %d`, result, envVarVal)
	}
}
