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

// Job represents a Cloud Run job.
type Job struct {
	// Name is the name of the job.
	Name string
	// Execution is the name of the execution.
	Execution string

	// TaskIndex is the index of the task within the job.
	TaskIndex int
	// TaskAttempt is the attempt number of the task.
	TaskAttempt int
	// TaskCount is the total number of tasks in the job.
	TaskCount int

	// Project is the Google Cloud Project ID.
	Project string
	// ProjectNumber is the Google Cloud Project Number.
	ProjectNumber string

	// Configs is a map of configuration values.
	Configs map[string]string
	// Clients is a map of client instances.
	Clients map[string]interface{}
}

// Notice logs a NOTICE message.
func (j *Job) Notice(message string) {
	j.Log("NOTICE", message)
}

// Noticef logs a NOTICE message with formatting.
func (j *Job) Noticef(message string, v ...any) {
	j.Logf("NOTICE", message, v...)
}

// Info logs an INFO message.
func (j *Job) Info(message string) {
	j.Log("INFO", message)
}

// Infof logs an INFO message with formatting.
func (j *Job) Infof(message string, v ...any) {
	j.Logf("INFO", message, v...)
}

// Debug logs a DEBUG message.
func (j *Job) Debug(message string) {
	j.Log("DEBUG", message)
}

// Debugf logs a DEBUG message with formatting.
func (j *Job) Debugf(message string, v ...any) {
	j.Logf("DEBUG", message, v...)
}

// Error logs an ERROR message.
func (j *Job) Error(err error) {
	j.Log("ERROR", err.Error())
}

// Fatal logs a FATAL message and exits the program.
func (j *Job) Fatal(err error) {
	Fatal(err)
}

// Log logs a message with the specified severity.
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

// Logf logs a message with the specified severity and formatting.
func (j *Job) Logf(severity string, format string, v ...any) {
	j.Log(severity, fmt.Sprintf(format, v...))
}
