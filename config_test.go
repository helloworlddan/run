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
	"os"
	"slices"
	"testing"
)

func Test_putConfig(t *testing.T) {
	config := make(map[string]string)

	putConfig(config, "some key", "some val")

	if len(config) != 1 {
		t.Fatalf("putConfig() failed to add config correctly")
	}
}

func Test_getConfig(t *testing.T) {
	config := make(map[string]string)

	configKey := "test.config"
	configVal := "test.config.val"
	putConfig(config, configKey, configVal)

	_, err := getConfig(config, "non-existent")
	if err == nil {
		t.Fatalf("getConfig() failed to err on non-existent config")
	}

	result, err := getConfig(config, configKey)
	if err != nil {
		t.Fatalf("getConfig() failed to retrieve existing config")
	}

	if result != configVal {
		t.Fatalf("getConfig() failed to store config correctly")
	}
}

func Test_listConfigKeys(t *testing.T) {
	config := make(map[string]string)

	keys := listConfigKeys(config)
	if len(keys) != 0 {
		t.Fatalf("listConfigKeys() failed to read config keys correctly")
	}

	testKeys := []string{"config.A", "config.B"}

	putConfig(config, testKeys[0], "")
	putConfig(config, testKeys[1], "")

	keys = listConfigKeys(config)
	if len(keys) != 2 {
		t.Fatalf("listConfigKeys() failed to read config keys correctly")
	}

	if !slices.Contains(keys, testKeys[0]) || !slices.Contains(keys, testKeys[1]) {
		t.Fatalf("listClientNames() doesn't contain stored config key")
	}
}

func Test_loadConfig(t *testing.T) {
	config := make(map[string]string)

	envVarKey := "some-test-key"
	envVarVal := "some-test-val"
	err := os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("loadConfig() failed to test: %v", err)
	}

	_, err = loadConfig(config, "non-existent-key")
	if err == nil {
		t.Fatal("loadConfig() didn't error on non-existent key")
	}

	result, err := loadConfig(config, envVarKey)
	if err != nil {
		t.Fatalf("loadConfig() failed to retrieve key: %v", err)
	}

	if result != envVarVal {
		t.Fatalf("loadConfig() failed to retrieve key correctly: want: %s, have: %s", envVarVal, result)
	}

	if len(config) != 1 {
		t.Fatal("loadConfig() failed to store config")
	}

	if config[envVarKey] != envVarVal {
		t.Fatal("loadConfig() stored config incorrectly")
	}
}
