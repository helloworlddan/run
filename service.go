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

// Service is intended to be instantiated once and kept around to access
// functionality related to the Cloud Run Service runtime.
type Service struct {
	server   *http.Server
	router   *http.ServeMux
	shutdown func(ctx context.Context, s *Service)
	configs  map[string]string
	clients  map[string]interface{}
}

// NewService creates a new Service instance.
func NewService() *Service {
	log.SetFlags(0)
	s := &Service{
		router:   &http.ServeMux{},
		server:   &http.Server{},
		shutdown: func(ctx context.Context, s *Service) {},
		configs:  make(map[string]string),
		clients:  make(map[string]interface{}),
	}
	s.server.Handler = s.router

	// simple uptime check handler
	s.router.HandleFunc("/uptimez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	return s
}

// ID returns the ID of the serving instance
func (s *Service) ID() string {
	id, err := instanceID()
	if err != nil {
		id = "00000"
	}
	return id
}

// Name returns the name of the service
func (s *Service) Name() string {
	name, err := kNativeService()
	if err != nil {
		name = "local"
	}
	return name
}

// Revision returns the name of the current revision of the service
func (s *Service) Revision() string {
	revision, err := kNativeRevision()
	if err != nil {
		revision = fmt.Sprintf("%s-00001-xxx", s.Name())
	}
	return revision
}

// Port returns the assigned port of the service
func (s *Service) Port() string {
	port, err := port()
	if err != nil {
		port = "8080"
	}
	return port
}

// ProjectID returns the name of the containing Google Cloud project or "local"
func (s *Service) ProjectID() string {
	project, err := projectID()
	if err != nil {
		project = "local"
	}
	return project
}

// ProjectNumber returns the 12-digit project number of the containing Google
// Cloud project or "000000000000"
func (s *Service) ProjectNumber() string {
	number, err := projectNumber()
	if err != nil {
		number = "000000000000"
	}
	return number
}

// Region returns the Google Cloud region in which the service is running or "local"
func (s *Service) Region() string {
	region, err := region()
	if err != nil {
		region = "local"
	}
	return region
}

// ServiceAccountEmail returns the email of the service account assigned to the
// service
func (s *Service) ServiceAccountEmail() string {
	email, err := serviceAccountEmail()
	if err != nil {
		email = "local"
	}
	return email
}

// ServiceAccountToken returns an authentication token for the assigned service
// account to authorize requests.
func (s *Service) ServiceAccountToken() string {
	token, err := serviceAccountToken()
	if err != nil {
		token = "local"
	}
	return token
}

// ListenAndServe starts the HTTP server, listens and serves requests
//
// It also traps SIGINT and SIGTERM. Both signals will cause a graceful
// shutdown of the HTTP server and executes the user supplied
// `run.ShutdownFunc`.
func (s *Service) ListenAndServe() error {
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	s.server.Addr = fmt.Sprintf(":%s", s.Port())

	go func(s *Service, errChan chan<- error) {
		s.Noticef(nil, "started and listening on port %s", s.Port())
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

	s.Info(nil, "shutdown complete")
	return nil
}

// ShutdownFunc registers a supplied function to be executed on server shutdown
//
// This is useful to run clean up routines, flush caches, drain and terminate
// connections, etc.
func (s *Service) ShutdownFunc(handler func(ctx context.Context, s *Service)) {
	s.shutdown = handler
}

// HandleFunc registers `http.HandleFunc` to respond to requests
func (s *Service) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}

// GetConfig retrieves a config value from the store
func (s *Service) GetConfig(key string) (string, error) {
	return getConfig(s.configs, key)
}

// PutConfig puts a config value in the store
func (s *Service) PutConfig(key string, val string) {
	putConfig(s.configs, key, val)
}

// LoadConfig looks up an environment variable puts it in the store and returns
// it's value
func (s *Service) LoadConfig(env string) (string, error) {
	return loadConfig(s.configs, env)
}

// ListConfigKeys returns a list of all available config keys
func (s *Service) ListConfigKeys() []string {
	return listConfigKeys(s.configs)
}

// GetClient resolves a client by name from the store
func (s *Service) GetClient(name string) (any, error) {
	return getClient(s.clients, name)
}

// AddClient add a client to the store
func (s *Service) AddClient(name string, client any) {
	addClient(s.clients, name, client)
}

// ListClientNames returns a list of all available clients
func (s *Service) ListClientNames() []string {
	return listClientNames(s.clients)
}

// Notice logs a message with NOTICE severity
func (s *Service) Notice(r *http.Request, message string) {
	s.Log(r, "NOTICE", message)
}

// Noticef logs a message with NOTICE severity and message
// interpolation/formatting
func (s *Service) Noticef(r *http.Request, format string, v ...any) {
	s.Logf(r, "NOTICE", format, v...)
}

// Info logs a message with INFO severity
func (s *Service) Info(r *http.Request, message string) {
	s.Log(r, "INFO", message)
}

// Infof logs a message with INFO severity and message
// interpolation/formatting
func (s *Service) Infof(r *http.Request, format string, v ...any) {
	s.Logf(r, "INFO", format, v...)
}

// Debug logs a message with DEBUG severity
func (s *Service) Debug(r *http.Request, message string) {
	s.Log(r, "DEBUG", message)
}

// Debugf logs a message with DEBUG severity and message
// interpolation/formatting
func (s *Service) Debugf(r *http.Request, format string, v ...any) {
	s.Logf(r, "DEBUG", format, v...)
}

// Error logs a message with ERROR severity
func (s *Service) Error(r *http.Request, err error) {
	s.Log(r, "ERROR", err.Error())
}

// Fatal logs a message and terminates the process.
func (s *Service) Fatal(r *http.Request, err error) {
	Fatal(err)
}

// Error logs a message
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

// Noticef logs a message with message interpolation/formatting
func (s *Service) Logf(r *http.Request, severity string, format string, v ...any) {
	s.Log(r, severity, fmt.Sprintf(format, v...))
}
