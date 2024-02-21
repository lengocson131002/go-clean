package bootstrap

import (
	"github.com/lengocson131002/go-clean/pkg/metrics"
	metric "github.com/lengocson131002/go-clean/pkg/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	METRIC_LABEL_REQUEST_TYPE                 = "request_type"
	METRIC_LABEL_REQUEST_STATUS               = "status"
	METRIC_LABEL_VALUE_REQUEST_STATUS_SUCCESS = "success"
	METRIC_LABEL_VALUE_REQUEST_STATUS_ERROR   = "error"
)

type Metricer struct {
	requestCountMetrics metrics.Counter
}

func NewMetricer() (*Metricer, error) {
	requestCountMetrics := metric.NewCounterFrom(prometheus.CounterOpts{
		Name: "request_count_total",
		Help: "Total count of requests",
	}, []string{
		METRIC_LABEL_REQUEST_TYPE,
		METRIC_LABEL_REQUEST_STATUS,
	})

	return &Metricer{
		requestCountMetrics: requestCountMetrics,
	}, nil
}
