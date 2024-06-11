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
	"os"
	"strconv"
)

// NewJob creates a new Job instance.
//
// The Job instance will be populated with information from the environment variables
// set by Cloud Run and data available on the GCE metadata server.
func NewJob() *Job {
	log.SetFlags(0)
	j := &Job{
		Configs: make(map[string]string),
		Clients: make(map[string]interface{}),
	}

	j.Name = os.Getenv("CLOUD_RUN_JOB")
	if j.Name == "" {
		j.Name = "local"
	}

	j.Execution = os.Getenv("CLOUD_RUN_EXECUTION")
	if j.Execution == "" {
		j.Execution = "local"
	}

	task, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_INDEX"))
	if err == nil {
		j.TaskIndex = task
	}

	attempt, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_ATTEMPT"))
	if err == nil {
		j.TaskAttempt = attempt
	}

	count, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TASK_COUNT"))
	if err == nil {
		j.TaskCount = count
	}

	project, err := ProjectID()
	if err != nil {
		project = "local"
	}
	j.Project = project

	projectNumber, err := ProjectNumber()
	if err != nil {
		projectNumber = "local"
	}
	j.ProjectNumber = projectNumber

	return j
}

type Job struct {
	Configs       map[string]string
	Clients       map[string]interface{}
	Name          string
	Execution     string
	Project       string
	ProjectNumber string
	TaskIndex     int
	TaskAttempt   int
	TaskCount     int
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
		Component: j.Name,
	})
}

func (j *Job) Logf(severity string, format string, v ...any) {
	j.Log(severity, fmt.Sprintf(format, v...))
}
