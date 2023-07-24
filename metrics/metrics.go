package metrics

import (
	"fmt"
	"math"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestStatusCodes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_status_code",
		Help: "Status codes returned by the API",
	},
		[]string{"status_code", "operation_name"},
	)
	timeToProcessRequest = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration",
		Help:    "Time spent processing requests",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0, math.Inf(1)},
	})
)

// Init initializes prometheus collectors.
func Init() error {
	if err := prometheus.Register(requestStatusCodes); err != nil {
		return fmt.Errorf("error registering requestStatusCodes metrics collector: %v", err)
	}

	if err := prometheus.Register(timeToProcessRequest); err != nil {
		return fmt.Errorf("error registering timeToProcessRequest metrics collector: %v", err)
	}

	return nil
}

// ObserveTimeToProcess records the time spent processing an operation.
func ObserveTimeToProcess(operation string, t float64) {
	timeToProcessRequest.Observe(t)
}

// RecordRequestStatusCode records the status code returned for each request.
func RecordRequestStatusCode(statusCode int, operationName string) {
	requestStatusCodes.WithLabelValues(strconv.Itoa(statusCode), operationName).Inc()
}
