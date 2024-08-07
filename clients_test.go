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
	"net/http"
	"slices"
	"testing"

	"github.com/helloworlddan/run"
)

type fakeClient struct {
	id          string
	initialized bool
}

func TestClient(t *testing.T) {
	run.ResetClients()

	run.Client("some key", "some client")

	if run.CountClients() != 1 {
		t.Fatalf("StoreClient() failed to add client correctly")
	}
}

func TestUseClient(t *testing.T) {
	run.ResetClients()

	run.Client("http", http.DefaultClient)

	var httpClient *http.Client
	httpClient, err := run.UseClient("http", httpClient)
	if err != nil {
		t.Fatalf("failed to retrieve client, err: %v", err)
	}

	if httpClient == nil {
		t.Fatal("failed to retrieve client")
	}

	run.LazyClient("fake", func() {
		run.Client("fake", &fakeClient{
			id:          "fake",
			initialized: true,
		})
	})

	var client *fakeClient
	client, err = run.UseClient("fake", client)
	if err != nil {
		t.Fatalf("failed to retrieve client, err: %v", err)
	}

	if client == nil {
		t.Fatal("failed to retrieve client")
	}

	if !client.initialized {
		t.Fatal("failed to retrieve client")
	}
}

func TestListClientNames(t *testing.T) {
	run.ResetClients()

	names := run.ListClientNames()
	if len(names) != 0 {
		t.Fatalf("ListClientNames() failed to read client names correctly")
	}

	testNames := []string{"client.A", "client.B"}

	run.Client(testNames[0], testNames[0])
	run.Client(testNames[1], testNames[1])

	names = run.ListClientNames()
	if len(names) != 2 {
		t.Fatalf("ListClientNames() failed to read client names correctly")
	}

	if !slices.Contains(names, testNames[0]) || !slices.Contains(names, testNames[1]) {
		t.Fatalf("ListClientNames() doesn't contain stored client name")
	}
}
