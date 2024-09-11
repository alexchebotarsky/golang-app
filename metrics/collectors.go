package metrics

import (
	"math"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsHandled = newCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_handled",
		Help: "Handled requests counter and metadata associated with them",
	},
		[]string{"status_code", "route_name"},
	))
	requestsDuration = newCollector(prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "requests_duration",
		Help:    "Time spent processing requests",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0, math.Inf(1)},
	}))
	eventsProcessed = newCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "events_processed",
		Help: "Handled PubSub events counter and metadata associated with them",
	},
		[]string{"event_name", "status"},
	))
	eventsDuration = newCollector(prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "events_duration",
		Help:    "Time spent processing events",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0, math.Inf(1)},
	},
		[]string{"event_name"},
	))
)

func RecordRequestHandled(statusCode int, routeName string) {
	requestsHandled.WithLabelValues(strconv.Itoa(statusCode), routeName).Inc()
}

func ObserveRequestDuration(duration time.Duration) {
	requestsDuration.Observe(duration.Seconds())
}

func RecordEventProcessed(eventName, status string) {
	eventsProcessed.WithLabelValues(eventName, status).Inc()
}

func ObserveEventDuration(eventName string, duration time.Duration) {
	eventsDuration.WithLabelValues(eventName).Observe(duration.Seconds())
}
