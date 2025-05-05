package clients

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NATSOptions struct {
	URL           string
	Name          string
	MaxReconnect  int
	ReconnectWait time.Duration
}

var DefaultNATSOptions = &NATSOptions{
	URL:           nats.DefaultURL,
	Name:          "my-nats-client",
	MaxReconnect:  10,
	ReconnectWait: 2 * time.Second,
}

func NewNATSClient(opts *NATSOptions, logger *zap.Logger) (*nats.Conn, error) {
	if opts == nil {
		opts = DefaultNATSOptions
	}

	nc, err := nats.Connect(opts.URL,
		nats.Name(opts.Name),
		nats.MaxReconnects(opts.MaxReconnect),
		nats.ReconnectWait(opts.ReconnectWait),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			logger.Info("NATS disconnected", zap.String("reason", err.Error()))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS connected", zap.String("to", nc.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Info("NATS connection closed", zap.Error(nc.LastError()))
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			logger.Error("NATS connection error", zap.Error(err))
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return nc, nil
}

func (o *NATSOptions) WithURL(urls string) *NATSOptions {
	o.URL = urls
	return o
}

func (o *NATSOptions) WithName(name string) *NATSOptions {
	o.Name = name
	return o
}

func (o *NATSOptions) WithMaxReconnect(max int) *NATSOptions {
	o.MaxReconnect = max
	return o
}

func (o *NATSOptions) WithReconnectWait(reconnect time.Duration) *NATSOptions {
	o.ReconnectWait = reconnect
	return o
}
