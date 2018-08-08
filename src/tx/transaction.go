package tx
//PubKey: 原生key PubKeyHash:hash后的key,不可逆 address:加上版本和校验和encode后的key可逆

import (
	."misc"
	"crypto/sha256"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"crypto/rand"
	"fmt"
	"math/big"
	"crypto/elliptic"
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
	Signature []byte
	PubKey  []byte //交易发起者的公钥
}

type UnSpendTxs struct {
	ID 		[]byte	//交易ID
	Outs	map[int32]TxOutput//map[交易输出index]交易输出
}
func (tx *Transaction) GenID(){
	data := Serialize(tx)
	hash := sha256.Sum256(data)
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinBase() bool{
	return len(tx.Vin) == 0
}

func (tx *Transaction) Hash() []byte{
	txCopy := *tx
	txCopy.ID = []byte{}
	hash := sha256.Sum256(Serialize(txCopy))
	return hash[:]
}

//
func (tx *Transaction)Sign(priKey ecdsa.PrivateKey, preTxs map[string]*Transaction){
	if tx.IsCoinBase(){
		return
	}
	txCopy := tx.TrimmedCopy()

	for index, in := range txCopy.Vin{
		preTx := preTxs[hex.EncodeToString(in.TxId)]
		txCopy.Vin[index].Signature = nil
		txCopy.Vin[index].PubKey = preTx.Vout[in.OutIndex].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[index].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &priKey, txCopy.ID)

		if nil != err{
			fmt.Errorf("Sign Failed [%s]",err.Error())
			return
		}
		signature := append(r.Bytes(),s.Bytes()...)
		tx.Vin[index].Signature = signature
	}
}

func (tx *Transaction) Verify(preTxs map[string]*Transaction) bool{
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for index, in := range txCopy.Vin {
		preTx := preTxs[hex.EncodeToString(in.TxId)]
		txCopy.Vin[index].Signature = nil
		txCopy.Vin[index].PubKey = preTx.Vout[in.OutIndex].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[index].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}
//生成个用于签名的交易副本
func (tx *Transaction)TrimmedCopy() Transaction{
	vIn := make([]TxInput, 0)
	for _, in := range tx.Vin{
		vIn = append(vIn, TxInput{in.TxId, in.OutIndex, nil,  nil})
	}
	return Transaction{Vin:vIn, Vout:tx.Vout}
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

func NewTxOutput(value int32, address string)*TxOutput{
	txo := TxOutput{
		Value:value,
	}
	txo.Lock(address)

	return &txo
}
