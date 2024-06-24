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

func TestPutConfig(t *testing.T) {
	config = make(map[string]string)

	PutConfig("some key", "some val")

	if len(config) != 1 {
		t.Fatalf("putConfig() failed to add config correctly")
	}
	config = make(map[string]string)
}

func TestGetConfig(t *testing.T) {
	config = make(map[string]string)

	configKey := "test.config"
	configVal := "test.config.val"
	PutConfig(configKey, configVal)

	_, err := GetConfig("non-existent")
	if err == nil {
		t.Fatalf("getConfig() failed to err on non-existent config")
	}

	result, err := GetConfig(configKey)
	if err != nil {
		t.Fatalf("getConfig() failed to retrieve existing config")
	}

	if result != configVal {
		t.Fatalf("getConfig() failed to store config correctly")
	}
	config = make(map[string]string)
}

func TestListConfigKeys(t *testing.T) {
	config = make(map[string]string)

	keys := ListConfigKeys()
	if len(keys) != 0 {
		t.Fatalf("listConfigKeys() failed to read config keys correctly")
	}

	testKeys := []string{"config.A", "config.B"}

	PutConfig(testKeys[0], "")
	PutConfig(testKeys[1], "")

	keys = ListConfigKeys()
	if len(keys) != 2 {
		t.Fatalf("listConfigKeys() failed to read config keys correctly")
	}

	if !slices.Contains(keys, testKeys[0]) || !slices.Contains(keys, testKeys[1]) {
		t.Fatalf("listClientNames() doesn't contain stored config key")
	}
	config = make(map[string]string)
}

func TestLoadConfig(t *testing.T) {
	config = make(map[string]string)

	envVarKey := "some-test-key"
	envVarVal := "some-test-val"
	err := os.Setenv(envVarKey, envVarVal)
	if err != nil {
		t.Fatalf("loadConfig() failed to test: %v", err)
	}

	_, err = LoadConfig("non-existent-key")
	if err == nil {
		t.Fatal("loadConfig() didn't error on non-existent key")
	}

	result, err := LoadConfig(envVarKey)
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
	config = make(map[string]string)
}
