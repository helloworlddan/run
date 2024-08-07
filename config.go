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
	"strings"
)

var config map[string]string

// ResetConfig deletes all previously configured config.
func ResetConfig() {
	config = make(map[string]string)
}

// CountConfig returns number of stored config elements.
func CountConfig() int {
	return len(config)
}

// PutConfig adds a K/V pair to the global config store.
func PutConfig(key string, val string) {
	ensureInitConfig()
	config[key] = val
}

// GetConfig retrieves a value for a key from the global config store.
func GetConfig(key string) (string, error) {
	ensureInitConfig()
	val, ok := config[key]
	if !ok {
		return "", fmt.Errorf("no config found for key: '%s'", key)
	}
	return val, nil
}

// ListConfigKeys returns a list of all currently available keys in the global
// config store.
func ListConfigKeys() []string {
	ensureInitConfig()
	keys := make([]string, 0, len(config))
	for key := range config {
		keys = append(keys, key)
	}
	return keys
}

// LoadConfig lookups the named environment variable, puts it's value into
// the global config store and returns the value.
func LoadConfig(env string) (string, error) {
	ensureInitConfig()
	val := os.Getenv(env)
	if val == "" {
		return "", fmt.Errorf("unable to find value for env var: '%s'", env)
	}

	PutConfig(env, val)

	return val, nil
}

// LoadAllConfig loads all available environment variables and puts it in the
// config store.
func LoadAllConfig() {
	ensureInitConfig()
	for _, pair := range os.Environ() {
		key, value, ok := strings.Cut(pair, "=")
		if ok {
			PutConfig(key, value)
		}
	}
}

func ensureInitConfig() {
	if config == nil {
		ResetConfig()
	}
}
