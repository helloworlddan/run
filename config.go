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
)

var config map[string]string

func ensureInitConfig() {
	if config == nil {
		resetConfig()
	}
}

func PutConfig(key string, val string) {
	ensureInitConfig()
	config[key] = val
}

func GetConfig(key string) (string, error) {
	ensureInitConfig()
	val, ok := config[key]
	if !ok {
		return "", fmt.Errorf("no config found for key: '%s'", key)
	}
	return val, nil
}

func ListConfigKeys() []string {
	ensureInitConfig()
	keys := make([]string, 0, len(config))
	for key := range config {
		keys = append(keys, key)
	}
	return keys
}

func LoadConfig(env string) (string, error) {
	ensureInitConfig()
	val := os.Getenv(env)
	if val == "" {
		return "", fmt.Errorf("unable to find value for env var: '%s'", env)
	}

	PutConfig(env, val)

	return val, nil
}

func resetConfig() {
	config = make(map[string]string)
}
