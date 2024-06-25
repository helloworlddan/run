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

// Package Run provides useful, yet opinionated integrations for workloads
// running on Cloud Run.
package run

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type instance struct {
	id                  string
	serviceName         string
	jobName             string
	serviceRevision     string
	jobExecution        string
	projectID           string
	projectNumber       string
	region              string
	serviceAccountEmail string
	serviceAccountToken string
	servicePort         string
	jobTaskIndex        int
	jobTaskAttempt      int
	jobTaskCount        int
}

var this instance // NOTE: acts as cache

func Name() string {
	name := ServiceName()
	if name != "local" {
		return name
	}
	name = JobName()
	if name != "local" {
		return name
	}
	return "local"
}

func ServiceName() string {
	if this.serviceName != "" {
		return this.serviceName
	}
	this.serviceName = os.Getenv("K_SERVICE")
	if this.serviceName == "" {
		this.serviceName = "local"
	}
	return this.serviceName
}

func JobName() string {
	if this.jobName != "" {
		return this.jobName
	}
	this.jobName = os.Getenv("CLOUD_RUN_JOB")
	if this.jobName == "" {
		this.jobName = "local"
	}
	return this.jobName
}

func ID() string {
	if this.id != "" {
		return this.id
	}
	id, err := metadata("instance/id")
	this.id = id
	if err != nil {
		this.id = "000000"
	}
	return this.id
}

func ServiceRevision() string {
	if this.serviceRevision != "" {
		return this.serviceRevision
	}
	this.serviceRevision = os.Getenv("K_REVISION")
	if this.serviceRevision == "" {
		this.serviceRevision = fmt.Sprintf("%s-00001-xxx", Name())
	}
	return this.serviceRevision
}

func JobExecution() string {
	if this.jobExecution != "" {
		return this.jobExecution
	}
	this.jobExecution = os.Getenv("CLOUD_RUN_EXECUTION")
	if this.jobExecution == "" {
		this.jobExecution = "local"
	}
	return this.jobExecution
}

func ProjectID() string {
	if this.projectID != "" {
		return this.projectID
	}
	this.projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if len(this.projectID) >= 6 { // ProjectID should be at least 6 chars
		return this.projectID
	}
	project, err := metadata("project/project-id")
	this.projectID = project
	if err != nil {
		this.projectID = "local"
	}
	return this.projectID
}

func ProjectNumber() string {
	if this.projectNumber != "" {
		return this.projectNumber
	}
	number, err := metadata("project/numeric-project-id")
	this.projectNumber = number
	if err != nil {
		this.projectNumber = "000000000000"
	}
	return this.projectNumber
}

func Region() string {
	if this.region != "" {
		return this.region
	}
	region, err := metadata("instance/region")
	this.region = region
	if err != nil {
		this.region = "local"
	}
	return this.region
}

func ServiceAccountEmail() string {
	if this.serviceAccountEmail != "" {
		return this.serviceAccountEmail
	}
	email, err := metadata("instance/service-accounts/default/email")
	this.serviceAccountEmail = email
	if err != nil {
		this.serviceAccountEmail = "local@localhost.com"
	}
	return this.serviceAccountEmail
}

func ServiceAccountToken() string {
	if this.serviceAccountToken != "" {
		return this.serviceAccountToken
	}
	token, err := metadata("instance/service-accounts/default/token")
	this.serviceAccountToken = token
	if err != nil {
		this.serviceAccountToken = "local-token"
	}
	return this.serviceAccountToken
}

func ServicePort() string {
	if this.servicePort != "" {
		return this.servicePort
	}
	this.servicePort = os.Getenv("PORT")
	if this.servicePort == "" {
		this.servicePort = "8080"
	}
	return this.servicePort
}

func JobTaskIndex() int {
	if this.jobTaskIndex != 0 {
		return this.jobTaskIndex
	}
	index, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_INDEX"))
	this.jobTaskIndex = index
	if err != nil {
		this.jobTaskIndex = -1
	}
	return this.jobTaskIndex
}

func JobTaskAttempt() int {
	if this.jobTaskAttempt != 0 {
		return this.jobTaskAttempt
	}
	attempt, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_ATTEMPT"))
	this.jobTaskAttempt = attempt
	if err != nil {
		this.jobTaskAttempt = -1
	}
	return this.jobTaskAttempt
}

func JobTaskCount() int {
	if this.jobTaskCount != 0 {
		return this.jobTaskCount
	}
	count, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_COUNT"))
	this.jobTaskCount = count
	if err != nil {
		this.jobTaskCount = -1
	}
	return this.jobTaskCount
}

func AddAuthHeader(r *http.Request) *http.Request {
	token := ServiceAccountToken()
	r.Header.Add("Authorization", fmt.Sprintf("bearer: %s", token))
	return r
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

func resetInstance() {
	this = instance{}
}
