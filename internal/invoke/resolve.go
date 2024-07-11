package invoke

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
)

const (
	scheme = "bifrost"
)

type CustomBuilder struct{}

func (cb *CustomBuilder) Scheme() string {
	return scheme
}

func (cb *CustomBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	cr := &customResolver{
		cc:     cc,
		target: target,
	}

	cr.ResolveNow(resolver.ResolveNowOptions{})

	return cr, nil
}

type customResolver struct {
	sync.Mutex
	target resolver.Target
	cc     resolver.ClientConn
	ips    map[string]string
}

func (cr *customResolver) ResolveNow(o resolver.ResolveNowOptions) {
	var (
		addrs = make([]resolver.Address, 0)
		hp    []string
	)

	cr.Lock()
	defer cr.Unlock()

	if hp = strings.Split(cr.target.URL.Host, ":"); len(hp) >= 2 {
		if ip, ok := pool[hp[0]]; ok {
			addr := fmt.Sprintf("%s:%s", ip, hp[1])
			addrs = append(addrs, resolver.Address{Addr: addr})
		}
	}

	_ = cr.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (cr *customResolver) Close() {}

var (
	cb   = &CustomBuilder{}
	pool = make(map[string]string)
)

func init() {
	resolver.Register(cb)
}

type CustomDomain struct {
	Domain string
	IP     string
}

func NewCustomBuilder(cds ...CustomDomain) resolver.Builder {
	locker.Lock()
	defer locker.Unlock()

	for _, cd := range cds {
		pool[cd.Domain] = cd.IP
	}

	return cb
}
