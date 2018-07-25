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
//return btcoin数, map[交易ID]输出Index
func (bc *BlockChain) getUnSpendInfo(address string, amount int32) (int32, map[string][]int32){
	outPuts := make(map[string][]int32, 0)
	unSpendTxs := bc.getUnSpendTransactions(address)
	total := int32(0)
	for _, tx:= range unSpendTxs{
		txId := hex.EncodeToString(tx.ID)
		for index, out := range tx.Outs{
			total += out.Value
			outPuts[txId] = append(outPuts[txId], index)
			if total >= amount{
				goto Complete
			}
		}
	}
	Complete:

	return total, outPuts
}

func (bc *BlockChain) getUnSpendTransactions(address string) []UpSpendTxs{
	bci := bc.Iterator()
	//记录未花费输出map[交易id]输出index
	unSpendOuts := make([]UpSpendTxs,0)
	//记录被引用的输出map[交易id]输出index
	spendOuts := make(map[string][]int32,0)
	for{
		b := bci.Next()
		//整理每笔交易的输入输出，交易是有先后顺序的.先有输出，才有输入
		for _, tx := range b.Transactions {
			txWithUpSpendOuts := UpSpendTxs{
				ID:tx.ID,
			}
			txID := hex.EncodeToString(tx.ID)
			//先处理输出
			for index, output := range tx.Vout {
				isSpend := false
				if _, ok := spendOuts[txID]; ok {
					for _, outindex := range spendOuts[txID] {
						if int32(index) == outindex {
							isSpend = true
							break
						}
					}
				}
				if !isSpend&&output.CanBeUnLockByAddr(address){
					txWithUpSpendOuts.Outs[int32(index)] = output
					unSpendOuts = append(unSpendOuts, txWithUpSpendOuts)
				}
			}

			//把这笔交易能解锁输出的输入记录下来
			for _, input := range tx.Vin {
				if input.CanUnLockOutPutByAddr(address){
					quoteTxId := hex.EncodeToString(input.TxId)
					spendOuts[quoteTxId] = append(spendOuts[quoteTxId], input.OutIndex)
				}
			}
		}

		if len(b.PreBlockHash) == 0 {
			break
		}
	}
	return unSpendOuts
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

func NewTransaction(from, to string, amount int32, bc BlockChain) *Transaction{
	inPuts := make([]TxInput, 0)
	outPuts := make([]TxOutput, 0)
	total, outPutInfos := bc.getUnSpendInfo(from, amount)
	if total < amount{
		fmt.Errorf("Not Enough Coins!")
		return nil
	}

	for txId, outIndexs := range outPutInfos{
		for _, index := range outIndexs{
			inPut := TxInput{
				TxId: []byte(txId),
				OutIndex: index,
				ScriptSig:from,
			}
			inPuts = append(inPuts, inPut)
		}
	}

	outPuts = append(outPuts, TxOutput{
		Value: amount,
		ScriptPubKey:to,
	})
	//如果有找零，再创建个输出
	if total > amount{
		outPuts = append(outPuts, TxOutput{
			Value: total - amount,
			ScriptPubKey: to,
		})
	}

	tx := &Transaction{
		Vin:inPuts,
		Vout:outPuts,
	}

	tx.SetID()

	return tx
}

func NewCoinBaseTX(to string, data string) *Transaction{
	outPuts := make([]TxOutput, 0)

	out := TxOutput{
		Value:misc.CoinAward,
		ScriptPubKey:to,
	}
}


