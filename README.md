# Run

Package Run provides useful, yet opinionated integrations for workloads running
on Cloud Run. It loosely follows the
[general development tips](https://cloud.google.com/run/docs/tips/general) from
the official documentation.

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

 "cloud.google.com/go/bigquery"
 "github.com/helloworlddan/run"
)

func main() {
 http.HandleFunc("/", indexHandler)

 // Store config
 run.PutConfig("some-key", "some-val")

 // Store client with lazy initialization
 var bqClient *bigquery.Client
 run.LazyClient("bigquery", func() {
  run.Debug(nil, "lazy init: bigquery")
  var err error
  ctx := context.Background()
  bqClient, err = bigquery.NewClient(ctx, run.ProjectID())
  if err != nil {
   run.Error(nil, err)
  }
  run.Client("bigquery", bqClient)
 })

 // Define shutdown behavior and serve HTTP
 err := run.ServeHTTP(func(ctx context.Context) {
  run.Debug(nil, "shutting down connections...")
  if bqClient != nil { // Maybe nil due to lazy loading
   bqClient.Close()
  }
  run.Debug(nil, "connections closed")
 }, nil)
 if err != nil {
  run.Fatal(nil, err)
 }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
 fmt.Fprintf(w, "Name: %s\n", run.ServiceName())
 fmt.Fprintf(w, "Revision: %s\n", run.ServiceRevision())
 fmt.Fprintf(w, "ProjectID: %s\n", run.ProjectID())

 // Access config
 cfg, err := run.GetConfig("some-key")
 if err != nil {
  run.Error(r, err)
 }

 // Access client
 var client *bigquery.Client
 client, err = run.UseClient("bigquery", client)
 if err != nil {
  run.Error(nil, err)
 }
 // NOTE: use client
 _ = client

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

 err := run.ServeGRPC(func(ctx context.Context) {
  run.Debug(nil, "shutting down connections...")
  time.Sleep(time.Second * 1) // Pretending to clean up
  run.Debug(nil, "connections closed")
 }, server)
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
 run.Client("bigquery", bqClient)
 defer bqClient.Close()

 // Later usage
 var bqClient2 *bigquery.Client
 bqClient2, err = run.UseClient("bigquery", bqClient2)
 if err != nil {
  run.Error(nil, err)
 }

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
