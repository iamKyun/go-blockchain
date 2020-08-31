package main

import (
	"iamkyun.com/go-blockchain/core"
)

func main() {
	bc := core.NewBlockchain()
	defer bc.DB.Close()

	cli := core.CLI{bc}
	cli.Run()
}
