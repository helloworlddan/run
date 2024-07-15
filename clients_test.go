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

package run_test

import (
	"slices"
	"testing"

	"github.com/helloworlddan/run"
)

type FakeClient struct {
	id          string
	initialized bool
}

func TestStoreClient(t *testing.T) {
	run.ResetClients()

	run.StoreClient("some key", nil, nil)

	if run.CountClients() != 1 {
		t.Fatalf("StoreClient() failed to add client correctly")
	}
}

// BUG: Possible concurrent access issue
func TestUseClient(t *testing.T) {
	run.ResetClients()

	var storedClient *FakeClient
	run.StoreClient("fake", storedClient, func() {
		storedClient = &FakeClient{
			id:          "fake",
			initialized: true,
		}
		run.InitializeClient("fake", &FakeClient{
			id:          "fake",
			initialized: true,
		})
	})

	var requestedClient *FakeClient
	requestedClient, err := run.UseClient("fake", requestedClient)
	if err != nil {
		t.Fatal("failed to retrieve client")
	}

	if requestedClient == nil {
		t.Fatal("failed to retrieve client")
	}

	if !requestedClient.initialized {
		t.Fatal("failed to initialize client")
	}
}

func TestListClientNames(t *testing.T) {
	run.ResetClients()

	names := run.ListClientNames()
	if len(names) != 0 {
		t.Fatalf("ListClientNames() failed to read client names correctly")
	}

	testNames := []string{"client.A", "client.B"}

	run.StoreClient(testNames[0], nil, nil)
	run.StoreClient(testNames[1], nil, nil)

	names = run.ListClientNames()
	if len(names) != 2 {
		t.Fatalf("ListClientNames() failed to read client names correctly")
	}

	if !slices.Contains(names, testNames[0]) || !slices.Contains(names, testNames[1]) {
		t.Fatalf("ListClientNames() doesn't contain stored client name")
	}
}
