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
	"testing"
)

func Test_String(t *testing.T) {
	message := "this is a log message"
	severity := "ALERT"
	trace := "some-trace-key"
	component := "test"

	le := LogEntry{
		Message:   message,
		Severity:  severity,
		Trace:     trace,
		Component: component,
	}

	expect := `{"message":"this is a log message","severity":"ALERT","logging.googleapis.com/trace":"some-trace-key","component":"test"}`
	line := le.String()

	if line != expect {
		t.Fatalf("String() produced bad log line: '%s', want '%s'", line, expect)
	}
}

func Test_isLogEntrySeverity(t *testing.T) {
	if !isLogEntrySeverity("DEBUG") {
		t.Fatalf("isLogEntrySeverity() returns false for 'DEBUG', want true")
	}
	if !isLogEntrySeverity("DEFAULT") {
		t.Fatalf("isLogEntrySeverity() returns false for 'DEFAULT', want true")
	}
	if !isLogEntrySeverity("ALERT") {
		t.Fatalf("isLogEntrySeverity() returns false for 'ALERT', want true")
	}
	if isLogEntrySeverity("BOGUS") {
		t.Fatalf("isLogEntrySeverity() returns true for 'BOGUS', want false")
	}
	if isLogEntrySeverity("") {
		t.Fatalf("isLogEntrySeverity() returns true for '', want false")
	}
}
