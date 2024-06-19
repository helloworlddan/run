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
	"io"
	"log"
	"net/http"
)

// Job is intended to be instantiated once and kept around to access
// functionality related to the Cloud Run Job runtime.
type Job struct {
	configs map[string]string
	clients map[string]interface{}
}

// NewJob creates a new Job instance.
func NewJob() *Job {
	log.SetFlags(0)
	j := &Job{
		configs: make(map[string]string),
		clients: make(map[string]interface{}),
	}

	return j
}

// ID returns the ID of the serving instance
func (j *Job) ID() string {
	id, err := instanceID()
	if err != nil {
		id = "00000"
	}
	return id
}

// Name returns the name of the job
func (j *Job) Name() string {
	name, err := jobName()
	if err != nil {
		name = "local"
	}
	return name
}

// String returns the name of the service to satisfy fmt.Stringer
func (j *Job) String() string {
	return j.Name()
}

// Execution returns the name of the current execution of the job
func (j *Job) Execution() string {
	execution, err := jobExecution()
	if err != nil {
		execution = "local"
	}
	return execution
}

// TaskIndex returns the task index assigned to the job
func (j *Job) TaskIndex() int {
	index, err := jobTaskIndex()
	if err != nil {
		index = 0
	}
	return index
}

// TaskAttempt returns the attempt/retry counter of the task
func (j *Job) TaskAttempt() int {
	attempt, err := jobTaskAttempt()
	if err != nil {
		attempt = 0
	}
	return attempt
}

// TaskCount returns the total count of tasks
func (j *Job) TaskCount() int {
	count, err := jobTaskCount()
	if err != nil {
		count = 0
	}
	return count
}

// ProjectID returns the name of the containing Google Cloud project or "local"
func (j *Job) ProjectID() string {
	project, err := projectID()
	if err != nil {
		project = "local"
	}
	return project
}

// ProjectNumber returns the 12-digit project number of the containing Google
// Cloud project or "000000000000"
func (j *Job) ProjectNumber() string {
	number, err := projectNumber()
	if err != nil {
		number = "000000000000"
	}
	return number
}

// Region returns the Google Cloud region in which the job is running or "local"
func (j *Job) Region() string {
	region, err := region()
	if err != nil {
		region = "local"
	}
	return region
}

// ServiceAccountEmail returns the email of the service account assigned to the
// job
func (j *Job) ServiceAccountEmail() string {
	email, err := serviceAccountEmail()
	if err != nil {
		email = "local"
	}
	return email
}

// ServiceAccountToken returns an authentication token for the assigned service
// account to authorize requests.
func (j *Job) ServiceAccountToken() string {
	token, err := serviceAccountToken()
	if err != nil {
		token = "local"
	}
	return token
}

// NewAuthenticatedRequest returns a new http request with an Authorization header
func (j *Job) NewAuthenticatedRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	return newAuthenticatedRequest(j, ctx, method, url, body)
}

// GetConfig retrieves a config value from the store
func (j *Job) GetConfig(key string) (string, error) {
	return getConfig(j.configs, key)
}

// PutConfig puts a config value in the store
func (j *Job) PutConfig(key string, val string) {
	putConfig(j.configs, key, val)
}

// LoadConfig looks up an environment variable puts it in the store and returns
// it's value
func (j *Job) LoadConfig(env string) (string, error) {
	return loadConfig(j.configs, env)
}

// ListConfigKeys returns a list of all available config keys
func (j *Job) ListConfig() []string {
	return listConfigKeys(j.configs)
}

// GetClient resolves a client by name from the store
func (j *Job) GetClient(name string) (any, error) {
	return getClient(j.clients, name)
}

// AddClient add a client to the store
func (j *Job) AddClient(name string, client any) {
	addClient(j.clients, name, client)
}

// ListClientNames returns a list of all available clients
func (j *Job) ListClientNames() []string {
	return listClientNames(j.clients)
}

// Log logs a message
func (j *Job) Log(severity string, message string) {
	logf(j, nil, severity, message)
}

// Logf logs a message with message interpolation/formatting
func (j *Job) Logf(severity string, format string, v ...any) {
	logf(j, nil, severity, format, v...)
}

// Default logs a message with DEFAULT severity
func (j *Job) Default(message string) {
	logf(j, nil, "DEFAULT", message)
}

// Defaultf logs a message with DEFAULT severity and message
// interpolation/formatting
func (j *Job) Defaultf(format string, v ...any) {
	logf(j, nil, "DEFAULT", format, v...)
}

// Debug logs a message with DEBUG severity
func (j *Job) Debug(message string) {
	logf(j, nil, "DEBUG", message)
}

// Debugf logs a message with DEBUG severity and message
// interpolation/formatting
func (j *Job) Debugf(format string, v ...any) {
	logf(j, nil, "DEBUG", format, v...)
}

// Info logs a message with INFO severity
func (j *Job) Info(message string) {
	logf(j, nil, "INFO", message)
}

// Infof logs a message with INFO severity and message
// interpolation/formatting
func (j *Job) Infof(format string, v ...any) {
	logf(j, nil, "INFO", format, v...)
}

// Notice logs a message with NOTICE severity
func (j *Job) Notice(message string) {
	logf(j, nil, "NOTICE", message)
}

// Noticef logs a message with NOTICE severity and message
// interpolation/formatting
func (j *Job) Noticef(format string, v ...any) {
	logf(j, nil, "NOTICE", format, v...)
}

// Warning logs a message with WARNING severity
func (j *Job) Warning(message string) {
	logf(j, nil, "WARNING", message)
}

// Warningf logs a message with WARNING severity and message
// interpolation/formatting
func (j *Job) Warningf(format string, v ...any) {
	logf(j, nil, "WARNING", format, v...)
}

// Error logs a message with ERROR severity
func (j *Job) Error(err error) {
	logf(j, nil, "ERROR", err.Error())
}

// Critical logs a message with CRITICAL severity
func (j *Job) Critical(message string) {
	logf(j, nil, "CRITICAL", message)
}

// Criticalf logs a message with CRITICAL severity and message
// interpolation/formatting
func (j *Job) Criticalf(format string, v ...any) {
	logf(j, nil, "CRITICAL", format, v...)
}

// Alert logs a message with ALERT severity
func (j *Job) Alert(message string) {
	logf(j, nil, "ALERT", message)
}

// Alertf logs a message with ALERT severity and message
// interpolation/formatting
func (j *Job) Alertf(format string, v ...any) {
	logf(j, nil, "ALERT", format, v...)
}

// Emergency logs a message with EMERGENCY severity
func (j *Job) Emergency(message string) {
	logf(j, nil, "EMERGENCY", message)
}

// Emergencyf logs a message with EMERGENCY severity and message
// interpolation/formatting
func (j *Job) Emergencyf(format string, v ...any) {
	logf(j, nil, "EMERGENCY", format, v...)
}

// Fatal logs a message and terminates the process.
func (j *Job) Fatal(err error) {
	log.Fatalf("fatal error: %v", err)
}
