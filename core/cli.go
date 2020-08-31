package core

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

type CLI struct {
	Bc *Blockchain
}

func (b *CLI) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("    add -data BLOCK_DATA : add a block to the blockchain ")
	fmt.Println("    print : print all the blocks of the blockchain")
}

func (b *CLI) validateArgs() {
	if len(os.Args) < 2 {
		b.printUsage()
		os.Exit(1)
	}
}

func (b *CLI) addBlock(data string) {
	b.Bc.AddBlock(data)
	fmt.Println("Succeeded!")
}

func (b *CLI) Run() {
	b.validateArgs()
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)

	addBlockData := addCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		b.printUsage()
		os.Exit(1)
	}

	if addCmd.Parsed() {
		if *addBlockData == "" {
			addCmd.Usage()
			os.Exit(1)
		}
		b.addBlock(*addBlockData)
	}

	if printCmd.Parsed() {
		b.print()
	}
}

func (b *CLI) print() {
	//创建迭代类对象blockchainIterator
	blockchainIterator := b.Bc.Iterator()

	for {
		//通过Next()方法获取当前区块，并更新区块链对象保存的Hash值为上一个区块的Hash值
		block := blockchainIterator.Next()

		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PreBlockHash:%x\n", block.PreBlockHash)
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)

		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}

	}

}
