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
	"reflect"
	"slices"
	"sync"
)

var clients map[string]*lazyClient

type lazyClient struct {
	clientPtr      any
	clientType     reflect.Type
	lazyInitialize func()
	clientOnce     sync.Once
}

func ensureInitClients() {
	if clients == nil {
		resetClients()
	}
}

// StoreClient is intended to store a pointer to a client under a given key.
// The caller can optionally specify lazyInit to be run once before the client
// is retrieved using run.UseClient for the first time. If lazyInit is nil, the
// client should be initialized before storing it.
func StoreClient(name string, client any, lazyInit func()) {
	ensureInitClients()
	clients[name] = &lazyClient{
		client,
		reflect.TypeOf(client),
		lazyInit,
		sync.Once{},
	}
}

// GetClient is intended to retrieve a pointer to a client for a given key name.
// I requires the name of a stored client and a nil pointer of it's type.
func UseClient[T any](name string, client T) (T, error) {
	ensureInitClients()
	lc, ok := clients[name]
	if !ok {
		return client, fmt.Errorf("no client found for name: '%s'", name)
	}

	if reflect.TypeOf(client) != lc.clientType {
		return client, fmt.Errorf("wrong type requested for client name: '%s'", name)
	}

	// Do lazy initialization
	if lc.lazyInitialize != nil {
		lc.clientOnce.Do(lc.lazyInitialize)
	}

	v := reflect.ValueOf(lc.clientPtr)
	actual := v.Interface().(T)

	return actual, nil
}

// TODO: build client in use for Close()

// ListClientNames returns a list of all available keys store in the global
// clients store.
func ListClientNames() []string {
	ensureInitClients()
	names := make([]string, 0, len(clients))
	for name := range clients {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func resetClients() {
	clients = make(map[string]*lazyClient)
}
