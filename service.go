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
	"time"
)

type Service struct {
	server   *http.Server
	router   *http.ServeMux
	shutdown func(ctx context.Context, s *Service)
	Configs  map[string]string
	Clients  map[string]interface{}
}

// NewService creates a new Service instance.
//
// The Service instance will be populated with information from the environment variables
// set by Cloud Run.
func NewService() *Service {
	log.SetFlags(0)
	s := &Service{
		router:   &http.ServeMux{},
		server:   &http.Server{},
		shutdown: func(ctx context.Context, s *Service) {},
		Configs:  make(map[string]string),
		Clients:  make(map[string]interface{}),
	}
	s.server.Handler = s.router

	// Simple uptime check handler
	s.router.HandleFunc("/uptimez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	return s
}

func (s *Service) Name() string {
	name, err := kNativeService()
	if err != nil {
		name = "local"
	}
	return name
}

func (s *Service) Revision() string {
	revision, err := kNativeRevision()
	if err != nil {
		revision = "001"
	}
	return revision
}

func (s *Service) Port() string {
	port, err := port()
	if err != nil {
		port = "8080"
	}
	return port
}

func (s *Service) ProjectID() string {
	project, err := projectID()
	if err != nil {
		project = "local"
	}
	return project
}

func (s *Service) ProjectNumber() string {
	number, err := projectNumber()
	if err != nil {
		number = "0000000000"
	}
	return number
}

func (s *Service) Region() string {
	region, err := region()
	if err != nil {
		region = "local"
	}
	return region
}

func (s *Service) ServiceAccountEmail() string {
	email, err := serviceAccountEmail()
	if err != nil {
		email = "local"
	}
	return email
}

func (s *Service) ServiceAccountToken() string {
	token, err := serviceAccountToken()
	if err != nil {
		token = "local"
	}
	return token
}

func (s *Service) ListenAndServe() error {
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	s.server.Addr = fmt.Sprintf(":%s", s.Port())

	go func(s *Service, errChan chan<- error) {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(s, errChan)

	select {
	case err := <-errChan:
		return err
	case sig := <-sigChan:
		s.Noticef(nil, "shutdown initiated by signal: %v", sig)
	}

	// Cloud Run 10 sec time out
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the http server by waiting on existing requests
	if err := s.server.Shutdown(ctx); err != nil {
		s.Fatal(nil, err)
	}

	// User-supplied shutdown
	s.shutdown(ctx, s)

	s.Notice(nil, "shutdown complete")
	return nil
}

func (s *Service) ShutdownFunc(handler func(ctx context.Context, s *Service)) {
	s.shutdown = handler
}

func (s *Service) HandleFunc(pattern string, handler func(w http.ResponseWriter, _ *http.Request)) {
	s.router.HandleFunc(pattern, handler)
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
		Component: s.Name(),
	}

	if r == nil {
		log.Println(le)
		return
	}

	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	ts := strings.Split(traceHeader, "/")
	if len(ts) > 0 && len(ts[0]) > 0 {
		le.Trace = fmt.Sprintf("projects/%s/traces/%s", s.ProjectID(), ts[0])
	}

	log.Println(le)
}

func (s *Service) Logf(r *http.Request, severity string, format string, v ...any) {
	s.Log(r, severity, fmt.Sprintf(format, v...))
}
