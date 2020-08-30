package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const (
	targetBits = 20
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{
		block:  block,
		target: target,
	}
}

func (w ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Calculating the block contains [%s]\n ", w.block.Data)
	for nonce < maxNonce {
		data := w.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(w.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

func (w ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			w.block.PreBlockHash,
			w.block.Data,
			IntToHex(w.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{})
}

func (w ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := w.prepareData(w.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(w.target) == -1
}
