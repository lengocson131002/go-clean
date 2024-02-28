package kafka

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
)

var (
	DefaultBrokerConfig  = sarama.NewConfig()
	DefaultClusterConfig = sarama.NewConfig()
)

type brokerConfigKey struct{}
type clusterConfigKey struct{}

func BrokerConfig(c *sarama.Config) broker.BrokerOption {
	return setBrokerOption(brokerConfigKey{}, c)
}

func ClusterConfig(c *sarama.Config) broker.BrokerOption {
	return setBrokerOption(clusterConfigKey{}, c)
}

type subscribeContextKey struct{}

// SubscribeContext set the context for broker.SubscribeOption
func SubscribeContext(ctx context.Context) broker.SubscribeOption {
	return setSubscribeOption(subscribeContextKey{}, ctx)
}

type subscribeConfigKey struct{}

func SubscribeConfig(c *sarama.Config) broker.SubscribeOption {
	return setSubscribeOption(subscribeConfigKey{}, c)
}

// consumerGroupHandler is the implementation of sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	logger  logger.Logger
	handler broker.Handler
	subopts broker.SubscribeOptions
	kopts   broker.BrokerOptions
	cg      sarama.ConsumerGroup
	sess    sarama.ConsumerGroupSession
}

func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := sess.Context()
	for msg := range claim.Messages() {
		var m = broker.Message{}

		for _, header := range msg.Headers {
			m.Headers[string(header.Key)] = string(header.Value)
		}

		s := strings.Replace(string(msg.Value), "\\", "", 0)

		fmt.Printf(s)
		//
		m.Body = []byte(msg.Value)

		p := &publication{m: &m, t: msg.Topic, km: msg, cg: h.cg, sess: sess}
		eh := h.kopts.ErrorHandler

		// if err := json.Unmarshal(msg.Value, &m); err != nil {
		// 	p.err = err
		// 	p.m.Body = msg.Value
		// 	if eh != nil {
		// 		eh(p)
		// 	} else if h.logger != nil {
		// 		h.logger.Errorf(ctx, "[kafka]: failed to unmarshal: %v", err)
		// 	}
		// 	continue
		// }

		err := h.handler(p)
		if err == nil && h.subopts.AutoAck {
			sess.MarkMessage(msg, "")
		} else if err != nil {
			p.err = err
			if eh != nil {
				eh(p)
			} else if h.logger != nil {
				h.logger.Errorf(ctx, "[kafka]: subscriber error: %v", err)
			}
		}
	}
	return nil
}
