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
	resetClients()

	AddClient("some key", nil)

	if len(clients) != 1 {
		t.Fatalf("AddCleint() failed to add client correctly")
	}
}

func TestGetClient(t *testing.T) {
	resetClients()

	clientName := "test.client"
	client := http.DefaultClient
	AddClient(clientName, client)

	_, err := GetClient("non-existent")
	if err == nil {
		t.Fatalf("GetClient() failed to err on non-existent client")
	}

	rawResult, err := GetClient(clientName)
	if err != nil {
		t.Fatalf("GetClient() failed to retrieve existing client")
	}

	result := rawResult.(*http.Client)
	if result != client {
		t.Fatalf("GetClient() failed to store client correctly")
	}
}

func TestListClientNames(t *testing.T) {
	resetClients()

	names := ListClientNames()
	if len(names) != 0 {
		t.Fatalf("ListClientNames() failed to read client names correctly")
	}

	testNames := []string{"client.A", "client.B"}

	AddClient(testNames[0], nil)
	AddClient(testNames[1], nil)

	names = ListClientNames()
	if len(names) != 2 {
		t.Fatalf("ListClientNames() failed to read client names correctly")
	}

	if !slices.Contains(names, testNames[0]) || !slices.Contains(names, testNames[1]) {
		t.Fatalf("ListClientNames() doesn't contain stored client name")
	}
}
