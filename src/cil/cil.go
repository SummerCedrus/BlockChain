//handler cmd
package cil

import (
	"flag"
	"os"
	"fmt"
	"blockChain"
)

type CLI struct {
	bc *blockChain.BlockChain
}

func NewCLI(bc *blockChain.BlockChain) *CLI{
	return &CLI{bc:bc}
}

func (cli *CLI)Run(){
	addblock := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChain := flag.NewFlagSet("print", flag.ExitOnError)

	if len(os.Args)<1{
		fmt.Errorf("Wrong Arg Number!!!")
		return
	}
	switch os.Args[1] {
	case "addblock":
		addblock.Parse(os.Args[2:])
	case "printChain":
		printChain.Parse(os.Args[2:])
	default:
		fmt.Errorf("Error Cmd [%s]", os.Args[1])
		panic("")
	}

	if addblock.Parsed(){
		addblockData := addblock.String("data","", "block info")
		cli.addblock(*addblockData)
	}
	
	if printChain.Parsed(){
		cli.printChain()
	}
}

func (cil *CLI)printUsage(){

}
func (cil *CLI)addblock(data string) bool {
	bc := cil.bc
	err := bc.AddBlock([]byte(data))
	if nil != err{
		fmt.Errorf("Add Block failed!")
		return false
	}

	return true
}

func (cil *CLI)printChain()  {
	bc := cil.bc
	bc.Print()
}

