package rpctest

import (
	"context"
	"io/fs"
	"log"
)

type rpcSyncLogPool map[string]*rpcSyncLogPeer

func NewSyncLog() Pool {
	return make(rpcSyncLogPool)
}

func (p rpcSyncLogPool) NewPeer(id string) RPC {
	log.Printf("[rpc] new peer %s", id)
	n := &rpcSyncLogPeer{
		pool: p,
		id:   id,
	}
	p[id] = n
	return n
}

type rpcSyncLogPeer struct {
	pool   rpcSyncLogPool
	id     string
	target func(context.Context, []byte) ([]byte, error)
}

func (r *rpcSyncLogPeer) All(ctx context.Context, data []byte) ([]any, error) {
	log.Printf("[rpc] %s → ALL: %x (all)", r.id, data)
	var res []any
	for _, p := range r.pool {
		buf, e := p.run(ctx, data)
		if e != nil {
			log.Printf("[rpc] %s → %s: error %s", p.id, r.id, e)
			res = append(res, e)
		} else {
			res = append(res, buf)
		}
	}
	return res, nil
}

func (r *rpcSyncLogPeer) Broadcast(ctx context.Context, data []byte) error {
	log.Printf("[rpc] %s → ALL: %x (responses ignored)", r.id, data)
	for _, p := range r.pool {
		p.run(ctx, data)
	}
	return nil
}

func (r *rpcSyncLogPeer) Request(ctx context.Context, id string, data []byte) ([]byte, error) {
	log.Printf("[rpc] %s → %s: %x (request)", r.id, id, data)
	p, ok := r.pool[id]
	if !ok {
		return nil, fs.ErrNotExist
	}
	res, err := p.run(ctx, data)
	if err != nil {
		log.Printf("[rpc] %s → %s: error %s (response)", p.id, r.id, err)
	} else {
		log.Printf("[rpc] %s → %s: %x (response)", p.id, r.id, res)
	}
	return res, err
}

func (r *rpcSyncLogPeer) Send(ctx context.Context, id string, data []byte) error {
	log.Printf("[rpc] %s → %s: %x (send)", r.id, id, data)
	if p, ok := r.pool[id]; ok {
		p.run(ctx, data)
	}
	return nil
}

func (r *rpcSyncLogPeer) Self() string {
	log.Printf("[rpc] %s: who am I?", r.id)
	return r.id
}

func (r *rpcSyncLogPeer) Connect(cb func(context.Context, []byte) ([]byte, error)) {
	log.Printf("[rpc] %s: connect endpoint", r.id)
	r.target = cb
}

func (r *rpcSyncLogPeer) ListOnlinePeers() []string {
	log.Printf("[rpc] %s: listing online peers")
	res := make([]string, 0, len(r.pool))
	for p := range r.pool {
		res = append(res, p)
	}
	return res
}

func (r *rpcSyncLogPeer) CountAllPeers() int {
	return len(r.pool)
}

func (r *rpcSyncLogPeer) run(ctx context.Context, data []byte) ([]byte, error) {
	t := r.target
	if t != nil {
		return t(ctx, data)
	}
	return nil, fs.ErrNotExist
}
