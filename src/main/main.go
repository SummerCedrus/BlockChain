package main

import "blockChain"

func main(){
	bc := blockChain.NewBlockChain()

	bc.AddBlock([]byte("Send Chengxs One coin"))

	bc.AddBlock([]byte("Send Chengxs Tow coin"))

	bc.Print()
}