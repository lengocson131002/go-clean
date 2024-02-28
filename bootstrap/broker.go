package bootstrap

import (
	"context"
	"strings"

	"github.com/lengocson131002/go-clean/pkg/env"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
	"github.com/lengocson131002/go-clean/pkg/transport/broker/kafka"
)

func GetKafkaBroker(cfg env.Configure, logger logger.Logger) broker.Broker {
	var addrs = cfg.GetString("KAFKA_BROKERS")
	var config = &kafka.KafkaBrokerConfig{
		Addresses: strings.Split(addrs, ","),
	}

	br, err := kafka.GetKafkaBroker(config, broker.WithLogger(logger))

	if err != nil {
		logger.Error(context.TODO(), "Failted to create kafka broker")
		panic(err)
	}

	return br
}
