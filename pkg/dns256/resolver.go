package dns256

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/inet256/inet256/pkg/inet256"
)

// Resolver performs path resolution
type Resolver struct {
	client Client
	roots  []inet256.Addr

	cache *Cache
}

func NewResolver(node inet256.Node, roots []inet256.Addr, opts ...ResolverOpt) *Resolver {
	return &Resolver{
		client: *NewClient(node),
		roots:  roots,
		cache:  NewCache(),
	}
}

func (r *Resolver) Resolve(ctx context.Context, x Path, opts ...ResolveOpt) ([]Entry, error) {
	config := resolveConfig{
		maxHops: 32,
	}
	roots := append([]inet256.Addr{}, r.roots...)
	for i := 0; len(roots) > 0; i++ {
		if i >= config.maxHops {
			return nil, fmt.Errorf("max hops exceeded")
		}
		var target inet256.Addr
		target, roots = roots[len(roots)-1], roots[:len(roots)-1]
		res, err := r.client.Do(ctx, target, Request{
			Query: &Query{
				Path: x,
			},
		})
		if err != nil {
			return nil, err
		}
		switch {
		case len(res.Entries) > 0:
			return res.Entries, nil
		case res.Next != nil:
			// TODO: add checks here.
			if !HasPrefix(x, res.Next.Prefix) {
				return nil, fmt.Errorf("invalid redirect %v from %v", res.Next, target)
			}
			x = x[len(res.Next.Prefix):]
			roots = res.Next.Addrs
		}
	}
	return nil, nil
}

func (r *Resolver) ResolveINET256(ctx context.Context, x Path) ([]inet256.Addr, error) {
	ents, err := r.Resolve(ctx, x, WithLabels(map[string]string{"type": "INET256"}))
	if err != nil {
		return nil, err
	}
	var ret []inet256.Addr
	for _, ent := range ents {
		addr, err := ent.AsINET256()
		if err != nil {
			continue
		}
		ret = append(ret, addr)
	}
	return ret, nil
}

func (r *Resolver) ResolveIP(ctx context.Context, x Path) ([]netip.Addr, error) {
	ents, err := r.Resolve(ctx, x, WithLabels(map[string]string{"type": "AA"}))
	if err != nil {
		return nil, err
	}
	var ret []netip.Addr
	for _, ent := range ents {
		addr, err := ent.AsIP()
		if err != nil {
			continue
		}
		ret = append(ret, addr)
	}
	return ret, nil
}

func (r *Resolver) ResolveIPPort(ctx context.Context, x Path) ([]netip.AddrPort, error) {
	ents, err := r.Resolve(ctx, x, WithLabels(map[string]string{"type": "SVC"}))
	if err != nil {
		return nil, err
	}
	var ret []netip.AddrPort
	for _, ent := range ents {
		ap, err := ent.AsIPPort()
		if err != nil {
			continue
		}
		ret = append(ret, ap)
	}
	return ret, nil
}

func HasSuffix[E comparable, S ~[]E](x, s S) bool {
	for i := len(s) - 1; i >= 0; i++ {
		if i >= len(x) {
			return false
		}
		if s[i] != x[i] {
			return false
		}
	}
	return true
}

func HasPrefix[E comparable, S ~[]E](x, p S) bool {
	for i := 0; i < len(p); i++ {
		if i >= len(x) {
			return false
		}
		if p[i] != x[i] {
			return false
		}
	}
	return true
}
