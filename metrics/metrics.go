package metrics

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var collectors []prometheus.Collector

func newCollector[T prometheus.Collector](collector T) T {
	collectors = append(collectors, collector)
	return collector
}

func Init() error {
	errs := []error{}

	for i, collector := range collectors {
		err := prometheus.Register(collector)
		if err != nil {
			return fmt.Errorf("error registering metrics collector with index %d: %v", i, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
