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
	eventProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_processed",
		Help: "Number of PubSub events processed",
	},
		[]string{"event_id"},
	)
	eventDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "event_duration",
		Help:    "Time spent processing events",
		Buckets: []float64{.001, .002, .003, .004, .005, .01, .02, .03, .04, .05, .1, .2, .3, .4, .5},
	})
)

func Init() error {
	if err := prometheus.Register(requestStatusCode); err != nil {
		return fmt.Errorf("error registering requestStatusCodes metrics collector: %v", err)
	}

	if err := prometheus.Register(requestDuration); err != nil {
		return fmt.Errorf("error registering timeToProcessRequest metrics collector: %v", err)
	}

	if err := prometheus.Register(eventProcessed); err != nil {
		return fmt.Errorf("error registering eventProcessed metrics collector: %v", err)
	}

	if err := prometheus.Register(eventDuration); err != nil {
		return fmt.Errorf("error registering eventDuration metrics collector: %v", err)
	}

	return nil
}

func RecordRequestStatusCode(statusCode int, routeID string) {
	requestStatusCode.WithLabelValues(strconv.Itoa(statusCode), routeID).Inc()
}

func ObserveRequestDuration(seconds float64) {
	requestDuration.Observe(seconds)
}

func RecordEventProcessed(eventID string) {
	eventProcessed.WithLabelValues(eventID).Inc()
}

func ObserveEventDuration(seconds float64) {
	eventDuration.Observe(seconds)
}
