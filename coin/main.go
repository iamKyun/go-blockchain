package main

import (
	"fmt"
	"iamkyun.com/go-blockchain/core"
)

func main() {
	bc := core.NewBlockchain()
	bc.AddBlock("Kyun 1BTC Mike")
	bc.AddBlock("Mike 1.5BTC Kyun")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev: %x\n", block.PreBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Time: %d\n", block.Timestamp)
		fmt.Printf("Data: %s\n\n", block.Data)
	}
}
