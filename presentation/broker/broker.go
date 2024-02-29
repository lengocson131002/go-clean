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

	// create t24 account

	// SERVICE A (Client):
	// Step 1: Gửi requset, tạo Signal ID

	// Step 2: Luôn lắng nghe response queue: khi có event -> trigger signal

	// SERVICE B (Server):
	// Step 1: Nhận request => Xử lý
	// Step 2: Gửi response vào response queue

	s.broker.Subscribe(RequestTopic, func(e broker.Event) error {
		if e.Message() == nil {
			// ignore
			return &broker.EmptyRequestError{}
		}

		body := e.Message().Body

		if len(body) == 0 {
			return nil
		}

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

		s.broker.Publish(ReplyTopic, &broker.Message{
			Body: resByte,
		})

		return nil
	})

	return err
}
