package dns256

import (
	"context"
	"net/netip"

	"github.com/inet256/inet256/pkg/inet256"
)

// Resolver performs path resolution
type Resolver struct {
	client Client
	roots  []inet256.Addr
}

func NewResolver(node inet256.Node, roots []inet256.Addr) *Resolver {
	return &Resolver{
		client: *NewClient(node),
		roots:  roots,
	}
}

func (r *Resolver) Resolve(ctx context.Context, x Path, keys []string) ([]Entry, error) {
	roots := append([]inet256.Addr{}, r.roots...)
	for len(roots) > 0 {
		var target inet256.Addr
		target, roots = roots[len(roots)-1], roots[:len(roots)-1]
		res, err := r.client.Do(ctx, target, Request{
			Query: &Query{
				Path: x,
				Keys: keys,
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
			x = res.Next.Path
			roots = res.Next.Addrs
		}
	}
	return nil, nil
}

func (r *Resolver) ResolveINET256(ctx context.Context, x Path) ([]inet256.Addr, error) {
	ents, err := r.Resolve(ctx, x, []string{"INET256"})
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
	ents, err := r.Resolve(ctx, x, []string{"AAAA", "AA"})
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
	ents, err := r.Resolve(ctx, x, []string{"SVC"})
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
