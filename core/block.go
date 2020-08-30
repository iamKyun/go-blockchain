package core

import (
	"time"
)

type Block struct {
	Timestamp    int64
	Data         []byte
	PreBlockHash []byte
	Hash         []byte
	Nonce        int
}

func NewBlock(data string, preBlockHash []byte) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		Data:         []byte(data),
		PreBlockHash: preBlockHash,
		Hash:         []byte{},
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
