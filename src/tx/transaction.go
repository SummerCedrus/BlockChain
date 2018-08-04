package tx
//PubKey: 原生key PubKeyHash:hash后的key,不可逆 address:加上版本和校验和encode后的key可逆

import (
	."misc"
	"crypto/sha256"
	"bytes"
)

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}
//交易输出,可用于计算余额
type TxOutput struct {
	Value int32		// 交易的bt币数量 单位satoshi(0.00000001 BTC)
	PubKeyHash []byte //交易目标的公钥hash
}
//每个交易输入引用一个之前交易的输出，coinbase交易除外
type TxInput struct {
	TxId	[]byte	//引用的交易ID
	OutIndex int32	//引用输出的索引
	Signature string
	PubKey  []byte //交易发起者的公钥
}

type UpSpendTxs struct {
	ID 		[]byte
	Outs	map[int32]TxOutput
}
func (tx *Transaction) GenID(){
	data := Serialize(tx)
	hash := sha256.Sum256(data)
	tx.ID = hash[:]
}

func (txIn *TxInput) CanUnLockOutPut(pubKeyHash []byte) bool{
	txInPubKeyHash := HashPubKey(txIn.PubKey)

	return bytes.Compare(txInPubKeyHash,pubKeyHash) == 0
}
//使用一个address锁定一个输出，即这个输出属于这个address
func (txOut *TxOutput) Lock(address string) {
	pubKeyHash := DecodeAddress(address)
	txOut.PubKeyHash = pubKeyHash
}

func (txOut *TxOutput) IsLockedByPubKeyHash(pubKeyHash []byte) bool{
	return bytes.Compare(txOut.PubKeyHash, pubKeyHash) == 0
}
