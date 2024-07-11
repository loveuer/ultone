package mq

import (
	"crypto/tls"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OptionFn func(*_client)

// WithURI
//   - amqp uri, e.g. amqp://guest:guest@127.0.0.1:5672/vhost
//   - tips: with tls, scheme should be amqps, amqps://xx:xx@127.0.0.1:5672/vhost
func WithURI(uri string) OptionFn {
	return func(c *_client) {
		c.uri = uri
	}
}

// WithTLS
//   - amqps tls config
//   - include client cert, client key, ca cert
func WithTLS(tlsCfg *tls.Config) OptionFn {
	return func(c *_client) {
		c.tlsCfg = tlsCfg
	}
}

type queueOption struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

func WithQueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) OptionFn {
	return func(c *_client) {
		c.queue = &queueOption{
			name:       name,
			durable:    durable,
			autoDelete: autoDelete,
			exclusive:  exclusive,
			noWait:     noWait,
			args:       args,
		}
	}
}
