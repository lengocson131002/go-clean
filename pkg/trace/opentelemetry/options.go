package opentelemetry

type Options struct {
	ServiceName string
	Endpoint    string
}

type Option func(*Options)

func WithServiceName(sn string) Option {
	return func(o *Options) {
		o.ServiceName = sn
	}
}

func WithEndpoint(ep string) Option {
	return func(o *Options) {
		o.Endpoint = ep
	}
}
