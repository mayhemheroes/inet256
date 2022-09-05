package dns256

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/inet256/inet256/pkg/inet256"
	"golang.org/x/crypto/sha3"
)

// Handler modifies res in response to req and returns true if a message should be sent in response.
// The res.RequestID will be set automatically
type Handler func(res *Response, req *Request) bool

// Serve uses h to handle and respond to requests received by node until a
// non-transient error occurs.
// Such an error could come from the Node, or the Context.
// Serve returns nil, only if node returns context.Cancelled.
func Serve(ctx context.Context, node inet256.Node, h Handler) error {
	for {
		if err := node.Receive(ctx, func(msg inet256.Message) {
			var req Request
			if err := json.Unmarshal(msg.Payload, &req); err != nil {
				return
			}
			var reqID RequestID
			sha3.ShakeSum256(reqID[:], msg.Payload)
			replyTo := msg.Src
			go func() {
				var res Response
				if h(&res, &req) {
					res.RequestID = reqID
					data, err := json.Marshal(res)
					if err != nil {
						panic(err)
					}
					if err := node.Send(ctx, replyTo, data); err != nil {
						log.Println(err)
					}
				}
			}()
		}); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
	}
}
