# Run

Package Run provides useful, yet opinionated integrations for workloads running
on Cloud Run.

## Cloud Run Services

### HTTP Services

Read more in
[run-examples](https://github.com/helloworlddan/run-examples/tree/main/run-http-service).

```golang
package main

import (
 "context"
 "fmt"
 "net/http"
 "time"

 "github.com/helloworlddan/run"
)

func main() {
 http.HandleFunc("/", indexHandler)

 run.PutConfig("some-key", "some-val")

 shutdown := func(ctx context.Context) {
  run.Debug(nil, "shutting down connections...")
  time.Sleep(time.Second * 1) // Pretending to clean up
  run.Debug(nil, "connections closed")
 }

 err := run.ServeHTTP(shutdown, nil)
 if err != nil {
  run.Fatal(nil, err)
 }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
 fmt.Fprintf(w, "Name: %s\n", run.ServiceName())
 fmt.Fprintf(w, "Revision: %s\n", run.ServiceRevision())
 fmt.Fprintf(w, "ProjectID: %s\n", run.ProjectID())
 cfg, err := run.GetConfig("some-key")
 if err != nil {
  run.Error(r, err)
 }
 fmt.Fprintf(w, "Config[some-key]: %s\n", cfg)
 run.Debugf(r, "request completed")
}
```

### GRPC Services

Read more in
[run-examples](https://github.com/helloworlddan/run-examples/tree/main/run-grpc-service).

```golang
package main

import (
 "context"
 "time"

 "github.com/helloworlddan/run"
 "github.com/helloworlddan/run-examples/run-grpc-service/runclock"
 "google.golang.org/grpc"
)

func main() {
 server := grpc.NewServer()
 runclock.RegisterRunClockServer(server, clockServer{})

 shutdown := func(ctx context.Context) {
  run.Debug(nil, "shutting down connections...")
  time.Sleep(time.Second * 1) // Pretending to clean up
  run.Debug(nil, "connections closed")
 }

 err := run.ServeGRPC(shutdown, server)
 if err != nil {
  run.Fatal(nil, err)
 }
}

type clockServer struct {
 runclock.UnimplementedRunClockServer
}

func (srv clockServer) GetTime(ctx context.Context, in *runclock.Empty) (*runclock.Time, error) {
 now := time.Now()
 run.Debug(nil, "received request")
 return &runclock.Time{
  Formatted: now.GoString(),
 }, nil
}
```

A client implementation can be found
[here](https://github.com/helloworlddan/run-examples/tree/main/run-grpc-client).

## Cloud Run Jobs

Read more in
[run-examples](https://github.com/helloworlddan/run-examples/tree/main/run-job).

```golang
package main

import (
 "context"
 "net/http"

 "cloud.google.com/go/bigquery"
 "github.com/helloworlddan/run"
)

func main() {
 // Store config
 run.PutConfig("my.app.key", "some value")
 cfgVal, err := run.GetConfig("my.app.key")
 if err != nil {
  run.Debugf(nil, "unable to read config: %v", err)
 }
 run.Infof(nil, "loaded config: %s", cfgVal)

 // Store client
 ctx := context.Background()
 bqClient, err := bigquery.NewClient(ctx, run.ProjectID())
 if err != nil {
  run.Error(nil, err)
 }
 run.AddClient("bigquery", bqClient)

 // Later usage
 clientRef, err := run.GetClient("bigquery")
 if err != nil {
  run.Error(nil, err)
 }
 bqClient2 := clientRef.(*bigquery.Client)
 _ = bqClient2

 // Make service account authenticated requests
 req, err := http.NewRequest(http.MethodGet, "https://google.com", nil)
 if err != nil {
  run.Error(nil, err)
 }
 req = run.AddOAuth2Header(req)
 resp, err := http.DefaultClient.Do(req)
 if err != nil {
  run.Error(nil, err)
 }
 defer resp.Body.Close()
 // read response
}
```

## TODO

- Deal with local develop on GCE, with metadata server available
