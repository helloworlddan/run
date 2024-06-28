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
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	http2 "golang.org/x/net/http2"
	http2clear "golang.org/x/net/http2/h2c"
	grpc "google.golang.org/grpc"
)

// ServeGRPC starts the GRPC server, listens and serves requests
//
// It also traps SIGINT and SIGTERM. Both signals will cause a graceful
// shutdown of the GRPC server and executes the user supplied
// shutdown func.
func ServeGRPC(shutdown func(context.Context), server *grpc.Server) error {
	if server == nil {
		return errors.New("cannot listen using ni GRPC server")
	}

	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func(errChan chan<- error) {
		listener, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", ServicePort()))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			errChan <- err
		}
	}(errChan)
	Noticef(nil, "started and listening on port %s", ServicePort())

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
	server.Stop()

	// User-supplied shutdown
	shutdown(ctx)

	Info(nil, "shutdown complete")
	return nil
}

// ServeHTTP starts the HTTP server, listens and serves requests
//
// It also traps SIGINT and SIGTERM. Both signals will cause a graceful
// shutdown of the HTTP server and executes the user supplied
// shutdown func.
func ServeHTTP(shutdown func(context.Context), server *http.Server) error {
	if server == nil {
		mux := http.DefaultServeMux
		// Add default uptime check handler
		mux.HandleFunc("GET /uptimez", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		server = &http.Server{
			Addr: net.JoinHostPort("0.0.0.0", ServicePort()),
			// Support HTTP2
			Handler: http2clear.NewHandler(mux, &http2.Server{}),
		}
	}

	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func(errChan chan<- error) {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(errChan)
	Noticef(nil, "started and listening on port %s", ServicePort())

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
	if err := server.Shutdown(ctx); err != nil {
		Fatal(nil, err)
	}

	// User-supplied shutdown
	shutdown(ctx)

	Info(nil, "shutdown complete")
	return nil
}
