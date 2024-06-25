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
	"fmt"
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

// ListenAndServe starts the GRPC server, listens and serves requests
//
// It also traps SIGINT and SIGTERM. Both signals will cause a graceful
// shutdown of the GRPC server and executes the user supplied
// shutdown func.
func ListenAndServeGRPC(shutdown func(context.Context), server *grpc.Server) error {
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	if server == nil {
		return errors.New("cannot listen using ni GRPC server")
	}

	go func(errChan chan<- error) {
		Noticef(nil, "started and listening on port %s", ServicePort())
		listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", ServicePort()))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			errChan <- err
		}
	}(errChan)

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
	return grpc.ErrServerStopped
}

// ListenAndServe starts the HTTP server, listens and serves requests
//
// It also traps SIGINT and SIGTERM. Both signals will cause a graceful
// shutdown of the HTTP server and executes the user supplied
// shutdown func.
func ListenAndServeHTTP(shutdown func(context.Context), server *http.Server) error {
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	if server == nil {
		server = &http.Server{
			Addr:    net.JoinHostPort("0.0.0.0", ServicePort()),
			Handler: http2clear.NewHandler(http.DefaultServeMux, &http2.Server{}),
		}
	}

	go func(errChan chan<- error) {
		Noticef(nil, "started and listening on port %s", ServicePort())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(errChan)

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
	return http.ErrServerClosed
}
