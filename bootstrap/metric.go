package bootstrap

import (
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/metrics/prometheous"
)

func GetPrometheusMetricer(logger logger.Logger) *prometheous.PrometheousMetricer {
	metricer, err := prometheous.NewPrometheusMetricer()

	if err != nil {
		logger.Error("Failed to create prometheous metricer")
		panic(err)
	}

	return metricer
}
