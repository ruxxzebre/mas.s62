package main

import "sync"

type Hash [32]byte

type Block struct {
	PrevHash Hash
	Name     string
	Nonce    string
}

// BlockChain is not actually a blockchain, it's just the tip.
// The chain itself only exists in a file.
type BlockChain struct {
	mtx   sync.Mutex
	tip   Block
	bchan chan Block
}