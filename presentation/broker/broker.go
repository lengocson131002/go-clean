package broker

import (
	"context"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
)

const (
	RequestTopic = "go.test.clean.request"
	ReplyTopic   = "go.test.clean.reply"
)

type BrokerServer struct {
	broker broker.Broker
	logger logger.Logger
}

func NewBrokerServer(broker broker.Broker, logger logger.Logger) *BrokerServer {
	return &BrokerServer{
		broker: broker,
		logger: logger,
	}
}

func (s *BrokerServer) Start(ctx context.Context) error {
	err := s.broker.Connect()
	if err == nil {
		s.logger.Infof(ctx, "Connected to broker server")
	}

	csGroupOpt := broker.WithSubscribeGroup("go_clean")

	s.broker.Subscribe(RequestTopic, func(e broker.Event) error {
		return HandleBrokerEvent[*domain.OpenAccountRequest, *domain.OpenAccountResponse](s.broker, e, ReplyTopic)
	}, csGroupOpt)

	return err
}
