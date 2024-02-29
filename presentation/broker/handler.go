package broker

import "github.com/lengocson131002/go-clean/pkg/transport/broker"

type errorHandler struct {
}

func (h *errorHandler) Handle(e broker.Event) error {
	return nil
}
