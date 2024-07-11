package mq

import (
	"crypto/tls"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/url"
	"sync"
)

// Init - init mq client:
//   - @param.uri: "{scheme[amqp/amqps]}://{username}:{password}@{endpoint}/{virtual_host}"
//   - @param.certs: with amqps, certs[0]=client crt bytes, certs[0]=client key bytes

type _client struct {
	sync.Mutex
	uri     string
	tlsCfg  *tls.Config
	conn    *amqp.Connection
	ch      *amqp.Channel
	consume <-chan amqp.Delivery
	queue   *queueOption
}

func (c *_client) open() error {
	var (
		err error
	)

	c.Lock()
	defer c.Unlock()

	if c.tlsCfg != nil {
		c.conn, err = amqp.DialTLS(c.uri, c.tlsCfg)
	} else {
		c.conn, err = amqp.Dial(c.uri)
	}

	if err != nil {
		return err
	}

	if c.ch, err = c.conn.Channel(); err != nil {
		return err
	}

	if client.queue != nil && client.queue.name != "" {
		if _, err = client.ch.QueueDeclare(
			client.queue.name,
			client.queue.durable,
			client.queue.autoDelete,
			client.queue.exclusive,
			client.queue.noWait,
			client.queue.args,
		); err != nil {
			return fmt.Errorf("declare queue: %s, err: %w", client.queue.name, err)
		}
	}

	return nil
}

var (
	client = &_client{
		uri:    "amqp://guest:guest@127.0.0.1:5672/",
		tlsCfg: nil,
	}
)

// Init - init mq client
func Init(opts ...OptionFn) error {
	var (
		err error
	)

	for _, fn := range opts {
		fn(client)
	}

	if _, err = url.Parse(client.uri); err != nil {
		return fmt.Errorf("url parse uri err: %w", err)
	}

	if err = client.open(); err != nil {
		return err
	}

	return nil
}
