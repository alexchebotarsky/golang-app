package metrics

import (
	"fmt"
	"math"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestStatusCode = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_status_code",
		Help: "Status codes returned by the API",
	},
		[]string{"status_code", "route_id"},
	)
	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration",
		Help:    "Time spent processing requests",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0, math.Inf(1)},
	})
)

func Init() error {
	if err := prometheus.Register(requestStatusCode); err != nil {
		return fmt.Errorf("error registering requestStatusCodes metrics collector: %v", err)
	}

	if err := prometheus.Register(requestDuration); err != nil {
		return fmt.Errorf("error registering timeToProcessRequest metrics collector: %v", err)
	}

	return nil
}

func RecordRequestStatusCode(statusCode int, routeID string) {
	requestStatusCode.WithLabelValues(strconv.Itoa(statusCode), routeID).Inc()
}

func ObserveRequestDuration(seconds float64) {
	requestDuration.Observe(seconds)
}
