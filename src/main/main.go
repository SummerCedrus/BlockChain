package main

import (
	"blockChain"
	"fmt"
)

func main(){
	bc := blockChain.OpenBlockChain("./db/blockchain.db","chain_1")
	if nil == bc{
		fmt.Println("open block chain failed!")
		return
	}
	bc.AddBlock([]byte("Send Chengxs One coin"))

	bc.AddBlock([]byte("Send Chengxs Tow coin"))

	bc.Print()
}