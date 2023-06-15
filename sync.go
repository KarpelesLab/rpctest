package rpctest

import (
	"context"
	"io/fs"
)

type rpcSyncPool map[string]*rpcSyncPeer

func NewSync() Pool {
	return make(rpcSyncPool)
}

func (p rpcSyncPool) NewPeer(id string) RPC {
	n := &rpcSyncPeer{
		pool: p,
		id:   id,
	}
	p[id] = n
	return n
}

type rpcSyncPeer struct {
	pool   rpcSyncPool
	id     string
	target func(context.Context, []byte) ([]byte, error)
}

func (r *rpcSyncPeer) All(ctx context.Context, data []byte) ([]any, error) {
	var res []any
	for _, p := range r.pool {
		r, e := p.run(ctx, data)
		if e != nil {
			res = append(res, e)
		} else {
			res = append(res, r)
		}
	}
	return res, nil
}

func (r *rpcSyncPeer) Broadcast(ctx context.Context, data []byte) error {
	for _, p := range r.pool {
		p.run(ctx, data)
	}
	return nil
}

func (r *rpcSyncPeer) Request(ctx context.Context, id string, data []byte) ([]byte, error) {
	p, ok := r.pool[id]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return p.run(ctx, data)
}

func (r *rpcSyncPeer) Send(ctx context.Context, id string, data []byte) error {
	if p, ok := r.pool[id]; ok {
		p.run(ctx, data)
	}
	return nil
}

func (r *rpcSyncPeer) Self() string {
	return r.id
}

func (r *rpcSyncPeer) Connect(cb func(context.Context, []byte) ([]byte, error)) {
	r.target = cb
}

func (r *rpcSyncPeer) ListOnlinePeers() []string {
	res := make([]string, 0, len(r.pool))
	for p := range r.pool {
		res = append(res, p)
	}
	return res
}

func (r *rpcSyncPeer) CountAllPeers() int {
	return len(r.pool)
}

func (r *rpcSyncPeer) run(ctx context.Context, data []byte) ([]byte, error) {
	t := r.target
	if t != nil {
		return t(ctx, data)
	}
	return nil, fs.ErrNotExist
}
