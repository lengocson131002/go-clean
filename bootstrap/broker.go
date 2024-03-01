package bootstrap

import (
	"context"
	"strings"

	"github.com/lengocson131002/go-clean/pkg/config"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
	"github.com/lengocson131002/go-clean/pkg/transport/broker/kafka"
)

func GetKafkaBroker(cfg config.Configure, logger logger.Logger) broker.Broker {
	var (
		addrs             = strings.Split(cfg.GetString("KAFKA_BROKERS"), ",")
		TLSEnabled        = cfg.GetBool("KAFKA_TLS_ENABLED")
		TLSSkipVerify     = cfg.GetBool("KAFKA_TLS_SKIP_VERIFY")
		TLSCaCertFile     = cfg.GetString("KAFKA_CA_CERT_FILE")
		TLSClientKeyFile  = cfg.GetString("KAFKA_CLIENT_KEY_FILE")
		TLSClientCertFile = cfg.GetString("KAFKA_CLIENT_CERT_FILE")
		SASLEnabled       = cfg.GetBool("KAFKA_SASL_ENABLED")
		SASLAlgorithm     = cfg.GetString("KAFKA_SASL_ALGORITHM")
		SASLUser          = cfg.GetString("KAFKA_SASL_USER")
		SASLPassword      = cfg.GetString("KAFKA_SASL_PASSWORD")
	)
	var config = &kafka.KafkaBrokerConfig{
		Addresses:         addrs,
		TLSEnabled:        TLSEnabled,
		TLSSkipVerify:     TLSSkipVerify,
		TLSCaCertFile:     TLSCaCertFile,
		TLSClientCertFile: TLSClientCertFile,
		TLSClientKeyFile:  TLSClientKeyFile,
		SASLEnabled:       SASLEnabled,
		SASLAlgorithm:     SASLAlgorithm,
		SASLUser:          SASLUser,
		SASLPassword:      SASLPassword,
	}

	br, err := kafka.GetKafkaBroker(
		config,
		broker.WithLogger(logger),
	)

	if err != nil {
		logger.Error(context.TODO(), "Failted to create kafka broker")
		panic(err)
	}

	return br
}
