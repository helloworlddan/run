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
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// Service is intended to be instantiated once and kept around to access
// functionality related to the Cloud Run Service runtime.
type Service struct {
	httpServer *http.Server
	httpRouter *http.ServeMux
	grpcServer *grpc.Server
	shutdown   func(ctx context.Context)
	configs    map[string]string
	clients    map[string]interface{}
}

// NewService creates a new Service instance.
func NewService(opt ...grpc.ServerOption) *Service {
	log.SetFlags(0)
	s := &Service{
		httpRouter: &http.ServeMux{},
		httpServer: &http.Server{},
		grpcServer: grpc.NewServer(opt...),
		shutdown:   func(ctx context.Context) {},
		configs:    make(map[string]string),
		clients:    make(map[string]interface{}),
	}
	s.httpServer.Handler = s.httpRouter

	// simple uptime check handler
	s.httpRouter.HandleFunc("/uptimez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	return s
}

func (s *Service) GRPCServer() *grpc.Server {
	return s.grpcServer
}

// String returns the name of the service to satisfy fmt.Stringer
func (s *Service) String() string {
	return Name()
}

// ListenAndServe starts the GRPC server, listens and serves requests
//
// It also traps SIGINT and SIGTERM. Both signals will cause a graceful
// shutdown of the GRPC server and executes the user supplied
// `run.ShutdownFunc`.
func (s *Service) ListenAndServeGRPC() error {
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func(s *Service, errChan chan<- error) {
		Noticef(nil, "started and listening on port %s", Port())
		listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", Port()))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		if err := s.grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			errChan <- err
		}
	}(s, errChan)

	select {
	case err := <-errChan:
		return err
	case sig := <-sigChan:
		Noticef(nil, "shutdown initiated by signal: %v", sig)
	}

	// Cloud Run 10 sec time out
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the http server by waiting on existing requests
	s.grpcServer.Stop()

	// User-supplied shutdown
	s.shutdown(ctx)

	Info(nil, "shutdown complete")
	return nil
}

// ListenAndServe starts the HTTP server, listens and serves requests
//
// It also traps SIGINT and SIGTERM. Both signals will cause a graceful
// shutdown of the HTTP server and executes the user supplied
// `run.ShutdownFunc`.
func (s *Service) ListenAndServeHTTP() error {
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	s.httpServer.Addr = fmt.Sprintf(":%s", Port())

	go func(s *Service, errChan chan<- error) {
		Noticef(nil, "started and listening on port %s", Port())
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(s, errChan)

	select {
	case err := <-errChan:
		return err
	case sig := <-sigChan:
		Noticef(nil, "shutdown initiated by signal: %v", sig)
	}

	// Cloud Run 10 sec time out
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the http server by waiting on existing requests
	if err := s.httpServer.Shutdown(ctx); err != nil {
		Fatal(nil, err)
	}

	// User-supplied shutdown
	s.shutdown(ctx)

	Info(nil, "shutdown complete")
	return nil
}

// ShutdownFunc registers a supplied function to be executed on server shutdown
//
// This is useful to run clean up routines, flush caches, drain and terminate
// connections, etc.
func (s *Service) ShutdownFunc(handler func(ctx context.Context)) {
	s.shutdown = handler
}

// HandleFunc registers `http.HandleFunc` to respond to requests
func (s *Service) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.httpRouter.HandleFunc(pattern, handler)
}

// HandleStatic registers a handle to servie static assets from `path`
func (s *Service) HandleStatic(pattern string, path string) {
	handler := http.FileServer(http.Dir(path))
	s.httpRouter.Handle(pattern, http.StripPrefix(pattern, handler))
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
