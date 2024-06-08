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
	"encoding/json"
	"log"
)

// LogEntry represents a log entry.
type LogEntry struct {
	// Message is the log message.
	Message string `json:"message"`
	// Severity is the severity of the log message.
	Severity string `json:"severity,omitempty"`
	// Trace is the trace ID of the log message.
	Trace string `json:"logging.googleapis.com/trace,omitempty"`
	// Component is the component that generated the log message.
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

// Fatal logs a fatal error and exits the program.
func Fatal(err error) {
	log.Fatalf("fatal error: %v", err)
}

// logEntrySeverities returns a list of valid log entry severities.
func logEntrySeverities() []string {
	return []string{"INFO", "NOTICE", "ERROR", "DEBUG", "FATAL"}
}

// isLogEntrySeverity checks if the given severity is a valid log entry severity.
func isLogEntrySeverity(severity string) bool {
	for _, elem := range logEntrySeverities() {
		if elem == severity {
			return true
		}
	}
	return false
}
