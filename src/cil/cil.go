//handler cmd
package cil

import (
	"flag"
	"os"
	"fmt"
	"blockChain"
	"tx"
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
	balance := flag.NewFlagSet("balance", flag.ExitOnError)
	send := flag.NewFlagSet("send", flag.ExitOnError)

	if len(os.Args)<2{
		fmt.Errorf("Wrong Arg Number!!!")
		return
	}
	switch os.Args[1] {
	case "addblock":
		addblock.Parse(os.Args[2:])
	case "printChain":
		printChain.Parse(os.Args[2:])
	case "balance":
		balance.Parse(os.Args[2:])
	case "send":
		send.Parse(os.Args[2:])
	default:
		fmt.Errorf("Error Cmd [%s]", os.Args[1])
		panic("")
	}

	//if addblock.Parsed(){
	//	addblockData := addblock.String("data","", "block info")
	//	cli.addblock(*addblockData)
	//}
	
	if printChain.Parsed(){
		cli.printChain()
	}

	if balance.Parsed(){
		address := addblock.String("addr","", "adress")
		cli.getBalance(*address)
	}
	if send.Parsed(){
		from := send.String("from", "", "from address")
		to := send.String("to","", "to address")
		amount := send.Int("amount", 0, "send coin amount")
		cli.send(*from, *to, int32(*amount))
	}
}

func (cil *CLI)printUsage(){

}
//func (cil *CLI)addblock(data string) bool {
//	bc := cil.bc
//	err := bc.AddBlock([]byte(data))
//	if nil != err{
//		fmt.Errorf("Add Block failed!")
//		return false
//	}
//
//	return true
//}

func (cil *CLI)printChain()  {
	bc := cil.bc
	bc.Print()
}

func (cil *CLI)getBalance(address string){
	bc := cil.bc
	fmt.Println("Balance:[%d]", bc.GetBalance(address))
}

func (cil *CLI)send(from, to string, amount int32){
	bc := cil.bc
	trans := bc.NewTransaction(from, to, amount)
	err := bc.AddBlock([]*tx.Transaction{trans})
	if nil != err{
		fmt.Errorf("send error [%s]",err.Error())
		return
	}

	fmt.Println("send SUCCESS")
}

