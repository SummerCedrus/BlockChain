package wallets

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	."misc"
	"misc"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}
//新钱包
func NewWallets() *Wallet{
	//生成曲线
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	if nil != err{
		fmt.Errorf("NewWallets Error [%v]", err.Error())
		return nil
	}

	publicKey := append(privateKey.X.Bytes(), privateKey.Y.Bytes() ...)

	return &Wallet{*privateKey, publicKey}
}

func (w *Wallet)GetAddress() []byte{
	pubKeyHash := HashPubKey(w.PublicKey)
	pubKeyHashWithVer := append([]byte(misc.Version), pubKeyHash ...)
	chechRes := CheckSum(pubKeyHashWithVer)

	payload := append(pubKeyHashWithVer, chechRes ...)
	address := misc.Base58Encode(payload)

	return address
}
