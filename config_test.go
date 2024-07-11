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

package run_test

import (
	"os"
	"slices"
	"testing"

	"github.com/helloworlddan/run"
)

func TestPutConfig(t *testing.T) {
	run.ResetConfig()

	run.PutConfig("some key", "some val")

	if run.CountConfig() != 1 {
		t.Fatalf("PutConfig() failed to add config correctly")
	}
}

func TestGetConfig(t *testing.T) {
	run.ResetConfig()

	configKey := "test.config"
	configVal := "test.config.val"
	run.PutConfig(configKey, configVal)

	_, err := run.GetConfig("non-existent")
	if err == nil {
		t.Fatalf("GetConfig() failed to err on non-existent config")
	}

	result, err := run.GetConfig(configKey)
	if err != nil {
		t.Fatalf("GetConfig() failed to retrieve existing config")
	}

	if result != configVal {
		t.Fatalf("GetConfig() failed to store config correctly")
	}
}

func TestListConfigKeys(t *testing.T) {
	run.ResetConfig()

	keys := run.ListConfigKeys()
	if len(keys) != 0 {
		t.Fatalf("ListConfigKeys() failed to read config keys correctly")
	}

	testKeys := []string{"config.A", "config.B"}

	run.PutConfig(testKeys[0], "")
	run.PutConfig(testKeys[1], "")

	keys = run.ListConfigKeys()
	if len(keys) != 2 {
		t.Fatalf("ListConfigKeys() failed to read config keys correctly")
	}

	if !slices.Contains(keys, testKeys[0]) || !slices.Contains(keys, testKeys[1]) {
		t.Fatalf("ListClientNames() doesn't contain stored config key")
	}
}

func TestLoadConfig(t *testing.T) {
	run.ResetConfig()

	envVarKey := "some-test-key"
	envVarVal := "some-test-val"
	err := os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("LoadConfig() failed to test: %v", err)
	}

	_, err = run.LoadConfig("non-existent-key")
	if err == nil {
		t.Fatal("LoadConfig() didn't error on non-existent key")
	}

	result, err := run.LoadConfig(envVarKey)
	if err != nil {
		t.Fatalf("LoadConfig() failed to retrieve key: %v", err)
	}

	if result != envVarVal {
		t.Fatalf("LoadConfig() failed to retrieve key correctly: want: %s, have: %s", envVarVal, result)
	}

	if run.CountConfig() != 1 {
		t.Fatal("LoadConfig() failed to store config")
	}
}
