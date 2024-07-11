package mq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func TestConsume(t *testing.T) {
	clientCert, err := tls.LoadX509KeyPair(
		"/Users/loveuer/codes/project/bifrost-pro/search_v3/internal/database/mq/tls/client.crt",
		"/Users/loveuer/codes/project/bifrost-pro/search_v3/internal/database/mq/tls/client.key",
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	ca, err := os.ReadFile("/Users/loveuer/codes/project/bifrost-pro/search_v3/internal/database/mq/tls/ca.crt")
	if err != nil {
		t.Fatal(err.Error())
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(ca) {
		t.Fatal("ca pool append ca crt err")
	}

	if err := Init(
		WithURI("amqps://admin:password@mq.dev:5671/export"),
		WithTLS(&tls.Config{
			Certificates:       []tls.Certificate{clientCert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}),
		WithQueueDeclare("export", false, false, false, false, amqp.Table{"x-max-priority": 100}),
	); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	ch, err := Consume(ctx, "export", &ConsumeOpt{MaxReconnection: -1})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[TEST] start consume msg")
	for msg := range ch {
		t.Logf("[TEST] [%s] [msg: %s]", time.Now().Format("060102T150405"), string(msg.Body))
		_ = msg.Ack(false)
	}
}

func TestPublish(t *testing.T) {
	clientCert, err := tls.LoadX509KeyPair(
		"/Users/loveuer/codes/project/bifrost-pro/search_v3/internal/database/mq/tls/client.crt",
		"/Users/loveuer/codes/project/bifrost-pro/search_v3/internal/database/mq/tls/client.key",
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	ca, err := os.ReadFile("/Users/loveuer/codes/project/bifrost-pro/search_v3/internal/database/mq/tls/ca.crt")
	if err != nil {
		t.Fatal(err.Error())
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(ca) {
		t.Fatal("ca pool append ca crt err")
	}

	if err := Init(
		WithURI("amqps://admin:password@mq.dev:5671/export"),
		WithTLS(&tls.Config{
			Certificates:       []tls.Certificate{clientCert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}),
		WithQueueDeclare("export", false, false, false, false, amqp.Table{"x-max-priority": 100}),
	); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	count := 1

	t.Log("[TEST] start publish msg...")

	for {
		if err = Publish(ctx, "export", amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(time.Now().Format(time.RFC3339) + " => hello_world@" + strconv.Itoa(count)),
		}); err != nil {
			t.Log(err.Error())
		}

		time.Sleep(11 * time.Second)
		count++
	}
}
