package blockChain

import (
	"block"
	"fmt"
	."github.com/bolt"
	"misc"
	"errors"
	"bytes"
	."tx"
	"encoding/hex"
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

//找出够消耗的未花费的输出
//return map[交易ID]输出Index
func (bc *BlockChain) getUnSpendList(address string, amount int32) (int32, map[string]int32){
	outPuts := make(map[string]int32, 0)
	return 0, outPuts
}

func (bc *BlockChain) getUnSpendTransactions(address string) []*Transaction{
	bci := bc.Iterator()
	//记录被引用的输出map[交易id]输出index
	spendOuts := make(map[string]int,0)
	for{
		b := bci.Next()
		//整理每笔交易的输入输出，交易是有先后顺序的.先有输出，才有输入
		for _, tx := range b.Transactions {
			txID := hex.EncodeToString(tx.ID)
			//先处理输出
			for index, output := range tx.Vout {
				isSpend := false
				if _, ok := spendOuts[txID]; ok {
					for _, outindex := range spendOuts[txID] {
						if index == outindex {
							isSpend = true
							break
						}
					}
				}
				if isSpend{

				}
			}
		}
		if len(b.PreBlockHash) == 0 {
			break
		}
	}

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

