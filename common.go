package rpctest

type Pool interface {
	NewPeer(id string) RPC
}
