package mq

import (
	"context"
	"fmt"
	"github.com/loveuer/esgo2dump/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"time"
	"ultone/internal/tool"
)

// ConsumeOpt
//   - Name: consumer's name, default unamed_<timestamp>
//   - MaxReconnection: when mq connection closed, max reconnection times, default 3, -1 for unlimited
type ConsumeOpt struct {
	Name            string // consumer's name, default unamed_<timestamp>
	AutoAck         bool
	Exclusive       bool
	NoLocal         bool
	NoWait          bool
	MaxReconnection int // when mq connection closed, max reconnection times, default 3, -1 for unlimited
	Args            amqp.Table
}

func Consume(ctx context.Context, queue string, opts ...*ConsumeOpt) (<-chan amqp.Delivery, error) {
	var (
		err error
		res = make(chan amqp.Delivery, 1)
		opt = &ConsumeOpt{
			Name:            os.Getenv("HOSTNAME"),
			AutoAck:         false,
			Exclusive:       false,
			NoLocal:         false,
			NoWait:          false,
			Args:            nil,
			MaxReconnection: 3,
		}
	)

	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	if opt.Name == "" {
		opt.Name = fmt.Sprintf("unamed_%d", time.Now().UnixMilli())
	}

	client.Lock()
	if client.consume, err = client.ch.Consume(queue, opt.Name, opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args); err != nil {
		client.Unlock()
		return nil, err
	}
	client.Unlock()

	go func() {
	Run:
		retry := 0
		for {
			select {
			case <-ctx.Done():
				close(res)
				return
			case m, ok := <-client.consume:
				if !ok {
					log.Warn("[mq] consume channel closed!!!")
					goto Reconnect
				}

				res <- m
			}
		}

	Reconnect:
		if opt.MaxReconnection == -1 || opt.MaxReconnection > retry {
			retry++

			log.Warn("[mq] try reconnect[%d/%d] to mq server after %d seconds...err: %v", retry, opt.MaxReconnection, tool.Min(60, retry*5), err)
			time.Sleep(time.Duration(tool.Min(60, retry*5)) * time.Second)
			if err = client.open(); err != nil {
				goto Reconnect
			}

			client.Lock()
			if client.consume, err = client.ch.Consume(queue, opt.Name, opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args); err != nil {
				client.Unlock()
				goto Reconnect
			}
			client.Unlock()

			log.Info("[mq] reconnect success!!!")
			goto Run
		}
	}()

	return res, nil
}
