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
	lazyInitialize func()
	clientOnce     *sync.Once
	initialized    bool
}

func ensureInitClients() {
	if clients == nil {
		ResetClients()
	}
}

// ResetClients deletes all previously configured clients.
func ResetClients() {
	clients = make(map[string]*lazyClient)
}

// CountClients returns number of stored clients.
func CountClients() int {
	return len(clients)
}

// Client registers and already initialized clients
func Client(name string, client any) {
	ensureInitClients()
	clients[name] = &lazyClient{
		client,
		nil,
		&sync.Once{},
		true,
	}
}

// LazyClient registers an uninitialized client name with an initialization
// function. The init func should call Client() with the initialized client.
func LazyClient(name string, init func()) {
	ensureInitClients()
	clients[name] = &lazyClient{
		nil,
		init,
		&sync.Once{},
		false,
	}
}

// GetClient is intended to retrieve a pointer to a client for a given key name.
// I requires the name of a stored client and a nil pointer of it's type.
func UseClient[T any](name string, client T) (T, error) {
	ensureInitClients()
	// Check if client is a pointer
	if !isPointer(client) {
		return client, fmt.Errorf("expected pointer to client, but got %T", client)
	}

	// Check if client is known
	lc, ok := clients[name]
	if !ok {
		return client, fmt.Errorf("no client found for name: '%s'", name)
	}

	if !lc.initialized && lc.lazyInitialize == nil {
		return client, fmt.Errorf("cannot initialize client '%s'", name)
	}

	lc.clientOnce.Do(lc.lazyInitialize)

	// refresh
	lc, ok = clients[name]
	if !ok {
		return client, fmt.Errorf("no client found for name: '%s'", name)
	}

	actual, ok := lc.clientPtr.(T)
	if !ok {
		return client, fmt.Errorf("failed to cast stored client to requested type: %T", actual)
	}

	client = actual

	return actual, nil
}

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

func isPointer(a any) bool {
	return reflect.ValueOf(a).Kind() == reflect.Ptr
}
