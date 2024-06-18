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
	"os"
	"testing"
)

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
