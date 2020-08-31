package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
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

	target = target.Lsh(target, 256-targetBits)

	return &ProofOfWork{block, target}
}

func (proofOfWork *ProofOfWork) Run(num int64) ([]byte, int64) {
	//初始化随机数nonce为0
	nonce := 0
	//存储新生成的hash值
	var hashInt big.Int
	//存储hash值
	var hash [32]byte

	for {
		//1. 将Block的属性拼接成字节数组，注意，参数为Nonce
		databytes := proofOfWork.prepareData(nonce)
		//2.将拼接后的字节数组生成Hash
		hash = sha256.Sum256(databytes)
		//3. 将hash存储至hashInt
		hashInt.SetBytes(hash[:])
		//4.判断hashInt是否小于Block里面的Target
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//此处需要满足hashInt(y)小于设置的target(x)即 x > y,则挖矿成功
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			fmt.Printf("第%d个区块，挖矿成功:%x\n", num, hash)
			fmt.Println(time.Now())
			time.Sleep(time.Second * 2)
			//挖矿成功，退出循环
			break
		}
		nonce++
	}

	return hash[:], int64(nonce)

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

func (w ProofOfWork) IsValid() bool {
	var hashInt big.Int

	hashInt.SetBytes(w.block.Hash)

	if w.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}
