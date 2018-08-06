package main

import (
	"fmt"
	"blockChain"
	"cil"
	"misc"
)

func main(){
	bc := blockChain.OpenBlockChain(misc.Block_Chain_Path,"chain_1")
	if nil == bc{
		fmt.Println("open block chain failed!")
		return
	}
	cil := cil.NewCLI(bc)
	cil.Run()
}