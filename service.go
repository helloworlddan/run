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
	"net/http"
	"os"
	"strings"
)

// NewService creates a new Service instance.
//
// The Service instance will be populated with information from the environment variables
// set by Cloud Run.
func NewService() *Service {
	log.SetFlags(0)
	s := &Service{
		Router:  &http.ServeMux{},
		Configs: make(map[string]string),
		Clients: make(map[string]interface{}),
	}

	s.Name = os.Getenv("K_SERVICE")
	if s.Name == "" {
		s.Name = "local"
	}

	s.Revision = os.Getenv("K_REVISION")
	if s.Revision == "" {
		s.Revision = "local"
	}

	s.Port = os.Getenv("PORT")
	if s.Port == "" {
		s.Port = "8080"
	}

	project, err := ProjectID()
	if err != nil {
		project = "local"
	}
	s.Project = project

	projectNumber, err := ProjectNumber()
	if err != nil {
		projectNumber = "local"
	}
	s.ProjectNumber = projectNumber

	return s
}

// Service represents a Cloud Run service.
type Service struct {
	// Name is the name of the service.
	Name string
	// Revision is the name of the revision.
	Revision string
	// Port is the port that the service listens on.
	Port string

	// Project is the Google Cloud Project ID.
	Project string
	// ProjectNumber is the Google Cloud Project Number.
	ProjectNumber string

	// Router is the HTTP router for the service.
	Router *http.ServeMux
	// Configs is a map of configuration values.
	Configs map[string]string
	// Clients is a map of client instances.
	Clients map[string]interface{}
}

// Notice logs a NOTICE message.
func (s *Service) Notice(r *http.Request, message string) {
	s.Log(r, "NOTICE", message)
}

// Noticef logs a NOTICE message with formatting.
func (s *Service) Noticef(r *http.Request, format string, v ...any) {
	s.Logf(r, "NOTICE", format, v...)
}

// Info logs an INFO message.
func (s *Service) Info(r *http.Request, message string) {
	s.Log(r, "INFO", message)
}

// Infof logs an INFO message with formatting.
func (s *Service) Infof(r *http.Request, format string, v ...any) {
	s.Logf(r, "INFO", format, v...)
}

// Debug logs a DEBUG message.
func (s *Service) Debug(r *http.Request, message string) {
	s.Log(r, "DEBUG", message)
}

// Debugf logs a DEBUG message with formatting.
func (s *Service) Debugf(r *http.Request, format string, v ...any) {
	s.Logf(r, "DEBUG", format, v...)
}

// Error logs an ERROR message.
func (s *Service) Error(r *http.Request, err error) {
	s.Log(r, "ERROR", err.Error())
}

// Fatal logs a FATAL message and exits the program.
func (s *Service) Fatal(r *http.Request, err error) {
	Fatal(err)
}

// Log logs a message with the specified severity.
func (s *Service) Log(r *http.Request, severity string, message string) {
	if !isLogEntrySeverity(severity) {
		Fatal(fmt.Errorf("unknown severitiy: %s", severity))
	}

	le := &LogEntry{
		Message:   message,
		Severity:  severity,
		Component: s.Name,
	}

	if r == nil {
		log.Println(le)
		return
	}

	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	ts := strings.Split(traceHeader, "/")
	if len(ts) > 0 && len(ts[0]) > 0 {
		le.Trace = fmt.Sprintf("projects/%s/traces/%s", s.Project, ts[0])
	}

	log.Println(le)
}

// Logf logs a message with the specified severity and formatting.
func (s *Service) Logf(r *http.Request, severity string, format string, v ...any) {
	s.Log(r, severity, fmt.Sprintf(format, v...))
}
