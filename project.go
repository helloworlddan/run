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
	"errors"
	"os"
)

// ProjectID returns the Google Cloud Project ID.
//
// It first checks the GOOGLE_CLOUD_PROJECT environment variable.
// If it is not set, it retrieves the project ID from the Google Cloud Metadata service.
func ProjectID() (string, error) {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project != "" {
		return project, nil
	}

	project, err := metadata("project/project-id")
	if err != nil {
		return "", err
	}

	if len(project) < 6 {
		return "", errors.New("unable to read project ID")
	}

	return project, nil
}

// ProjectNumber returns the Google Cloud Project Number.
//
// It retrieves the project number from the Google Cloud Metadata service.
func ProjectNumber() (string, error) {
	project, err := metadata("project/numeric-project-id")
	if err != nil {
		return "", err
	}

	return project, err
}
