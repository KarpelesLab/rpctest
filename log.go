package rpctest

import (
	"context"
	"log"
)

type rpcLogPeer struct {
	RPC
}

func NewLogPeer(r RPC) RPC {
	if r == nil {
		log.Printf("[rpc] Instanciated with nil RPC, communications disabled")
		return nil
	}
	log.Printf("[rpc] Instanciated peer %s", r.Self())
	return &rpcLogPeer{RPC: r}
}

func (r *rpcLogPeer) All(ctx context.Context, data []byte) ([]any, error) {
	log.Printf("[rpc] %s → ALL: %x (all)", r.RPC.Self(), data)
	return r.RPC.All(ctx, data)
}

func (r *rpcLogPeer) Broadcast(ctx context.Context, data []byte) error {
	log.Printf("[rpc] %s → ALL: %x (responses ignored)", r.RPC.Self(), data)
	return r.RPC.Broadcast(ctx, data)
}

func (r *rpcLogPeer) Request(ctx context.Context, id string, data []byte) ([]byte, error) {
	log.Printf("[rpc] %s → %s: %x (request)", r.RPC.Self(), id, data)
	return r.RPC.Request(ctx, id, data)
}

func (r *rpcLogPeer) Send(ctx context.Context, id string, data []byte) error {
	log.Printf("[rpc] %s → %s: %x (send)", r.RPC.Self(), id, data)
	return r.RPC.Send(ctx, id, data)
}

func (r *rpcLogPeer) Connect(cb func(context.Context, []byte) ([]byte, error)) {
	log.Printf("[rpc] %s: connect endpoint", r.RPC.Self())
	r.RPC.Connect(cb)
}

func (r *rpcLogPeer) ListOnlinePeers() []string {
	log.Printf("[rpc] %s: listing online peers")
	return r.RPC.ListOnlinePeers()
}
