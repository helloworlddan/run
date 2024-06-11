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
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// NewService creates a new Service instance.
//
// The Service instance will be populated with information from the environment variables
// set by Cloud Run.
func NewService() *Service {
	log.SetFlags(0)
	s := &Service{
		Router:  &http.ServeMux{},
		server:  &http.Server{},
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

	s.server.Addr = fmt.Sprintf(":%s", s.Port)
	s.server.Handler = s.Router

	s.Signals = make(chan os.Signal, 1)
	signal.Notify(s.Signals, syscall.SIGTERM, syscall.SIGINT)

	return s
}

type Service struct {
	Router        *http.ServeMux
	server        *http.Server
	Signals       chan (os.Signal)
	Configs       map[string]string
	Clients       map[string]interface{}
	Name          string
	Revision      string
	Port          string
	Project       string
	ProjectNumber string
}

func (s *Service) ListenAndServe() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Service) Notice(r *http.Request, message string) {
	s.Log(r, "NOTICE", message)
}

func (s *Service) Noticef(r *http.Request, format string, v ...any) {
	s.Logf(r, "NOTICE", format, v...)
}

func (s *Service) Info(r *http.Request, message string) {
	s.Log(r, "INFO", message)
}

func (s *Service) Infof(r *http.Request, format string, v ...any) {
	s.Logf(r, "INFO", format, v...)
}

func (s *Service) Debug(r *http.Request, message string) {
	s.Log(r, "DEBUG", message)
}

func (s *Service) Debugf(r *http.Request, format string, v ...any) {
	s.Logf(r, "DEBUG", format, v...)
}

func (s *Service) Error(r *http.Request, err error) {
	s.Log(r, "ERROR", err.Error())
}

func (s *Service) Fatal(r *http.Request, err error) {
	Fatal(err)
}

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

func (s *Service) Logf(r *http.Request, severity string, format string, v ...any) {
	s.Log(r, severity, fmt.Sprintf(format, v...))
}
