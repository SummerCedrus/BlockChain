package tx

import (
	"fmt"
	"misc"
	"crypto/sha256"
	"blockChain"
)

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

type TxOutput struct {
	Value int32		// 交易的bt币数量 单位satoshi(0.00000001 BTC)
	ScriptPubKey string
}
//每个交易输入引用一个之前交易的输出，coinbase交易除外
type TxInput struct {
	TxId	[]byte	//引用的交易ID
	OutIndex int32	//引用输出的索引
	ScriptSig string
}

func NewTransaction(from, to string, amount int32, bc blockChain.BlockChain) *Transaction{
	inPuts := make([]TxInput, 0)
	outPuts := make([]TxOutput, 0)
	total, outPutInfos := getUnSpendList(from, amount)
	if total < amount{
		fmt.Errorf("Not Enough Coins!")
		return nil
	}

	for txId, outIndex := range outPutInfos{
		inPut := TxInput{
			TxId: []byte(txId),
			OutIndex: outIndex,
			ScriptSig:from,
		}
		inPuts = append(inPuts, inPut)
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

	tx.setID()

	return tx

}

func (tx *Transaction) setID(){
	data := misc.Serialize(tx)
	hash := sha256.Sum256(data)
	tx.ID = hash[:]
}

func (txIn *TxInput) CanUnLockOutPutByAddr(address string) bool{
	return txIn.ScriptSig == address
}

func (txOut *TxOutput) CanBeUnLockByAddr(address string) bool{
	return txOut.ScriptPubKey == address
}