package invoke

import (
	"github.com/sirupsen/logrus"
	"sync"

	"google.golang.org/grpc/resolver"
)

type Builder struct{}

func (b *Builder) Scheme() string {
	return SCHEME
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	cr := &Resolver{
		cc:     cc,
		target: target,
	}

	cr.ResolveNow(resolver.ResolveNowOptions{})

	return cr, nil
}

type Resolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {
	logrus.Tracef("resolve_v2 ResolveNow => target: %s, %v", r.target.URL.Host, ips)
	_ = r.cc.UpdateState(ips[r.target.URL.Host])
}

func (cr *Resolver) Close() {}

var (
	locker    = &sync.Mutex{}
	myBuilder = &Builder{}
	ips       = map[string]resolver.State{}
)
