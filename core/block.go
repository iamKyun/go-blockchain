package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

//定义区块
type Block struct {
	//1.区块高度，也就是区块的编号，第几个区块
	Height int64
	//2.上一个区块的Hash值
	PreBlockHash []byte
	//3.交易数据（最终都属于transaction 事务）
	Data []byte
	//4.创建时间的时间戳
	Timestamp int64
	//5.当前区块的Hash值
	Hash []byte
	//6.Nonce 随机数，用于验证工作量证明
	Nonce int64
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func NewBlock(data string, height int64, PreBlockHash []byte) *Block {
	//根据传入参数创建区块
	block := &Block{
		height,
		PreBlockHash,
		[]byte(data),
		time.Now().Unix(),
		nil,
		0,
	}
	//调用工作量证明的方法，并且返回有效的Hash和Nonce值
	//创建pow对象
	pow := NewProofOfWork(block)
	//通过Run()方法进行挖矿验证
	hash, nonce := pow.Run(height)
	//将Nonce,Hash赋值给类对象属性
	block.Hash = hash[:]
	block.Nonce = nonce
	return block

}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", 1, []byte{})
}
