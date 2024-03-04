package kafka

import (
	"encoding/json"
	"math"
	"math/rand"
	"testing"

	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/logger/logrus"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
)

func getKafkaBroker() broker.Broker {
	var config = &KafkaBrokerConfig{
		Addresses: []string{"localhost:9092"},
	}

	br, err := GetKafkaBroker(
		config,
	)

	if err != nil {
		panic(err)
	}

	return br
}

func getLogger() logger.Logger {
	return logrus.NewLogrusLogger()
}

type KRequestType struct {
	Number int
}

type KResponseType struct {
	Result float64
}

func TestKafkaPublishAndReceive(b *testing.T) {
	var (
		kBroker      = getKafkaBroker()
		requestTopic = "go.clean.test.benchmark.request"
		replyTopic   = "go.clean.test.benchmark.reply"
	)

	err := kBroker.Connect()
	if err != nil {
		b.Fail()
	}

	_, err = kBroker.Subscribe(requestTopic, func(e broker.Event) error {
		msg := e.Message()
		if msg == nil {
			return broker.EmptyMessageError{}
		}

		var req KRequestType
		err := json.Unmarshal(msg.Body, &req)
		if err != nil {
			return broker.InvalidDataFormatError{}
		}

		b.Logf("Received request: %v", req)

		result := KResponseType{
			Result: math.Pow(float64(req.Number), 2),
		}

		resultByte, err := json.Marshal(result)
		if err != nil {
			return broker.InvalidDataFormatError{}
		}

		// pubish to response topic
		err = kBroker.Publish(replyTopic, &broker.Message{
			Headers: msg.Headers,
			Body:    resultByte,
		})

		if err != nil {
			b.Error(err)
		}

		return nil
	}, broker.WithSubscribeGroup("benchmark.test"))

	if err != nil {
		b.Error(err)
	}

	req := KRequestType{
		Number: rand.Intn(100),
	}
	reqByte, err := json.Marshal(req)
	if err != nil {
		b.Error(err)
	}
	msg, err := kBroker.PublishAndReceive(requestTopic, &broker.Message{
		Body: reqByte,
	}, broker.WithPublishReplyToTopic(replyTopic))

	if err != nil {
		b.Logf("error: %v", err)
	} else {
		b.Logf("Received message: %s", string(msg.Body))
	}

	msg2, err2 := kBroker.PublishAndReceive(requestTopic, &broker.Message{
		Body: reqByte,
	}, broker.WithPublishReplyToTopic(replyTopic))

	if err2 != nil {
		b.Logf("error: %v", err2)
	} else {
		b.Logf("Received message: %s", string(msg2.Body))
	}

}
