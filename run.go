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

	knative "knative.dev/serving/pkg/apis/serving/v1"
)

type RunResourceType int

const (
	LocalResource   RunResourceType = iota
	ServiceResource RunResourceType = iota
	JobResource     RunResourceType = iota
)

type cache struct {
	instanceID          string
	serviceName         string
	jobName             string
	serviceRevision     string
	jobExecution        string
	projectID           string
	projectNumber       string
	region              string
	serviceAccountEmail string
	servicePort         string
	serviceURL          string
	jobTaskIndex        int
	jobTaskAttempt      int
	jobTaskCount        int
	knativeService      *knative.Service
}

var this cache // NOTE: acts as cache

// ResetCache resets the cached metadata of this instance
func ResetCache() {
	this = cache{}
}

func ResourceType() RunResourceType {
	if ServiceName() != "local" {
		return ServiceResource
	}
	if JobName() != "local" {
		return JobResource
	}
	return LocalResource
}

// ID returns the unique instance ID of the Cloud Run instance serving the
// running job or service by referring to the metadata server.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// simply return `000000`.
func InstanceID() string {
	if this.instanceID != "" {
		return this.instanceID
	}
	id, err := metadata("instance/id")
	this.instanceID = id
	if err != nil {
		this.instanceID = "000000"
	}
	return this.instanceID
}

// Name returns a preferred name for the currently running Cloud Run service or
// job. This will be either the service or job name or simply 'local'.
func Name() string {
	if ResourceType() == ServiceResource {
		return ServiceName()
	}
	if ResourceType() == JobResource {
		return JobName()
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

// Revision returns the revision identifier of the currently running
// Cloud Run service by looking up the `K_REVISION` environment variable.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return a deterministic identifier in the form of `<SERVICE_NAME>-00001-xxx`.
func Revision() string {
	if this.serviceRevision != "" {
		return this.serviceRevision
	}
	this.serviceRevision = os.Getenv("K_REVISION")
	if this.serviceRevision == "" {
		this.serviceRevision = fmt.Sprintf("%s-00001-xxx", Name())
	}
	return this.serviceRevision
}

// Execution returns the execution identifier of the currently running
// Cloud Run job by looking up the `CLOUD_RUN_EXECUTION` environment variable.
func Execution() string {
	if this.jobExecution != "" {
		return this.jobExecution
	}
	this.jobExecution = os.Getenv("CLOUD_RUN_EXECUTION")
	if this.jobExecution == "" {
		this.jobExecution = "local"
	}
	return this.jobExecution
}

// Port looks up and returns the configured service `$PORT` for
// this Cloud Run service.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// return the default value `8080`.
func Port() string {
	if this.servicePort != "" {
		return this.servicePort
	}
	this.servicePort = os.Getenv("PORT")
	if this.servicePort == "" {
		this.servicePort = "8080"
	}
	return this.servicePort
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
// Run
func DefaultServiceURL() string {
	if this.serviceURL != "" {
		return this.serviceURL
	}

	url := "http://localhost:8080"
	region := Region()
	if region != "local" {
		url = fmt.Sprintf("https://%s-%s.%s.run.app", ServiceName(), ProjectNumber(), region)
	}
	this.serviceURL = url
	return this.serviceURL
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
func AddOAuth2Header(r *http.Request) *http.Request {
	token := ServiceAccountAccessToken()
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return r
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
func AddOIDCHeader(r *http.Request, audience string) *http.Request {
	token := ServiceAccountIdentityToken(audience)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return r
}

// TaskIndex looks up and returns the current task index for the running
// Cloud Run job.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// sumply return `-1`.
func TaskIndex() int {
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

// TaskAttempt looks up and returns the current task attempt for the running
// Cloud Run job.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// sumply return `-1`.
func TaskAttempt() int {
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

// TaskCount looks up and returns the current task count for the running
// Cloud Run job.
//
// If the current process does not seem to be hosted on Cloud Run, it will
// sumply return `-1`.
func TaskCount() int {
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
	if this.knativeService != nil {
		return *this.knativeService, nil
	}

	err := loadKNativeService()
	if err != nil {
		return knative.Service{}, err
	}

	return *this.knativeService, nil
}

func Creator() (string, error) {
	annotationKey := "serving.knative.dev/creator"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func LastModifier() (string, error) {
	annotationKey := "serving.knative.dev/lastModifier"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func LaunchStage() (string, error) {
	annotationKey := "run.googleapis.com/launch-stage"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func Description() (string, error) {
	annotationKey := "run.googleapis.com/description"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func Ingress() (string, error) {
	annotationKey := "run.googleapis.com/ingress"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func BinaryAuthorizationPolicy() (string, error) {
	annotationKey := "run.googleapis.com/binary-authorization"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

// BinaryAuthorizationBreakglassJustification returns the justification for
// circumventing the configured Binary Authorization policy.
func BinaryAuthorizationBreakglassJustification() (string, error) {
	annotationKey := "run.googleapis.com/binary-authorization-breakglass"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func ServiceMinimumInstances() (int, error) {
	annotationKey := "run.googleapis.com/minScale"
	knativeService, err := KNativeService()
	if err != nil {
		return 0, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return 0, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strconv.Atoi(annotationValue)
}

func FunctionEntryPoint() (string, error) {
	annotationKey := "run.googleapis.com/function-target"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func InvokerIAMDisabled() (bool, error) {
	annotationKey := "run.googleapis.com/invoker-iam-disabled"
	knativeService, err := KNativeService()
	if err != nil {
		return false, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return false, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strconv.ParseBool(annotationValue)
}

func IAPEnabled() (bool, error) {
	annotationKey := "run.googleapis.com/iap-enabled"
	knativeService, err := KNativeService()
	if err != nil {
		return false, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return false, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strconv.ParseBool(annotationValue)
}

func ScalingMode() (string, error) {
	annotationKey := "run.googleapis.com/scalingMode"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func ManualInstances() (int, error) {
	annotationKey := "run.googleapis.com/manualInstanceCount"
	knativeService, err := KNativeService()
	if err != nil {
		return -1, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Annotations[annotationKey]
	if annotationValue == "" {
		return -1, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strconv.Atoi(annotationValue)
}

func RevisionMinimumInstances() (int, error) {
	annotationKey := "autoscaling.knative.dev/minScale"
	knativeService, err := KNativeService()
	if err != nil {
		return 0, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return 0, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strconv.Atoi(annotationValue)
}

func RevisionMaximumInstances() (int, error) {
	annotationKey := "autoscaling.knative.dev/maxScale"
	knativeService, err := KNativeService()
	if err != nil {
		return 0, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return 0, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strconv.Atoi(annotationValue)
}

func CPUThrottling() (bool, error) {
	annotationKey := "run.googleapis.com/cpu-throttling"
	knativeService, err := KNativeService()
	if err != nil {
		return true, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return true, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strconv.ParseBool(annotationValue)
}

func StartupCPUBoost() (bool, error) {
	annotationKey := "run.googleapis.com/cpu-throttling"
	knativeService, err := KNativeService()
	if err != nil {
		return false, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return false, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strconv.ParseBool(annotationValue)
}

func SessionAffinity() (string, error) {
	annotationKey := "run.googleapis.com/SessionAffinity"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func CloudSQLInstances() ([]string, error) {
	annotationKey := "run.googleapis.com/cloudsql-instances"
	knativeService, err := KNativeService()
	if err != nil {
		return []string{}, fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return []string{}, fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return strings.Split(annotationValue, ","), nil
}

func ExecutionEnvironment() (string, error) {
	annotationKey := "run.googleapis.com/execution-environment"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func VPCAccessConnector() (string, error) {
	annotationKey := "run.googleapis.com/vpc-access-connector"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func VPCAccessEgress() (string, error) {
	annotationKey := "run.googleapis.com/vpc-access-egress"
	knativeService, err := KNativeService()
	if err != nil {
		return "all-traffic", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return "all-traffic", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	if annotationValue == "all" {
		return "all-traffic", nil
	}

	return annotationValue, nil
}

func VPCNetworkInterfaces() (string, error) {
	annotationKey := "run.googleapis.com/network-interfaces"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func EncryptionKey() (string, error) {
	annotationKey := "run.googleapis.com/encryption-key"
	knativeService, err := KNativeService()
	if err != nil {
		return "", fmt.Errorf("error loading property '%s': %v", annotationKey, err)
	}
	annotationValue := knativeService.Spec.Template.Annotations[annotationKey]
	if annotationValue == "" {
		return "", fmt.Errorf("error reading property '%s'", annotationKey)
	}

	return annotationValue, nil
}

func loadKNativeService() error {
	if ResourceType() != ServiceResource {
		return errors.New("skipping KNative endpoint, assuming local")
	}

	url := fmt.Sprintf(
		"https://%s-run.googleapis.com/apis/serving.knative.dev/v1/namespaces/%s/services/%s",
		Region(),
		ProjectID(),
		ServiceName(),
	)
	Debugf(nil, "requesting: %s", url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request = AddOAuth2Header(request)
	Debugf(nil, "authenticating request: %#v", request)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf(
			"failed to read KNative endpoints, do we have roles/run.viewer on our service? error: %v",
			err,
		)
	}
	defer resp.Body.Close()

	Debugf(nil, "status: %s", resp.Status)

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	Debugf(nil, "raw response: %s", string(content))

	var service knative.Service
	err = json.Unmarshal(content, &service)
	if err != nil {
		Debugf(nil, "failed to unmarshal: %v", err)
		return err
	}

	this.knativeService = &service // Setting global
	return nil
}

func metadata(path string) (string, error) {
	if ResourceType() == LocalResource {
		return "", errors.New("skipping GCE metadata server, assuming local")
	}

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
