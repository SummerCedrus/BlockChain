package blockChain

import (
	"block"
	"fmt"
	."github.com/bolt"
	"misc"
	"errors"
	"bytes"
)

type BlockChain struct {
	tip []byte		//hash of last block
	db  *DB
	bucketName []byte
}

type BlockChainIter struct{
	currHash	[]byte
	db  *DB
	bucketName []byte
}

type CLI struct {
	bc *BlockChain
}
func (bc *BlockChain) AddBlock(data []byte) error{
	nb := block.NewBlock(bc.tip, data)
	db := bc.db
	err := db.Update(func(tx *Tx) error {
		bk := tx.Bucket(bc.bucketName)
		if nil == bk{
			return errors.New("bucket does not exist")
		}
		bk.Put(nb.Hash, misc.Serialize(nb))
		bk.Put([]byte("tip"), nb.Hash)

		return nil
	})

	bc.tip = nb.Hash
	return err
}

func (bc *BlockChain) Iterator() *BlockChainIter{
	return &BlockChainIter{currHash:bc.tip, db:bc.db, bucketName:bc.bucketName}
}

func (bc *BlockChain) Print() {
	iter := bc.Iterator()
	for {
		block := iter.Next()
		fmt.Printf("[%x] [%x] [%s]\n", block.PreBlockHash, block.Hash, string(block.Data))

		if 0 == bytes.Compare(block.PreBlockHash, []byte{}) {
			fmt.Println("chain end")
			break
		}
	}
}
func (bc *BlockChain) NewCLI() *CLI{
	return &CLI{bc:bc}
}
func (bci *BlockChainIter) Next() *block.Block{
	db := bci.db
	curBlock := new(block.Block)
	db.View(func(tx *Tx) error {
		bk := tx.Bucket(bci.bucketName)
		data := bk.Get(bci.currHash)
		if nil == data{
			return errors.New("can't find block")
		}
		err := misc.Deserialize(data, curBlock)

		if nil != err{
			fmt.Errorf("Deserialize failed error[%v]", err.Error())
			return err
		}

		return nil
	})

	bci.currHash = curBlock.PreBlockHash
	//fmt.Printf("[%x] [%x]\n",curBlock.PreBlockHash,curBlock.Hash)
	return curBlock
}

//func (cli *CLI)Run(){
//	addblock := flag.NewFlagSet("addblock", flag.ExitOnError)
//	print := flag.NewFlagSet("print", flag.ExitOnError)
//
//	if len(os.Args)>
//	switch os.Args[1] {
//	case "addblock":
//		addblock.Parse(os.Args[2:])
//	case "print":
//		print.Parse(os.Args[2:])
//	}
//
//}

func OpenBlockChain(filePath string, bucketName string) *BlockChain{
	bc := new(BlockChain)
	db, err := Open(filePath, 0600, nil)
	if nil != err||nil == db{
		fmt.Printf("Open db [%s] failed error[%s]!",filePath, err.Error())
		panic(fmt.Sprintf("Open db [%s] failed error[%s]!",filePath, err.Error()))
		return nil
	}
	bk := new(Bucket)
	err = db.Update(func(tx *Tx) error {
		bk = tx.Bucket([]byte(bucketName))
		if nil == bk{
			bk, err = tx.CreateBucket([]byte(bucketName))
			genBlock := block.NewGenesisBlock()
			bk.Put([]byte("tip"), genBlock.Hash)
			bk.Put(genBlock.Hash, misc.Serialize(genBlock))
			bc.tip = genBlock.Hash
			return err
		}else{
			bc.tip = bk.Get([]byte("tip"))
		}
		return nil
	})
	bc.db = db
	bc.bucketName = []byte(bucketName)
	if nil != err{
		fmt.Errorf("Open Bucket [%d] failed error[%s]", bucketName, err.Error())
		return nil
	}

	return bc
}

