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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

func projectID() (string, error) {
	project, err := metadata("project/project-id")
	if err == nil {
		return project, nil
	}

	project = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if len(project) >= 6 { // ProjectID should be at least 6 chars
		return project, nil
	}

	return "", errors.New("unable to read project ID")
}

func projectNumber() (string, error) {
	project, err := metadata("project/numeric-project-id")
	if err != nil {
		return "", err
	}

	return project, err
}

func region() (string, error) {
	region, err := metadata("instance/region")
	if err != nil {
		return "", err
	}

	return region, err
}

func instanceID() (string, error) {
	id, err := metadata("instance/id")
	if err != nil {
		return "", err
	}

	return id, err
}

func serviceAccountEmail() (string, error) {
	email, err := metadata("instance/service-accounts/default/email")
	if err == nil {
		return email, err
	}

	return "", err
}

func serviceAccountToken() (string, error) {
	token, err := metadata("instance/service-accounts/default/token")
	if err == nil {
		return token, err
	}

	return "", err
}

func port() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", errors.New("unable to read PORT")
	}
	return port, nil
}

func kNativeService() (string, error) {
	name := os.Getenv("K_SERVICE")
	if name == "" {
		return "", errors.New("unable to read KNATIVE_SERVICE")
	}
	return name, nil
}

func kNativeRevision() (string, error) {
	revision := os.Getenv("K_REVISION")
	if revision == "" {
		return "", errors.New("unable to read KNATIVE_REVISION")
	}
	return revision, nil
}

func jobName() (string, error) {
	job := os.Getenv("CLOUD_RUN_JOB")
	if job == "" {
		return "", errors.New("unable to read CLOUD_RUN_JOB")
	}
	return job, nil
}

func jobExecution() (string, error) {
	execution := os.Getenv("CLOUD_RUN_EXECUTION")
	if execution == "" {
		return "", errors.New("unable to read CLOUD_RUN_EXECUTION")
	}
	return execution, nil
}

func jobTaskIndex() (int, error) {
	index, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_INDEX"))
	if err != nil {
		return 0, errors.New("unable to read CLOUD_RUN_TASK_INDEX")
	}
	return index, nil
}

func jobTaskAttempt() (int, error) {
	attempt, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_ATTEMPT"))
	if err != nil {
		return 0, errors.New("unable to read CLOUD_RUN_TASK_ATTEMPT")
	}
	return attempt, nil
}

func jobTaskCount() (int, error) {
	count, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_COUNT"))
	if err != nil {
		return 0, errors.New("unable to read CLOUD_RUN_TASK_COUNT")
	}
	return count, nil
}

func getConfig(config map[string]string, key string) (string, error) {
	val, ok := config[key]
	if !ok {
		return "", fmt.Errorf("no config found for key: '%s'", key)
	}
	return val, nil
}

func putConfig(config map[string]string, key string, val string) {
	config[key] = val
}

func loadConfig(config map[string]string, env string) (string, error) {
	val := os.Getenv(env)
	if val == "" {
		return "", fmt.Errorf("unable to find value for env var: '%s'", env)
	}

	config[env] = val

	return val, nil
}

func listConfigKeys(config map[string]string) []string {
	keys := make([]string, 0, len(config))
	for key := range config {
		keys = append(keys, key)
	}
	return keys
}
