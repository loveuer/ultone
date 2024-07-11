package invoke

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
	"ultone/internal/tool"
)

const (
	SCHEME = "sonar"
)

type Client[T any] struct {
	domain    string
	endpoints []string
	fn        func(grpc.ClientConnInterface) T
	opts      []grpc.DialOption

	cc *grpc.ClientConn
}

func (c *Client[T]) Session() T {
	return c.fn(c.cc)
}

var (
	clients = &sync.Map{}
)

// NewClient
/*
 * domain    => Example: sonar_search
 * endpoints => Example: []string{"sonar_search:8080", "sonar_search:80801"} or []string{"10.10.10.10:32000", "10.10.10.10:32001"}
 * fn        => Example: system.NewSystemSrvClient
 * opts      => Example: grpc.WithTransportCredentials(insecure.NewCredentials()),
 */
func NewClient[T any](
	domain string,
	endpoints []string,
	fn func(grpc.ClientConnInterface) T,
	opts ...grpc.DialOption,
) (*Client[T], error) {

	cached, ok := clients.Load(domain)
	if ok {
		if client, ok := cached.(*Client[T]); ok {
			return client, nil
		}
	}

	resolved := resolver.State{Addresses: make([]resolver.Address, 0)}

	locker.Lock()
	for _, item := range endpoints {
		resolved.Addresses = append(resolved.Addresses, resolver.Address{Addr: item})
	}
	ips[domain] = resolved
	locker.Unlock()

	fullAddress := fmt.Sprintf("%s://%s", SCHEME, domain)

	opts = append(opts,
		grpc.WithResolvers(myBuilder),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithChainUnaryInterceptor(retryInterceptor(3, 3*time.Second)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.DialContext(
		tool.Timeout(3),
		fullAddress,
		opts...,
	)

	if err != nil {
		return nil, err
	}

	c := &Client[T]{
		cc: conn,
		fn: fn,
	}

	clients.Store(domain, c)

	return c, nil
}
