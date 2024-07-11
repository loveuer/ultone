package invoke

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func retryInterceptor(maxAttempt int, interval time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		if maxAttempt == 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		duration := interval

		for attempt := 1; attempt <= maxAttempt; attempt++ {

			if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
					logrus.Debugf("Connection failed err: %v, retry %d after %fs", err, attempt, duration.Seconds())

					time.Sleep(duration)
					duration *= 2

					continue
				}

				return err
			}

			return nil // 请求成功，不需要重试
		}

		return fmt.Errorf("max retry attempts reached")
	}
}
