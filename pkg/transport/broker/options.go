package broker

import (
	"context"
	"crypto/tls"

	"github.com/lengocson131002/go-clean/pkg/logger"
)

type BrokerOption func(*BrokerOptions)

type BrokerOptions struct {
	Context context.Context

	// underlying logger
	Logger logger.Logger

	// Handler executed when error happens in broker mesage
	// processing
	ErrorHandler Handler

	Addrs []string

	TLSConfig *tls.Config
}

func WithBrokerContext(ctx context.Context) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Context = ctx
	}
}

func WithBrokerAddresses(addrs ...string) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Addrs = addrs
	}
}

func WithLogger(log logger.Logger) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Logger = log
	}
}

func WithBrokerErrorHandler(handler Handler) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.ErrorHandler = handler
	}
}

func WithBrokerTLSConfig(t *tls.Config) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.TLSConfig = t
	}
}

type PublishOption func(*PublishOptions)

type PublishOptions struct {
	Context context.Context
}

func WithPublishContext(ctx context.Context) PublishOption {
	return func(opts *PublishOptions) {
		opts.Context = ctx
	}
}

type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	Context context.Context

	Queue string

	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
}

func WithSubscribeContext(ctx context.Context) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Context = ctx
	}
}

func WithSubscribeQueue(q string) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Queue = q
	}
}

func WithSubscribeAutoAck(autoAck bool) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.AutoAck = autoAck
	}
}
