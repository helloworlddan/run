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

// Package run provides a simple way to use Google Cloud Run.
package run

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// metadata retrieves a value from the Google Cloud GCE Metadata service.
//
// The path argument should be a path relative to the root of the metadata service.
// For example, to retrieve the project ID, you would use the path "project/project-id".
func metadata(path string) (string, error) {
	path = fmt.Sprintf("http://metadata.google.internal/computeMetadata/v1/%s", path)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Metadata-Flavor", "Google")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(string(raw)), nil
}
