# Run

Package Run provides useful, yet opinionated integrations for workloads running
on Cloud Run.

## Cloud Run Services

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
 service := run.NewService()

 service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Name: %s\n", service.Name())
  fmt.Fprintf(w, "Revision: %s\n", service.Revision())
  fmt.Fprintf(w, "ProjectID: %s\n", service.ProjectID())
  service.Debugf(r, "request completed")
 })

 service.ShutdownFunc(func(ctx context.Context, s *run.Service) {
  s.Debug(nil, "shutting down connections...")
  time.Sleep(time.Second * 1) // Pretending to clean up
  s.Debug(nil, "connections closed")
 })

 err := service.ListenAndServe()
 if err != nil {
  service.Fatal(nil, err)
 }
}
```

## Cloud Run Jobs

```golang
package main

import (
 "context"
 "net/http"

 "cloud.google.com/go/bigquery"
 "github.com/helloworlddan/run"
)

func main() {
 job := run.NewJob()

 // Store config
 job.PutConfig("my.app.key", "some value")
 cfgVal, err := job.GetConfig("my.app.key")
 if err != nil {
  job.Debugf("unable to read config: %v", err)
 }
 job.Infof("loaded config: %s", cfgVal)

 // Store client
 bqClient, err := bigquery.NewClient()
 if err != nil {
  job.Error(err)
 }
 job.AddClient("bigquery", bqClient)

 // Later usage
 clientRef, err := job.GetClient("bigquery")
 if err != nil {
  job.Error(err)
 }
 bqClient2 := clientRef.(*bigquery.Client)
 _ = bqClient2

 // Make service account authenticated requests
 ctx := context.Background()
 req, err := job.NewAuthenticatedRequest(ctx, http.MethodGet, "https://example.com", nil)
 if err != nil {
  job.Error(err)
 }
 resp, err := http.DefaultClient.Do(req)
 if err != nil {
  job.Error(err)
 }
 defer resp.Body.Close()
 // read response
}

```

## TODO

- Deal with local develop on GCE, with metadata server available
- Accessing knative spec and exposing config
