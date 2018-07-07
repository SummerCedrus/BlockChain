package blockChain

import (
	"block"
	"fmt"
)

type BlockChain struct {
	Blocks []*block.Block
}

func (bc *BlockChain) AddBlock(data []byte){
	lastBlock := bc.Blocks[len(bc.Blocks) - 1]

	nb := block.NewBlock(lastBlock.PreBlockHash, data)

	bc.Blocks = append(bc.Blocks, nb)
}

func (bc *BlockChain) Print(){
	fmt.Println("Block Length: ", len(bc.Blocks))
	for _, b := range bc.Blocks{
		fmt.Println(string(b.Data))
	}
}

func NewBlockChain() *BlockChain{
	return &BlockChain{Blocks:[]*block.Block{block.NewGenesisBlock()}}
}

