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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	knative "knative.dev/serving/pkg/apis/serving/v1"
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
	servicePort         string
	url                 string
	jobTaskIndex        int
	jobTaskAttempt      int
	jobTaskCount        int
}

var (
	this               instance // NOTE: acts as cache
	knativeService     *knative.Service
	knativeServiceOnce sync.Once
)

// ResetInstance resets the cached metadata of this instance
func ResetInstance() {
	this = instance{}
}

// ID returns the unique instance ID of the Cloud Run instance serving the
// running job or service by referring to the metadata server.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// simply return `000000`.
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

// Local checks the presence of Cloud Run provided variables to determine if
// this instance is running locally or on Cloud Run.
func Local() bool {
	if len(os.Getenv("K_SERVICE")) != 0 {
		return false
	}
	if len(os.Getenv("CLOUD_RUN_JOB")) != 0 {
		return false
	}
	return true
}

// Name returns a preferred name for the currently running Cloud Run service or
// job. This will be either the service or job name or simply 'local'.
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

// ServiceName returns the name of the currently running Cloud Run service by
// looking up the `K_SERVICE` environment variable.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// simply return `local`.
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

// JobName returns the name of the currently running Cloud Run job by
// looking up the `CLOUD_RUN_JOB` environment variable.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// simply return `local`.
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

// ServiceRevision returns the revision identifier of the currently running
// Cloud Run service by looking up the `K_REVISION` environment variable.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return a deterministic identifier in the form of `<SERVICE_NAME>-00001-xxx`.
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

// JobExecution returns the execution identifier of the currently running
// Cloud Run job by looking up the `CLOUD_RUN_EXECUTION` environment variable.
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

// ProjectID attempts to resolve the alphanumeric Google Cloud project ID that
// is hosting the current Cloud Run instance.
//
// It loosely does so by looking up the following established precedence:
// - The environment variable `GOOGLE_CLOUD_PROJECT`
// - Querying the metadata server
// - Simply returning `local`
func ProjectID() string {
	if len(this.projectID) >= 6 { // ProjectID should be at least 6 chars
		return this.projectID
	}
	this.projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if len(this.projectID) >= 6 {
		return this.projectID
	}
	project, _ := metadata("project/project-id")
	this.projectID = project
	if len(this.projectID) >= 6 {
		return this.projectID
	}
	this.projectID = "local"
	return this.projectID
}

// ProjectNumber looks up the numeric project number of the current Google Cloud
// project hosting the Cloud Run instance.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return `000000000000`.
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

// Region looks up the actively serving region for this Cloud Run service.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return `local`.
func Region() string {
	if this.region != "" {
		return this.region
	}
	region, err := metadata("instance/region")
	regionComponents := strings.Split(region, "/")
	this.region = regionComponents[len(regionComponents)-1]
	if err != nil {
		this.region = "local"
	}
	return this.region
}

// ServiceAccountEmail looks up and returns the email of the service account
// configured for this Cloud Run instance.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return `local@localhost.com`.
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

// URL infers the URL with which this service will be addressable. This will
// either be 'http://localhost:8080' or the deterministic URL provided by Cloud
// Run.
func URL() string {
	if this.url != "" {
		return this.url
	}

	url := "http://localhost:8080"
	region := Region()
	if region != "local" {
		url = fmt.Sprintf("https://%s-%s.%s.run.app", ServiceName(), ProjectNumber(), region)
	}
	this.url = url
	return this.url
}

// ServiceAccountAccessToken looks up and returns a fresh OAuth2 access token
// for the service account configured for this Cloud Run instance.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return `local-access-token`.
func ServiceAccountAccessToken() string {
	token, err := metadata("instance/service-accounts/default/token")
	if err != nil {
		return "local-access-token"
	}
	return token
}

// AddOAuth2Header injects an `Authorization` header  with a valid access token
// for the configured service account into the supplied HTTP request and returns
// it.
func AddOAuth2Header(request *http.Request) *http.Request {
	token := ServiceAccountAccessToken()
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return request
}

// ServiceAccountIdentityToken attempts to mint an OIDC Identity Token for the
// specified `audience` using the metadata server.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return `local-identity-token`.
func ServiceAccountIdentityToken(audience string) string {
	token, err := metadata(fmt.Sprintf(
		"instance/service-accounts/default/identity?audience=%s",
		audience,
	))
	if err != nil {
		return "local-identity-token"
	}
	return token
}

// AddOIDCHeader injects an `Authorization` header  with a valid identity token
// for the configured service account into the supplied HTTP request and returns
// it.
func AddOIDCHeader(request *http.Request, audience string) *http.Request {
	token := ServiceAccountIdentityToken(audience)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return request
}

// ServicePort looks up and returns the configured service `$PORT` for
// this Cloud Run service.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return the default value `8080`.
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

// JobTaskIndex looks up and returns the current task index for the running
// Cloud Run job.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// sumply return `-1`.
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

// JobTaskAttempt looks up and returns the current task attempt for the running
// Cloud Run job.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// sumply return `-1`.
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

// JobTaskCount looks up and returns the current task count for the running
// Cloud Run job.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// sumply return `-1`.
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

// KNativeService loads and returns a KNative Serving representation of the
// current service. Requires at least roles/run.Viewer on itself.
func KNativeService() (knative.Service, error) {
	// TODO : test this
	var err error
	if *knativeService == nil {
		knativeServiceOnce.Do(func() {
			err = loadKNativeService()
		})
	}

	if err != nil {
		return *knative.Service{}, nil
	}

	return *knativeService, nil
}

func loadKNativeService() error {
	if Local() {
		return errors.New("skipping kNative endpoints, assuming local")
	}

	url := fmt.Sprintf(
		"https://%s-run.googleapis.com/apis/serving.knative.dev/v1/namespaces/%s/routes/%s",
		Region(),
		ProjectID(),
		ServiceName(),
	)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request = AddOAuth2Header(request)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var service knative.Service
	err = json.Unmarshal(content, &service)
	if err != nil {
		return err
	}

	knativeService = &service // Setting global
	return nil
}

func metadata(path string) (string, error) {
	if Local() {
		return "", errors.New("skipping GCE metadata server, assuming local")
	}

	url := fmt.Sprintf("http://metadata.google.internal/computeMetadata/v1/%s", path)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Metadata-Flavor", "Google")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}
