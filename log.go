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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type LogEntry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	// Trace is the trace ID of the log message which will be propagated into
	// Cloud Trace.
	Trace string `json:"logging.googleapis.com/trace,omitempty"`
	// Component is the name of the service or job that produces the log entry.
	Component string `json:"component,omitempty"`
}

// String returns a JSON representation of the log entry.
func (le LogEntry) String() string {
	log.SetFlags(0)
	jsonBytes, err := json.Marshal(le)
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}

	return string(jsonBytes)
}

type logger interface {
	Name() string
	ProjectID() string
}

func logf(instance logger, r *http.Request, severity string, format string, v ...any) {
	if !isLogEntrySeverity(severity) {
		// Defaulting to the default
		severity = "DEFAULT"
	}

	message := fmt.Sprintf(format, v...)

	le := &LogEntry{
		Message:   message,
		Severity:  severity,
		Component: instance.Name(),
	}

	if r == nil {
		log.Println(le)
		return
	}

	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	ts := strings.Split(traceHeader, "/")
	if len(ts) > 0 && len(ts[0]) > 0 {
		le.Trace = fmt.Sprintf("projects/%s/traces/%s", instance.ProjectID(), ts[0])
	}

	log.Println(le)
}

func logEntrySeverities() []string {
	return []string{"DEFAULT", "INFO", "NOTICE", "ERROR", "DEBUG", "FATAL"}
}

func isLogEntrySeverity(severity string) bool {
	for _, elem := range logEntrySeverities() {
		if elem == severity {
			return true
		}
	}
	return false
}
