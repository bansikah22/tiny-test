package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

type AppInfo struct {
	AppName             string
	Version             string
	PodName             string
	Timestamp           string
	Uptime              string
	TotalRequests       int64
	RequestsPerEndpoint map[string]int64
}

type Metrics struct {
	TotalRequests       int64
	RequestsPerEndpoint map[string]int64
	StartTime           time.Time
	mu                  sync.RWMutex
}

var metrics = &Metrics{
	RequestsPerEndpoint: make(map[string]int64),
	StartTime:           time.Now(),
}

func main() {
	// Get environment variables with defaults
	appVersion := getEnv("APP_VERSION", "1.0.0")
	podName := getEnv("POD_NAME", "unknown")

	// Create file server for static files
	fs := http.FS(staticFiles)
	fileServer := http.FileServer(fs)

	// Serve static files
	http.Handle("/static/", trackMetrics("/static/", http.StripPrefix("/static/", fileServer)))

	// Root handler - serve UI
	http.HandleFunc("/", trackMetricsFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		uptime := time.Since(metrics.StartTime)
		metrics.mu.RLock()
		endpointStats := make(map[string]int64)
		for k, v := range metrics.RequestsPerEndpoint {
			endpointStats[k] = v
		}
		metrics.mu.RUnlock()

		info := AppInfo{
			AppName:             "Tiny Test App",
			Version:             appVersion,
			PodName:             podName,
			Timestamp:           time.Now().Format(time.RFC3339),
			Uptime:              formatUptime(uptime),
			TotalRequests:       atomic.LoadInt64(&metrics.TotalRequests),
			RequestsPerEndpoint: endpointStats,
		}

		// Read and parse template
		tmplData, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Error reading template", http.StatusInternalServerError)
			log.Printf("Error reading template: %v", err)
			return
		}

		tmpl, err := template.New("index").Parse(string(tmplData))
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			log.Printf("Error parsing template: %v", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, info); err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}
	}))

	// Health check endpoint
	http.HandleFunc("/healthz", trackMetricsFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	// Version endpoint
	http.HandleFunc("/version", trackMetricsFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version": appVersion,
		})
	}))

	// Info endpoint - returns JSON with pod information
	http.HandleFunc("/info", trackMetricsFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(metrics.StartTime)
		metrics.mu.RLock()
		endpointStats := make(map[string]int64)
		for k, v := range metrics.RequestsPerEndpoint {
			endpointStats[k] = v
		}
		metrics.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"app_name":              "Tiny Test App",
			"version":               appVersion,
			"pod_name":              podName,
			"timestamp":             time.Now().Format(time.RFC3339),
			"uptime_seconds":        int64(uptime.Seconds()),
			"uptime_formatted":      formatUptime(uptime),
			"total_requests":        atomic.LoadInt64(&metrics.TotalRequests),
			"requests_per_endpoint": endpointStats,
		})
	}))

	// Metrics endpoint - Prometheus-style simple metrics
	http.HandleFunc("/metrics", trackMetricsFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(metrics.StartTime)
		metrics.mu.RLock()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "# HELP http_requests_total Total number of HTTP requests\n")
		fmt.Fprintf(w, "# TYPE http_requests_total counter\n")
		fmt.Fprintf(w, "http_requests_total %d\n", atomic.LoadInt64(&metrics.TotalRequests))
		fmt.Fprintf(w, "\n# HELP http_requests_per_endpoint_total Total number of HTTP requests per endpoint\n")
		fmt.Fprintf(w, "# TYPE http_requests_per_endpoint_total counter\n")
		for endpoint, count := range metrics.RequestsPerEndpoint {
			fmt.Fprintf(w, "http_requests_per_endpoint_total{endpoint=\"%s\"} %d\n", endpoint, count)
		}
		fmt.Fprintf(w, "\n# HELP uptime_seconds Application uptime in seconds\n")
		fmt.Fprintf(w, "# TYPE uptime_seconds gauge\n")
		fmt.Fprintf(w, "uptime_seconds %.2f\n", uptime.Seconds())
		metrics.mu.RUnlock()
	}))

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	log.Printf("Version: %s, Pod: %s", appVersion, podName)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func trackMetricsFunc(path string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&metrics.TotalRequests, 1)
		metrics.mu.Lock()
		metrics.RequestsPerEndpoint[path]++
		metrics.mu.Unlock()
		handler(w, r)
	}
}

func trackMetrics(path string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&metrics.TotalRequests, 1)
		metrics.mu.Lock()
		metrics.RequestsPerEndpoint[path]++
		metrics.mu.Unlock()
		handler.ServeHTTP(w, r)
	})
}

func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
