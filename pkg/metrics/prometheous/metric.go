package prometheous

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheousMetricer struct {
	RequestTotalCounter *prometheus.CounterVec
}

func NewPrometheusMetricer() (*PrometheousMetricer, error) {
	requestTotalCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%srequests_total", DEFAULT_METRIC_PREFIX),
		Help: "Total requests processed, partitioned by endpoint and status",
	}, []string{
		fmt.Sprintf("%s%s", DEFAULT_METRIC_LABEL_PREFIX, METRIC_LABEL_SERVICE),
		fmt.Sprintf("%s%s", DEFAULT_METRIC_LABEL_PREFIX, METRIC_LABEL_ENDPOINT),
		fmt.Sprintf("%s%s", DEFAULT_METRIC_LABEL_PREFIX, METRIC_LABEL_STATUS),
	})

	return &PrometheousMetricer{
		RequestTotalCounter: requestTotalCounter,
	}, nil
}
