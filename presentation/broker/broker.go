package broker

import (
	"context"
	"encoding/json"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/pipeline"
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

	consumerGroupOption := broker.WithSubscribeGroup("group_name")

	s.broker.Subscribe(RequestTopic, func(e broker.Event) error {
		// Step 1: Nhận request và xử lý
		if e.Message() == nil || len(e.Message().Body) == 0 {
			return broker.EmptyRequestError{}
		}

		body := e.Message().Body
		var t24RequestModel *domain.OpenAccountRequest
		err := json.Unmarshal(body, &t24RequestModel)
		if err != nil {
			return err
		}

		// business logic
		res, err := pipeline.Send[*domain.OpenAccountRequest, *domain.OpenAccountResponse](context.TODO(), t24RequestModel)
		if err != nil {
			return err
		}

		resByte, err := json.Marshal(res)
		if err != nil {
			return err
		}

		// Step 2:
		s.broker.Publish(ReplyTopic, &broker.Message{
			Body: resByte,
		})

		return nil
	}, consumerGroupOption)

	return err
}
