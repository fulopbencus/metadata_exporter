package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Create a custom registry to exclude default Go metrics
var customRegistry = prometheus.NewRegistry()

// Create a gauge metric to track the metadata count
var metadataCountGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "metadata_count",
		Help: "Number of metadata entries in the host's response",
	},
)

func init() {
	// Register the custom metric with the custom Prometheus registry
	customRegistry.MustRegister(metadataCountGauge)
}

// Function to fetch the metadata count from the given URL and count lines ending with "200"
func getMetadataCount(url string) (int, error) {
	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Scan the response line by line
	scanner := bufio.NewScanner(resp.Body)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line ends with "200"
		if strings.HasSuffix(line, "200") {
			count++
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading response body from %s: %v", url, err)
	}

	return count, nil
}

// Update the gauge with the latest metadata count
func updateMetadataCount(url string) {
	for {
		metadataCount, err := getMetadataCount(url)
		if err != nil {
			log.Printf("Error fetching metadata count from %s: %v", url, err)
			metadataCountGauge.Set(0) // Set gauge to 0 in case of error
			continue
		}
		// Update the Prometheus gauge with the metadata count
		metadataCountGauge.Set(float64(metadataCount))

		// Sleep for 20 seconds before fetching the count again
		time.Sleep(20 * time.Second)
	}
}

func main() {
	// Command-line flags for URL and port
	port := flag.String("port", "9091", "Port to expose metrics")
	url := flag.String("url", "empty", "URL to fetch metadata count from")
	flag.Parse()

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to regularly update the metadata count
	go updateMetadataCount(*url)

	// Set up HTTP handler for Prometheus metrics endpoint, using the custom registry
	http.Handle("/metrics", promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{}))

	// Start the HTTP server for Prometheus to scrape
	go func() {
		log.Printf("Target : %s ... Starting server on :%s", *url, *port)
		if err := http.ListenAndServe(":"+*port, nil); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	// Block until we receive a signal for graceful shutdown
	sig := <-sigs
	log.Printf("Received signal: %v, shutting down", sig)
}
