package bootstrap

import (
	"context"

	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/metrics/prome"
)

func GetPrometheusMetricer(logger logger.Logger) *prome.PrometheusMetricer {
	metricer, err := prome.NewPrometheusMetricer()

	if err != nil {
		logger.Error(context.TODO(), "Failed to create prometheous metricer")
		panic(err)
	}

	return metricer
}
