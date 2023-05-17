package main

import (
	"backendify/cache"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	status := flag.Int("status", 204, "The option to turn on the status endpoint")
	port := flag.Int("port", 9000, "Port the service is listening on")
	loglevel := flag.String("loglevel", "InfoLevel", "Valid options are: InfoLevel|DebugLevel")
	timeout := flag.Duration("timeout", 200*time.Millisecond, "")
	flag.Parse()
	if *loglevel == "DebugLevel" {
		log.SetLevel(log.DebugLevel)
		log.Debug("Running in Debug mode!")
	}
	log.Infof("The code returned from the status endpoint %d", *status)
	log.Infof("The backend request timeout is set to %v", *timeout)

	// Configure StatsD metric
	config := &statsd.ClientConfig{
		Address: "127.0.0.1:8125",
	}
	client, err := statsd.NewClientWithConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	countryEndpointsMap := make(map[string]string, 1)
	endpointsMap(countryEndpointsMap, os.Args)
	cache := cache.NewCache()
	apiHandler := newApiHandler(countryEndpointsMap, cache, client, *timeout)

	apiHandler.collectMetric(1)

	mux := http.NewServeMux()
	mux.Handle("/company", apiHandler)
	mux.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(*status)
	})

	log.Infof("Starting the adapter-proxy server on port %d ... \n", *port)
	listenAddress := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}

func endpointsMap(backendAddrs map[string]string, args []string) {
	for _, arg := range os.Args[1:] {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			continue
		}
		if len(parts[0]) != 2 {
			log.Infof("Skipping invalid backend address for arg '%s'", arg)
			continue
		}
		//validate address
		_, err := url.ParseRequestURI(parts[1])
		if err != nil {
			log.Infof("Skipping invalid backend address '%s' for country '%s'", parts[1], parts[0])
			continue
		}
		backendAddrs[parts[0]] = parts[1]
		log.Infof("Added address '%s' for country '%s'\n", parts[1], parts[0])
	}
}
