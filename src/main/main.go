package main

import (
	"fmt"
	"blockChain"
	"cil"
)

func main(){
	bc := blockChain.OpenBlockChain("./db/blockchain.db","chain_1")
	if nil == bc{
		fmt.Println("open block chain failed!")
		return
	}
	cil := cil.NewCLI(bc)
	cil.Run()
}