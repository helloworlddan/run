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
	"net/http"
	"slices"
	"testing"
)

func TestAddClient(t *testing.T) {
	clients = make(map[string]any)
	AddClient("some key", nil)

	if len(clients) != 1 {
		t.Fatalf("addCleint() failed to add client correctly")
	}
	clients = make(map[string]any)
}

func TestGetClient(t *testing.T) {
	clients = make(map[string]any)

	clientName := "test.client"
	client := http.DefaultClient
	AddClient(clientName, client)

	_, err := GetClient("non-existent")
	if err == nil {
		t.Fatalf("getClient() failed to err on non-existent client")
	}

	rawResult, err := GetClient(clientName)
	if err != nil {
		t.Fatalf("getClient() failed to retrieve existing client")
	}

	result := rawResult.(*http.Client)
	if result != client {
		t.Fatalf("getClient() failed to store client correctly")
	}
	clients = make(map[string]any)
}

func TestListClientNames(t *testing.T) {
	clients = make(map[string]any)

	names := ListClientNames()
	if len(names) != 0 {
		t.Fatalf("listClientNames() failed to read client names correctly")
	}

	testNames := []string{"client.A", "client.B"}

	AddClient(testNames[0], nil)
	AddClient(testNames[1], nil)

	names = ListClientNames()
	if len(names) != 2 {
		t.Fatalf("listClientNames() failed to read client names correctly")
	}

	if !slices.Contains(names, testNames[0]) || !slices.Contains(names, testNames[1]) {
		t.Fatalf("listClientNames() doesn't contain stored client name")
	}
	clients = make(map[string]any)
}
