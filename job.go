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
	"log"
)

// Job is intended to be instantiated once and kept around to access
// functionality related to the Cloud Run Job runtime.
type Job struct {
	configs map[string]string
	clients map[string]interface{}
}

// NewJob creates a new Job instance.
func NewJob() *Job {
	log.SetFlags(0)
	j := &Job{
		configs: make(map[string]string),
		clients: make(map[string]interface{}),
	}

	return j
}

// GetConfig retrieves a config value from the store
func (j *Job) GetConfig(key string) (string, error) {
	return getConfig(j.configs, key)
}

// PutConfig puts a config value in the store
func (j *Job) PutConfig(key string, val string) {
	putConfig(j.configs, key, val)
}

// LoadConfig looks up an environment variable puts it in the store and returns
// it's value
func (j *Job) LoadConfig(env string) (string, error) {
	return loadConfig(j.configs, env)
}

// ListConfigKeys returns a list of all available config keys
func (j *Job) ListConfig() []string {
	return listConfigKeys(j.configs)
}

// GetClient resolves a client by name from the store
func (j *Job) GetClient(name string) (any, error) {
	return getClient(j.clients, name)
}

// AddClient add a client to the store
func (j *Job) AddClient(name string, client any) {
	addClient(j.clients, name, client)
}

// ListClientNames returns a list of all available clients
func (j *Job) ListClientNames() []string {
	return listClientNames(j.clients)
}
