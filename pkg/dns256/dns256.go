package dns256

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"regexp"
	"strings"

	"github.com/brendoncarroll/go-tai64"
	"github.com/inet256/inet256/pkg/inet256"
)

type Path []string

var validElem = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

func ParseDots(x string) (Path, error) {
	if x == "" {
		return Path{}, nil
	}
	parts := strings.Split(x, ".")
	for _, elem := range parts {
		if !validElem.MatchString(elem) {
			return nil, fmt.Errorf("%q is not a valid path element", elem)
		}
	}
	for i := 0; i < len(parts)/2; i++ {
		j := len(parts) - 1 - i
		parts[i], parts[j] = parts[j], parts[i]
	}
	return Path(parts), nil
}

func MustParseDots(x string) Path {
	p, err := ParseDots(x)
	if err != nil {
		panic(err)
	}
	return p
}

type RequestID [16]byte

func (id RequestID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id[:])
}

func (id *RequestID) UnmarshalJSON(x []byte) error {
	var s []byte
	if err := json.Unmarshal(x, &s); err != nil {
		return err
	}
	copy(id[:], s)
	return nil
}

type Request struct {
	Query *Query `json:"query"`
}

type Query struct {
	Path   Path   `json:"path"`
	Filter Filter `json:"filter"`
}

type Filter map[string]string

type Response struct {
	// RequestID is the ID of the corresponding Request
	RequestID RequestID `json:"req_id"`
	// Now is the server time
	Now tai64.TAI64 `json:"now"`

	Next    *Redirect `json:"next,omitempty"`
	Entries []Entry   `json:"entries,omitempty"`
}

type Redirect struct {
	Addrs  []inet256.Addr `json:"addrs"`
	Prefix Path           `json:"prefix"`
	// TTL is the time to live in seconds
	TTL uint32 `json:"ttl"`
}

type Entry struct {
	Data map[string]json.RawMessage `json:"data"`
	// TTL is the time to live in seconds
	TTL uint32 `json:"ttl"`
}

func (e Entry) AsUint64(key string) (ret uint64, err error) {
	err = json.Unmarshal(e.Data[key], &ret)
	return ret, err
}

func (e Entry) AsString(key string) (ret string, err error) {
	err = json.Unmarshal(e.Data[key], &ret)
	return ret, err
}

func (e Entry) AsINET256(key string) (ret inet256.Addr, err error) {
	err = json.Unmarshal(e.Data[key], &ret)
	return ret, err
}

func (e Entry) AsIP() (netip.Addr, error) {
	s, err := e.AsString()
	if err != nil {
		return netip.Addr{}, err
	}
	return netip.ParseAddr(s)
}

func (e Entry) AsIPPort() (netip.AddrPort, error) {
	s, err := e.AsString()
	if err != nil {
		return netip.AddrPort{}, err
	}
	return netip.ParseAddrPort(s)
}

func NewValue(x any) json.RawMessage {
	data, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return data
}
