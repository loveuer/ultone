package mq

import (
	"context"
	"errors"
	"github.com/loveuer/esgo2dump/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
	"ultone/internal/tool"
)

// PublishOpt
//   - MaxReconnect: publish msg auto retry with reconnect, should not be big, case memory leak
type PublishOpt struct {
	Exchange     string
	Mandatory    bool
	Immediate    bool
	MaxReconnect uint8 // publish msg auto retry with reconnect, should not be big(default 1), case memory leak
}

func Publish(ctx context.Context, queue string, msg amqp.Publishing, opts ...*PublishOpt) error {
	var (
		err error
		opt = &PublishOpt{
			Exchange:     amqp.DefaultExchange,
			Mandatory:    false,
			Immediate:    false,
			MaxReconnect: 1,
		}
		retry = 0
	)

	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	for ; retry <= int(opt.MaxReconnect); retry++ {
		if err = client.ch.PublishWithContext(ctx, opt.Exchange, queue, opt.Mandatory, opt.Immediate, msg); err == nil {
			return nil
		}

		if errors.Is(err, amqp.ErrClosed) {
			sleep := tool.Min(120, (retry+1)*30)

			log.Warn("[mq] connection closed, reconnect[%d/%d] after %d seconds", retry+1, opt.MaxReconnect, sleep)

			time.Sleep(time.Duration(sleep) * time.Second)

			if oerr := client.open(); oerr != nil {
				log.Error("[mq] reconnect[%d/%d] mq err: %v", oerr, retry+1, opt.MaxReconnect)
			} else {
				log.Info("[mq] reconnect mq success!!!")
			}

			continue
		}

		return err
	}

	return err
}
