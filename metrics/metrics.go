package metrics

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsHandled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_handled",
		Help: "Handled requests counter and metadata associated with them",
	},
		[]string{"status_code", "route_name"},
	)
	requestsDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "requests_duration",
		Help:    "Time spent processing requests",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0, math.Inf(1)},
	})
	eventsHandled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "events_processed",
		Help: "Handled PubSub events counter and metadata associated with them",
	},
		[]string{"event_name", "status"},
	)
	eventsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "events_duration",
		Help:    "Time spent processing events",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0, math.Inf(1)},
	},
		[]string{"event_name"},
	)
)

func Init() error {
	err := prometheus.Register(requestsHandled)
	if err != nil {
		return fmt.Errorf("error registering requestsHandled metrics collector: %v", err)
	}

	err = prometheus.Register(requestsDuration)
	if err != nil {
		return fmt.Errorf("error registering requestsDuration metrics collector: %v", err)
	}

	err = prometheus.Register(eventsHandled)
	if err != nil {
		return fmt.Errorf("error registering eventsHandled metrics collector: %v", err)
	}

	err = prometheus.Register(eventsDuration)
	if err != nil {
		return fmt.Errorf("error registering eventsDuration metrics collector: %v", err)
	}

	return nil
}

func RecordRequestStatusCode(statusCode int, routeName string) {
	requestsHandled.WithLabelValues(strconv.Itoa(statusCode), routeName).Inc()
}

func ObserveRequestDuration(duration time.Duration) {
	requestsDuration.Observe(duration.Seconds())
}

func RecordEventProcessed(eventName, status string) {
	eventsHandled.WithLabelValues(eventName, status).Inc()
}

func ObserveEventDuration(eventName string, duration time.Duration) {
	eventsDuration.WithLabelValues(eventName).Observe(duration.Seconds())
}
