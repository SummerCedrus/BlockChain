package blockChain

import (
	"fmt"
	"block"
	."github.com/bolt"
	."misc"
	."tx"
	"encoding/hex"
)

type UTXOSet struct {
	bc *BlockChain
}
//重建UTXO集
func (u *UTXOSet)ReBuild()  {
	db, err := Open(UTXO_Set_File_Path, 0600, nil)
	if nil != err||nil == db{
		panic(fmt.Sprintf("Open db [%s] failed error[%s]!",UTXO_Set_File_Path, err.Error()))
		return
	}
	utxo := u.bc.GetUTXO()
	err = db.Update(func(tx *Tx) error {
		err := tx.DeleteBucket([]byte(UTXO_Set_Name))
		bk, err:= tx.CreateBucket([]byte(UTXO_Set_Name))
		for _, txo := range utxo{
			err = bk.Put(txo.ID, Serialize(txo))
		}

		return err
	})

	if nil != err{
		fmt.Errorf("ReBulid Error [%s]", err.Error())
	}
}
//通过block的交易信息去更新UTXO集
func (u *UTXOSet)Update(block block.Block){
	db, err := Open(UTXO_Set_File_Path, 0600, nil)
	if nil != err||nil == db{
		panic(fmt.Sprintf("Open db [%s] failed error[%s]!",UTXO_Set_File_Path, err.Error()))
		return
	}

	err = db.Update(func(tx *Tx) error {
		bk := tx.Bucket([]byte(UTXO_Set_Name))
		if nil == bk{
			bk, err =tx.CreateBucket([]byte(UTXO_Set_Name))
		}

		for _, t := range block.Transactions{
			usp := UnSpendTxs{
				ID: t.ID,
				Outs:make(map[int32]TxOutput),
			}
			//将交易里面的新的输出更新到UTXO集
			for index, out := range t.Vout{
				usp.Outs[int32(index)] = out
			}

			err = bk.Put(t.ID, Serialize(usp))
			//将交易里面输入引用的输出从UTXO集删除
			if t.IsCoinBase(){
				continue
			}

			for _, in := range t.Vin{
				outsBytes := bk.Get(in.TxId)
				outs := new(UnSpendTxs)
				Deserialize(outsBytes, outs)
				delete(outs.Outs, in.OutIndex)
				//如果交易对应的未花费输出空了，把交易也删除
				if len(outs.Outs) == 0{
					bk.Delete(in.TxId)
				}else{
					bk.Put(in.TxId, Serialize(outs))
				}
			}
		}

		return err
	})

	if nil != err{
		fmt.Errorf("Update Error [%s]", err.Error())
	}
}

func (u *UTXOSet)getUnspendTxs(pubKeyHash []byte)[]UnSpendTxs{
	unSpendTxs := make([]UnSpendTxs, 0)
	db, err := Open(UTXO_Set_File_Path, 0600, nil)
	db.View(func(tx *Tx) error {
		bk := tx.Bucket([]byte(UTXO_Set_Name))
		if nil == bk{
			bk, err = tx.CreateBucket([]byte(UTXO_Set_Name))
			return err
		}
		c := bk.Cursor()
		for k,v := c.First();k != nil;k,v = c.Next(){
			unSpendTx := new(UnSpendTxs)
			Deserialize(v, unSpendTx)
			resUnSpendTx := UnSpendTxs{
				ID:unSpendTx.ID,
				Outs:make(map[int32]TxOutput),
			}
			for index, out := range unSpendTx.Outs{
				if out.IsLockedByPubKeyHash(pubKeyHash){
					resUnSpendTx.Outs[index] = out
				}
			}
			unSpendTxs = append(unSpendTxs, resUnSpendTx)
		}

		return nil
	})

	return unSpendTxs
}

func (u *UTXOSet)Print()  {
	db, err := Open(UTXO_Set_File_Path, 0600, nil)
	db.View(func(tx *Tx) error {
		bk := tx.Bucket([]byte(UTXO_Set_Name))
		c := bk.Cursor()
		for k,v := c.First();nil != k;k,v=c.Next(){
			unSpendTx := new(UnSpendTxs)
			Deserialize(v, unSpendTx)
			fmt.Printf("TxID[%v] Outs:[%v]" ,hex.EncodeToString(unSpendTx.ID), unSpendTx.Outs)
		}

		return nil
	})

	if nil != err{
		fmt.Printf("Error [%s]", err.Error())
	}
}
