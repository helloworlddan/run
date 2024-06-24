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
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func NewAuthenticatedRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	token := ServiceAccountToken()
	req.Header.Add("Authorization", fmt.Sprintf("bearer: %s", token))

	return req, nil
}

func ProjectID() string {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if len(project) >= 6 { // ProjectID should be at least 6 chars
		return project
	}

	project, err := metadata("project/project-id")
	if err != nil {
		return "local"
	}
	return project
}

func ProjectNumber() string {
	number, err := metadata("project/numeric-project-id")
	if err != nil {
		return "000000000000"
	}
	return number
}

func Region() string {
	region, err := metadata("instance/region")
	if err != nil {
		return "local"
	}
	return region
}

func ID() string {
	id, err := metadata("instance/id")
	if err != nil {
		return "00000"
	}
	return id
}

func ServiceAccountEmail() string {
	email, err := metadata("instance/service-accounts/default/email")
	if err != nil {
		return "local"
	}
	return email
}

func ServiceAccountToken() string {
	token, err := metadata("instance/service-accounts/default/token")
	if err != nil {
		return "local"
	}
	return token
}

func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func Name() string {
	name := KNativeService()
	if name != "local" {
		return name
	}

	name = JobName()
	if name != "local" {
		return name
	}

	return "local"
}

func KNativeService() string {
	name := os.Getenv("K_SERVICE")
	if name == "" {
		return "local"
	}
	return name
}

func KNativeRevision() string {
	revision := os.Getenv("K_REVISION")
	if revision == "" {
		return fmt.Sprintf("%s-00001-xxx", Name())
	}
	return revision
}

func JobName() string {
	job := os.Getenv("CLOUD_RUN_JOB")
	if job == "" {
		return "local"
	}
	return job
}

func JobExecution() string {
	execution := os.Getenv("CLOUD_RUN_EXECUTION")
	if execution == "" {
		return "local"
	}
	return execution
}

func JobTaskIndex() int {
	index, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_INDEX"))
	if err != nil {
		return 0
	}
	return index
}

func JobTaskAttempt() int {
	attempt, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_ATTEMPT"))
	if err != nil {
		return 0
	}
	return attempt
}

func JobTaskCount() int {
	count, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_COUNT"))
	if err != nil {
		return 0
	}
	return count
}

func metadata(path string) (string, error) {
	path = fmt.Sprintf("http://metadata.google.internal/computeMetadata/v1/%s", path)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Metadata-Flavor", "Google")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(raw)), nil
}
