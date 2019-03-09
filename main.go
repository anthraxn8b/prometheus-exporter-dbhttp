package main

import (
     "flag"
      "net/http"

  log "github.com/Sirupsen/logrus"
      "github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":9876", "The address to listen on for HTTP requests.")

func main() {
  flag.Parse()
  http.Handle("/metrics", promhttp.Handler())
  log.Info("Beginning to serve on port ...")
  log.Fatal(http.ListenAndServe(*addr, nil))
}
