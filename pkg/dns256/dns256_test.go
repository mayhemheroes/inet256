package dns256

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/brendoncarroll/go-p2p/p2ptest"
	"github.com/inet256/inet256/pkg/inet256"
	"github.com/inet256/inet256/pkg/inet256mem"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

var ctx = context.Background()

func TestParseDots(t *testing.T) {
	for _, tc := range []struct {
		In  string
		Out Path
		Err error
	}{
		{In: "", Out: Path{}},
		{In: "com", Out: Path{"com"}},
		{In: "www.example.com", Out: Path{"com", "example", "www"}},
	} {
		y, err := ParseDots(tc.In)
		if err != nil {
			require.Nil(t, y)
			require.Equal(t, tc.Err, err)
		} else {
			require.Equal(t, tc.Out, y)
		}
	}
}

func TestClientServer(t *testing.T) {
	n1, n2 := newPair(t)
	res := doOne(t, n1, n2, Request{
		Query: &Query{
			Path: MustParseDots("www.example.com"),
		},
	}, func(res *Response, req *Request) bool {
		res.Entries = []Entry{
			{Value: NewValue("data"), TTL: 300},
		}
		return true
	})
	require.NotNil(t, res)
	t.Log(res.Entries)
	require.Len(t, res.Entries, 1)
	require.Equal(t, must(res.Entries[0].AsString()), "data")
}

func TestResolveChain(t *testing.T) {
	nodes := makeNodes(t, 4)
	cnode := nodes[0]
	snodes := nodes[1:]

	ctx, cf := context.WithCancel(ctx)
	defer cf()
	eg := errgroup.Group{}
	eg.Go(func() error {
		return Serve(ctx, snodes[0], func(res *Response, req *Request) bool {
			if HasPrefix(req.Query.Path, Path{"a"}) {
				res.Next = &Redirect{
					Addrs: []inet256.Addr{
						snodes[1].LocalAddr(),
					},
					Prefix: Path{"a"},
				}
				return true
			}
			return false
		})
	})
	eg.Go(func() error {
		return Serve(ctx, snodes[1], func(res *Response, req *Request) bool {
			if HasPrefix(req.Query.Path, Path{"b"}) {
				res.Next = &Redirect{
					Addrs: []inet256.Addr{
						snodes[2].LocalAddr(),
					},
					Prefix: Path{"b"},
				}
				return true
			}
			return false
		})
	})
	eg.Go(func() error {
		return Serve(ctx, snodes[2], func(res *Response, req *Request) bool {
			if slices.Equal(req.Query.Path, Path{"c"}) {
				res.Entries = []Entry{
					{Value: NewValue("hello world"), TTL: 300},
				}
				return true
			}
			return false
		})
	})

	r := NewResolver(cnode, []inet256.Addr{snodes[0].LocalAddr()})
	ents, err := r.Resolve(ctx, Path{"a", "b", "c"})
	require.NoError(t, err)
	t.Log(ents)
	require.Len(t, ents, 1)
	require.Equal(t, "hello world", must(ents[0].AsString()))
	cf()
	require.NoError(t, eg.Wait())
}

func newPair(t testing.TB) (n1, n2 inet256.Node) {
	ns := makeNodes(t, 2)
	return ns[0], ns[1]
}

func makeNodes(t testing.TB, n int) (ns []inet256.Node) {
	s := inet256mem.New()
	ns = make([]inet256.Node, n)
	var err error
	for i := range ns {
		pk := p2ptest.NewTestKey(t, i)
		ns[i], err = s.Open(ctx, pk)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, s.Drop(ctx, pk)) })
	}
	return ns
}

func doOne(t testing.TB, cnode, snode inet256.Node, req Request, h Handler) Response {
	ctx, cf := context.WithCancel(ctx)
	defer cf()
	c := NewClient(cnode)
	eg := errgroup.Group{}
	eg.Go(func() error {
		return Serve(ctx, snode, h)
	})
	var res Response
	eg.Go(func() error {
		defer cf()
		req := Request{
			Query: &Query{
				Path: MustParseDots("www.example.com"),
			},
		}
		res2, err := c.Do(ctx, snode.LocalAddr(), req)
		if err != nil {
			return err
		}
		res = *res2
		return nil
	})
	require.NoError(t, eg.Wait())
	return res
}

func must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}
