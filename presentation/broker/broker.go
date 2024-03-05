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

	go func() {
		defer func(ctx context.Context) {
			if err := s.broker.Disconnect(); err != nil {
				s.logger.Errorf(ctx, "Failed to shutdown broker server: %v", err)
			}
			s.logger.Info(ctx, "Stop Broker Server")
		}(ctx)

		<-ctx.Done()
	}()

	csGroupOpt := broker.WithSubscribeGroup("go_clean")

	_, err = s.broker.Subscribe(RequestTopic, func(e broker.Event) error {
		return HandleBrokerEvent[*domain.OpenAccountRequest, *domain.OpenAccountResponse](s.broker, e, ReplyTopic)
	}, csGroupOpt)

	return err
}
