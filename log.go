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
	"log"
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
		Fatal(err)
	}

	return string(jsonBytes)
}

func Fatal(err error) {
	log.Fatalf("fatal error: %v", err)
}

func logEntrySeverities() []string {
	return []string{"INFO", "NOTICE", "ERROR", "DEBUG", "FATAL"}
}

func isLogEntrySeverity(severity string) bool {
	for _, elem := range logEntrySeverities() {
		if elem == severity {
			return true
		}
	}
	return false
}
