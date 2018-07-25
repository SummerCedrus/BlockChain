package tx

import (
	"misc"
	"crypto/sha256"
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

type UpSpendTxs struct {
	ID 		[]byte
	Outs	map[int32]TxOutput
}
func (tx *Transaction) SetID(){
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
