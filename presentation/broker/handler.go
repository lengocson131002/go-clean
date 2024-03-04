package broker

import (
	"context"
	"encoding/json"

	"github.com/lengocson131002/go-clean/pkg/pipeline"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
)

func HandleBrokerEvent[TReq any, TRes any](b broker.Broker, e broker.Event, replyTopic string) error {
	res, err := handleEvent[TReq, TRes](e)

	var brokerRes interface{}
	if err != nil {
		brokerRes = broker.FailureResponse(err)
	} else {
		brokerRes = broker.SuccessResponse[TRes](res)
	}

	body, err := json.Marshal(brokerRes)
	if err != nil {
		return err
	}

	fMsg := broker.Message{
		Body: body,
	}

	if e.Message() != nil {
		fMsg.Headers = e.Message().Headers
	}

	return b.Publish(replyTopic, &fMsg)
}

func handleEvent[TReq any, TRes any](e broker.Event) (res TRes, err error) {
	if e.Message() == nil || len(e.Message().Body) == 0 {
		return *new(TRes), broker.EmptyMessageError{}
	}

	body := e.Message().Body
	var request TReq
	err = json.Unmarshal(body, &request)
	if err != nil {
		return *new(TRes), broker.InvalidDataFormatError{}
	}

	// business logic
	res, err = pipeline.Send[TReq, TRes](context.TODO(), request)
	if err != nil {
		return *new(TRes), err
	}

	return res, nil
}
