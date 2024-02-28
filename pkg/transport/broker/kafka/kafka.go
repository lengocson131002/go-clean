package kafka

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/logger/logrus"
	"github.com/lengocson131002/go-clean/pkg/transport/broker"
)

type kBroker struct {
	addrs []string

	c sarama.Client       // broker connection client
	p sarama.SyncProducer // producer

	sc []sarama.Client // subscription connectiion clients

	connected bool
	scMutex   sync.Mutex
	opts      broker.BrokerOptions
}

func NewKafkaBroker(opts ...broker.BrokerOption) broker.Broker {
	options := broker.BrokerOptions{
		Context: context.Background(),
		Logger:  logrus.NewLogrusLogger(), // using logrus logging by default
	}

	for _, o := range opts {
		o(&options)
	}

	var cAddrs []string
	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}

	if len(cAddrs) == 0 {
		cAddrs = []string{DefaultKafkaBroker}
	}

	return &kBroker{
		addrs: cAddrs,
		opts:  options,
	}
}

type subscriber struct {
	cg   sarama.ConsumerGroup
	t    string
	opts broker.SubscribeOptions
}

type publication struct {
	t    string
	err  error
	cg   sarama.ConsumerGroup
	km   *sarama.ConsumerMessage
	m    *broker.Message
	sess sarama.ConsumerGroupSession
}

func (p *publication) Topic() string {
	return p.t
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) Ack() error {
	p.sess.MarkMessage(p.km, "")
	return nil
}

func (p *publication) Error() error {
	return p.err
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	return s.cg.Close()
}

func (k *kBroker) Address() string {
	if len(k.addrs) > 0 {
		return k.addrs[0]
	}
	return DefaultKafkaBroker
}

func (k *kBroker) Connect() error {
	if k.connected {
		return nil
	}

	k.scMutex.Lock()
	if k.c != nil {
		k.scMutex.Unlock()
		return nil
	}
	k.scMutex.Unlock()

	pconfig := k.getBrokerConfig()
	// For implementation reasons, the SyncProducer requires
	// `Producer.Return.Errors` and `Producer.Return.Successes`
	// to be set to true in its configuration.
	pconfig.Producer.Return.Successes = true
	pconfig.Producer.Return.Errors = true

	c, err := sarama.NewClient(k.addrs, pconfig)
	if err != nil {
		return err
	}

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		return err
	}

	k.scMutex.Lock()
	k.c = c
	k.p = p
	k.sc = make([]sarama.Client, 0)
	k.connected = true
	defer k.scMutex.Unlock()

	return nil
}

func (k *kBroker) Disconnect() error {
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	for _, client := range k.sc {
		client.Close()
	}
	k.sc = nil
	k.p.Close()
	if err := k.c.Close(); err != nil {
		return err
	}
	k.connected = false
	return nil
}

func (k *kBroker) Init(opts ...broker.BrokerOption) error {
	for _, o := range opts {
		o(&k.opts)
	}
	var cAddrs []string
	for _, addr := range k.opts.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{DefaultKafkaBroker}
	}
	k.addrs = cAddrs
	return nil
}

func (k *kBroker) Options() broker.BrokerOptions {
	return k.opts
}

func (k *kBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, _, err = k.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	})

	return err
}

func (k *kBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
		Queue:   uuid.New().String(),
	}

	for _, o := range opts {
		o(&opt)
	}
	// we need to create a new client per consumer
	c, err := k.getSaramaClusterClient(topic)
	if err != nil {
		return nil, err
	}
	cg, err := sarama.NewConsumerGroupFromClient(opt.Queue, c)
	if err != nil {
		return nil, err
	}
	h := &consumerGroupHandler{
		handler: handler,
		subopts: opt,
		kopts:   k.opts,
		cg:      cg,
		logger:  k.getLogger(),
	}

	ctx := context.Background()
	topics := []string{topic}
	go func() {
		for {
			select {
			case err := <-cg.Errors():
				if err != nil {
					if logger := k.getLogger(); logger != nil {
						logger.Errorf(ctx, "consumer error:", err)
					}
				}
			default:
				err := cg.Consume(ctx, topics, h)
				switch err {
				case sarama.ErrClosedConsumerGroup:
					return
				case nil:
					continue
				default:
					if logger := k.getLogger(); logger != nil {
						logger.Errorf(ctx, "consumer error:", err)
					}
				}
			}
		}
	}()
	return &subscriber{cg: cg, opts: opt, t: topic}, nil
}

func (k *kBroker) getBrokerConfig() *sarama.Config {
	if c, ok := k.opts.Context.Value(brokerConfigKey{}).(*sarama.Config); ok {
		return c
	}
	return DefaultBrokerConfig
}

func (k *kBroker) getClusterConfig() *sarama.Config {
	if c, ok := k.opts.Context.Value(clusterConfigKey{}).(*sarama.Config); ok {
		return c
	}
	clusterConfig := DefaultClusterConfig

	// the oldest supported version is V0_10_2_0
	if !clusterConfig.Version.IsAtLeast(sarama.V0_10_2_0) {
		clusterConfig.Version = sarama.V0_10_2_0
	}

	clusterConfig.Consumer.Return.Errors = true
	clusterConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	return clusterConfig
}

func (k *kBroker) getSaramaClusterClient(topic string) (sarama.Client, error) {
	config := k.getClusterConfig()
	cs, err := sarama.NewClient(k.addrs, config)
	if err != nil {
		return nil, err
	}
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	k.sc = append(k.sc, cs)
	return cs, nil
}

func (k *kBroker) getLogger() logger.Logger {
	return k.opts.Logger
}

func (k *kBroker) String() string {
	return "kafka broker implementation"
}