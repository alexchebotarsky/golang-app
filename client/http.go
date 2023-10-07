package client

import (
	"net/http"
	"time"

	"github.com/goodleby/golang-server/tracing"
)

type Parameters struct {
	Timeout time.Duration
}

func NewHTTPClient(params Parameters) *http.Client {
	c := http.Client{
		Timeout:   params.Timeout,
		Transport: tracing.NewTracedTransport(http.DefaultTransport),
	}

	return &c
}
