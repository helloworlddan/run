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
	"sync"
)

var clients map[string]*lazyClient

type lazyClient struct {
	clientPtr      any
	lazyInitialize func()
	clientOnce     sync.Once
}

func ensureInitClients() {
	if clients == nil {
		resetClients()
	}
}

// AddClient is intended to store a pointer to an initialized API client under a
// given key.
func AddClient(name string, client any) {
	ensureInitClients()
	clients[name] = &lazyClient{
		client,
		nil,
		sync.Once{},
	}
}

// AddClient is intended to store a pointer to an initialized API client under a
// given key.
func AddLazyClient(name string, client any, lazyInit func()) {
	ensureInitClients()
	clients[name] = &lazyClient{
		client,
		lazyInit,
		sync.Once{},
	}
}

// GetClient is intended to retrieve a pointer to an initialized API for a given
// key name.
func GetClient(name string) (any, error) {
	ensureInitClients()
	lc, ok := clients[name]
	if !ok {
		return nil, fmt.Errorf("no client found for name: '%s'", name)
	}

	if lc.lazyInitialize != nil {
		lc.clientOnce.Do(lc.lazyInitialize)
	}

	return lc.clientPtr, nil
}

// ListClientNames returns a list of all available keys store in the global
// clients store.
func ListClientNames() []string {
	ensureInitClients()
	names := make([]string, 0, len(clients))
	for name := range clients {
		names = append(names, name)
	}
	return names
}

func resetClients() {
	clients = make(map[string]*lazyClient)
}
