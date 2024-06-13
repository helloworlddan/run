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
	"fmt"
	"log"
)

type Job struct {
	configs map[string]string
	clients map[string]interface{}
}

// NewJob creates a new Job instance.
//
// The Job instance will be populated with information from the environment variables
// set by Cloud Run and data available on the GCE metadata server.
func NewJob() *Job {
	log.SetFlags(0)
	j := &Job{
		configs: make(map[string]string),
		clients: make(map[string]interface{}),
	}

	return j
}

func (j *Job) Name() string {
	name, err := jobName()
	if err != nil {
		name = "local"
	}
	return name
}

func (j *Job) Execution() string {
	execution, err := jobExecution()
	if err != nil {
		execution = "local"
	}
	return execution
}

func (j *Job) TaskIndex() int {
	index, err := jobTaskIndex()
	if err != nil {
		index = 0
	}
	return index
}

func (j *Job) TaskAttempt() int {
	attempt, err := jobTaskAttempt()
	if err != nil {
		attempt = 0
	}
	return attempt
}

func (j *Job) TaskCount() int {
	count, err := jobTaskCount()
	if err != nil {
		count = 0
	}
	return count
}

func (j *Job) ProjectID() string {
	project, err := projectID()
	if err != nil {
		project = "local"
	}
	return project
}

func (j *Job) ProjectNumber() string {
	number, err := projectNumber()
	if err != nil {
		number = "000000000000"
	}
	return number
}

func (j *Job) Region() string {
	region, err := region()
	if err != nil {
		region = "local"
	}
	return region
}

func (j *Job) ServiceAccountEmail() string {
	email, err := serviceAccountEmail()
	if err != nil {
		email = "local"
	}
	return email
}

func (j *Job) ServiceAccountToken() string {
	token, err := serviceAccountToken()
	if err != nil {
		token = "local"
	}
	return token
}

func (j *Job) GetConfig(key string) (string, error) {
	return getConfig(j.configs, key)
}

func (j *Job) PutConfig(key string, val string) {
	putConfig(j.configs, key, val)
}

func (j *Job) LoadConfig(env string) (string, error) {
	return loadConfig(j.configs, env)
}

func (j *Job) ListConfig() []string {
	return listConfig(j.configs)
}

func (j *Job) Notice(message string) {
	j.Log("NOTICE", message)
}

func (j *Job) Noticef(message string, v ...any) {
	j.Logf("NOTICE", message, v...)
}

func (j *Job) Info(message string) {
	j.Log("INFO", message)
}

func (j *Job) Infof(message string, v ...any) {
	j.Logf("INFO", message, v...)
}

func (j *Job) Debug(message string) {
	j.Log("DEBUG", message)
}

func (j *Job) Debugf(message string, v ...any) {
	j.Logf("DEBUG", message, v...)
}

func (j *Job) Error(err error) {
	j.Log("ERROR", err.Error())
}

func (j *Job) Fatal(err error) {
	Fatal(err)
}

func (j *Job) Log(severity string, message string) {
	if !isLogEntrySeverity(severity) {
		Fatal(fmt.Errorf("unknown severity: %s", severity))
	}

	log.Println(LogEntry{
		Message:   message,
		Severity:  severity,
		Component: j.Name(),
	})
}

func (j *Job) Logf(severity string, format string, v ...any) {
	j.Log(severity, fmt.Sprintf(format, v...))
}
