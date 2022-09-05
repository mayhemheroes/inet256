package dns256

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/brendoncarroll/go-p2p/p2ptest"
	"github.com/inet256/inet256/pkg/inet256"
	"github.com/inet256/inet256/pkg/inet256mem"
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
			{Key: "a", Value: NewValue("data"), TTL: 300},
		}
		return true
	})
	require.NotNil(t, res)
	require.Len(t, res.Entries, 1)
	require.Equal(t, res.Entries[0].Key, "a")
}

func TestResolveChain(t *testing.T) {
	nodes := makeNodes(t, 3)
	cnode := nodes[0]
	snodes := nodes[1:]

	ctx, cf := context.WithCancel(ctx)
	defer cf()
	eg := errgroup.Group{}
	for i, snode := range snodes {
		if i < len(snodes)-1 {
			eg.Go(func() error {
				return Serve(ctx, snode, func(res *Response, req *Request) bool {
					return true
				})
			})
		} else {
			eg.Go(func() error {
				return Serve(ctx, snode, func(res *Response, req *Request) bool {
					return true
				})
			})
		}
	}

	r := NewResolver(cnode, []inet256.Addr{snodes[0].LocalAddr()})
	ents, err := r.Resolve(ctx, Path{"a", "b", "c"}, []string{"key1"})
	require.NoError(t, err)
	require.Len(t, ents, 1)
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
