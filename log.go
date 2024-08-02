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
	"runtime"
	"strings"
)

// LogEntry is the structured version of a single log entry intended to be
// stored in Google Cloud Logging in JSON-serialized form.
type LogEntry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	// Trace is the trace ID of the log message which will be propagated into
	// Cloud Trace.
	Trace string `json:"logging.googleapis.com/trace,omitempty"`
	// SourceLocation holds the location within the source where the log message
	// was generated.
	SourceLocation *SourceLocation `json:"logging.googleapis.com/sourceLocation,omitempty"`
	// Component is the name of the service or job that produces the log entry.
	Component string `json:"component,omitempty"`
}

// SourceLocation is the structured version of a location in the source code (at
// compile time) which emits a log entry (at runtime). It is intended to be
// embedded in a run.LogEntry in JSON-serialized form.
type SourceLocation struct {
	File     string `json:"file,omitempty"`
	Function string `json:"function,omitempty"`
	Line     string `json:"line,omitempty"`
}

// String returns a JSON representation of the log entry.
func (le LogEntry) String() string {
	if Name() == "local" {
		return fmt.Sprintf("%s %s", le.Severity, le.Message)
	}
	log.SetFlags(0)
	jsonBytes, err := json.Marshal(le)
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}

	return string(jsonBytes)
}

// Log logs a message
func Log(r *http.Request, severity string, message string) {
	logf(r, severity, message)
}

// Logf logs a message with message interpolation/formatting
func Logf(r *http.Request, severity string, format string, v ...any) {
	logf(r, severity, format, v...)
}

// Default logs a message with DEFAULT severity
func Default(r *http.Request, message string) {
	logf(r, "DEFAULT", message)
}

// Defaultf logs a message with DEFAULT severity and message
// interpolation/formatting
func Defaultf(r *http.Request, format string, v ...any) {
	logf(r, "DEFAULT", format, v...)
}

// Debug logs a message with DEBUG severity
func Debug(r *http.Request, message string) {
	logf(r, "DEBUG", message)
}

// Debugf logs a message with DEBUG severity and message
// interpolation/formatting
func Debugf(r *http.Request, format string, v ...any) {
	logf(r, "DEBUG", format, v...)
}

// Info logs a message with INFO severity
func Info(r *http.Request, message string) {
	logf(r, "INFO", message)
}

// Infof logs a message with INFO severity and message
// interpolation/formatting
func Infof(r *http.Request, format string, v ...any) {
	logf(r, "INFO", format, v...)
}

// Notice logs a message with NOTICE severity
func Notice(r *http.Request, message string) {
	logf(r, "NOTICE", message)
}

// Noticef logs a message with NOTICE severity and message
// interpolation/formatting
func Noticef(r *http.Request, format string, v ...any) {
	logf(r, "NOTICE", format, v...)
}

// Warning logs a message with WARNING severity
func Warning(r *http.Request, message string) {
	logf(r, "WARNING", message)
}

// Warningf logs a message with WARNING severity and message
// interpolation/formatting
func Warningf(r *http.Request, format string, v ...any) {
	logf(r, "WARNING", format, v...)
}

// Error logs a message with ERROR severity
func Error(r *http.Request, err error) {
	logf(r, "ERROR", err.Error())
}

// Critical logs a message with CRITICAL severity
func Critical(r *http.Request, message string) {
	logf(r, "CRITICAL", message)
}

// Criticalf logs a message with CRITICAL severity and message
// interpolation/formatting
func Criticalf(r *http.Request, format string, v ...any) {
	logf(r, "CRITICAL", format, v...)
}

// Alert logs a message with ALERT severity
func Alert(r *http.Request, message string) {
	logf(r, "ALERT", message)
}

// Alertf logs a message with ALERT severity and message
// interpolation/formatting
func Alertf(r *http.Request, format string, v ...any) {
	logf(r, "ALERT", format, v...)
}

// Emergency logs a message with EMERGENCY severity
func Emergency(r *http.Request, message string) {
	logf(r, "EMERGENCY", message)
}

// Emergencyf logs a message with EMERGENCY severity and message
// interpolation/formatting
func Emergencyf(r *http.Request, format string, v ...any) {
	logf(r, "EMERGENCY", format, v...)
}

// Fatal logs a message and terminates the process.
func Fatal(r *http.Request, err error) {
	log.Fatalf("fatal error: %v", err)
}

func logf(r *http.Request, severity string, format string, v ...any) {
	log.SetFlags(0)
	if !isLogEntrySeverity(severity) {
		// Defaulting to the default, duh
		severity = "DEFAULT"
	}

	location := &SourceLocation{}
	caller, file, line, ok := runtime.Caller(2)
	if ok {
		location.File = file
		location.Line = fmt.Sprintf("%d", line)
		location.Function = runtime.FuncForPC(caller).Name()
	}

	message := fmt.Sprintf(format, v...)
	component := Name()

	le := &LogEntry{
		Severity:       severity,
		SourceLocation: location,
		Message:        message,
		Component:      component,
	}

	if r == nil {
		log.Println(le)
		return
	}

	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	ts := strings.Split(traceHeader, "/")
	if len(ts) > 0 && len(ts[0]) > 0 {
		le.Trace = fmt.Sprintf("projects/%s/traces/%s", ProjectID(), ts[0])
	}

	log.Println(le)
}

func logEntrySeverities() []string {
	// reference: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity
	return []string{
		"DEFAULT",
		"DEBUG",
		"INFO",
		"NOTICE",
		"WARNING",
		"ERROR",
		"CRITICAL",
		"ALERT",
		"EMERGENCY",
		"FATAL",
	}
}

func isLogEntrySeverity(severity string) bool {
	for _, elem := range logEntrySeverities() {
		if elem == severity {
			return true
		}
	}
	return false
}
