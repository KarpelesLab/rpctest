package rpctest

import (
	"context"
	"io/fs"
)

type rpcAsyncPool map[string]*rpcAsyncPeer

func NewAsync() Pool {
	return make(rpcAsyncPool)
}

func (p rpcAsyncPool) NewPeer(id string) RPC {
	n := &rpcAsyncPeer{
		pool: p,
		id:   id,
	}
	p[id] = n
	return n
}

type rpcAsyncPeer struct {
	pool   rpcAsyncPool
	id     string
	target func(context.Context, []byte) ([]byte, error)
}

func (r *rpcAsyncPeer) All(ctx context.Context, data []byte) ([]any, error) {
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

func (r *rpcAsyncPeer) Broadcast(ctx context.Context, data []byte) error {
	for _, p := range r.pool {
		go p.run(ctx, data)
	}
	return nil
}

func (r *rpcAsyncPeer) Request(ctx context.Context, id string, data []byte) ([]byte, error) {
	p, ok := r.pool[id]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return p.run(ctx, data)
}

func (r *rpcAsyncPeer) Send(ctx context.Context, id string, data []byte) error {
	if p, ok := r.pool[id]; ok {
		go p.run(ctx, data)
	}
	return nil
}

func (r *rpcAsyncPeer) Self() string {
	return r.id
}

func (r *rpcAsyncPeer) Connect(cb func(context.Context, []byte) ([]byte, error)) {
	r.target = cb
}

func (r *rpcAsyncPeer) ListOnlinePeers() []string {
	res := make([]string, 0, len(r.pool))
	for p := range r.pool {
		res = append(res, p)
	}
	return res
}

func (r *rpcAsyncPeer) CountAllPeers() int {
	return len(r.pool)
}

func (r *rpcAsyncPeer) run(ctx context.Context, data []byte) ([]byte, error) {
	t := r.target
	if t != nil {
		return t(ctx, data)
	}
	return nil, fs.ErrNotExist
}
