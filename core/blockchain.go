package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const (
	dbFile      = "blockchain.db"
	blockBucket = "blocks"
)

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

func (blc *Blockchain) AddBlock(data string) {
	//获取db对象blc.DB
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		//1.获取表
		b := tx.Bucket([]byte(blockBucket))
		//2.创建新区块
		if b != nil {
			//通过Key:blc.Tip获取Value(区块序列化字节数组)
			byteBytes := b.Get(blc.Tip)
			//反序列化出最新区块(上一个区块)对象
			block := DeserializeBlock(byteBytes)

			//3.通过NewBlock进行挖矿生成新区块newBlock
			newBlock := NewBlock(data, block.Height+1, block.Hash)
			//4.将最新区块序列化并且存储到数据库中(key=新区块的Hash值，value=新区块序列化)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			/*5.更新数据库中"l"对应的Hash为新区块的Hash值
			  用途:便于通过该Hash值找到对应的Block序列化，从而找到上一个Block对象，为生成新区块函数NewBlock提供高度Height与上一个区块的Hash值PreBlockHash
			*/
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			//6. 更新Tip值为新区块的Hash值
			blc.Tip = newBlock.Hash
		}

		return nil

	})
	if err != nil {
		log.Panic(err)
	}
}

type BlockChainIterator struct {
	CurrentHash []byte   // 保存当前的区块Hash值
	DB          *bolt.DB //DB对象
}

func (blockchainIterator *BlockChainIterator) Next() *Block {
	//1.定义Block对象block
	var block *Block
	//2.操作DB对象blockchainIterator.DB
	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		//3.打开表对象blockTableName
		b := tx.Bucket([]byte(blockBucket))

		if b != nil {
			//Get()方法通过Key:当前区块的Hash值获取当前区块的序列化信息
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)
			//反序列化出当前的区块
			block = DeserializeBlock(currentBlockBytes)
			//更新迭代器里面的CurrentHash
			blockchainIterator.CurrentHash = block.PreBlockHash
		}
		return nil

	})

	if err != nil {
		log.Panic(err)
	}
	return block
}

//定义函数DeserializeBlock()，传入参数为字节数组，返回值为Block
func DeserializeBlock(blockBytes []byte) *Block {
	//1.定义一个Block指针对象
	var block Block
	//2.初始化反序列化对象decoder
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	//3.通过Decode()进行反序列化
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	//4.返回block对象
	return &block
}

func (b *Blockchain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{b.Tip, b.DB}
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b == nil {
			fmt.Println("DB not found. Create new one")
			genesisBlock := NewGenesisBlock()

			tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("1"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesisBlock.Hash

		} else {
			tip = b.Get([]byte("1"))
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}
