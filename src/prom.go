package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)

	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"path"})

	pageLoadTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "page_load_time_seconds",
		Help:    "Page load time in seconds.",
		Buckets: prometheus.DefBuckets,
	})

	resourceUtilizationCPU = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resource_utilization_cpu",
		Help: "Resource utilization of the application CPU.",
	})

	resourceUtilizationMem = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resource_utilization_mem",
		Help: "Resource utilization of the application memory.",
	})
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		start := time.Now()
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()
		httpDuration.WithLabelValues(path).Observe(duration)
		pageLoadTime.Observe(duration)

		// Retrieve CPU and memory utilization metrics
		cpuUsage, memoryUsage, err := getResourceUtilization()
		if err != nil {
			// Log the error, but don't stop the execution
			log.Printf("Error retrieving resource utilization: %v", err)
		} else {
			resourceUtilizationCPU.Set(cpuUsage)
			resourceUtilizationMem.Set(memoryUsage)
		}

	})
}

// getResourceUtilization retrieves CPU and memory utilization metrics
func getResourceUtilization() (float64, float64, error) {
	// Get CPU usage
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return 0, 0, err
	}

	// Get memory usage
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}

	// Return CPU and memory usage percentages
	return cpuUsage[0], memInfo.UsedPercent, nil
}

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(responseStatus)
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(pageLoadTime)
	prometheus.MustRegister(resourceUtilizationCPU)
	prometheus.MustRegister(resourceUtilizationMem)
	// prometheus.MustRegister(throughput)
}
