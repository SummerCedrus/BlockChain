package wallets

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	."misc"
	."github.com/bolt"
	"encoding/gob"
	"bytes"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
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

func GetWallet(address string) *Wallet{
	ws := new(Wallets)
	ws.InitWallets()
	w := ws.GetWallet(address)

	return w
}
func (w *Wallet)GetAddress() string{
	pubKeyHash := HashPubKey(w.PublicKey)
	pubKeyHashWithVer := append([]byte{Version}, pubKeyHash ...)
	chechRes := CheckSum(pubKeyHashWithVer)

	payload := append(pubKeyHashWithVer, chechRes ...)
	address := Base58Encode(payload)

	return string(address)
}

func (ws *Wallets)CreateWallet() string{
	wallet := NewWallets()
	if wallet == nil{
		fmt.Errorf("CreateWallet Error!!!")
	}
	addr := wallet.GetAddress()

	ws.Wallets[addr] = wallet

	return addr
}
//根据地址获取钱包
func (ws *Wallets)GetWallet(address string) *Wallet{
	if w, ok := ws.Wallets[address]; ok{
		return w
	}
	fmt.Printf("Cant Find Wallet By Address[%s]", address)

	return nil
}

func (ws *Wallets) InitWallets(){
	ws.Wallets = make(map[string]*Wallet, 0)
	db, err := Open(Wallet_File_Path, 0600, nil)
	if nil != err||nil == db{
		panic(fmt.Sprintf("Open db [%s] failed error[%s]!",Wallet_File_Path, err.Error()))
		return
	}
	defer func() {
		db.Close()
	}()
	err = db.Update(func(tx *Tx) error {
		b := tx.Bucket([]byte(Wallet_Bucket_Name))
		if nil == b{
			b, err = tx.CreateBucket([]byte(Wallet_Bucket_Name))
			if nil != err{
				return err
			}
		}
		//键:address 值:wallet
		err := b.ForEach(func(k, v []byte) error {
			//凡是解码的包里面有interface，必须把interface可能的类型注册下
			gob.Register(elliptic.P256())
			decoder := gob.NewDecoder(bytes.NewReader(v))
			wallet := new(Wallet)
			err := decoder.Decode(wallet)
			if nil != err{
				fmt.Errorf("Decode Wallet Failed address[%s]", k)
				return err
			}

			ws.Wallets[string(k)] = wallet

			return nil
		})
		return err
	})
	 if nil != err{
		 panic(fmt.Sprintf("load wallet data failed error[%s]!", err.Error()))
	 }
	return
}

func (ws *Wallets) SaveWallets(){
	db, err := Open(Wallet_File_Path, 0600, nil)

	if nil != err||nil == db{
		panic(fmt.Sprintf("Open db [%s] failed error[%s]!",Wallet_File_Path, err.Error()))
		return
	}
	defer func() {
		db.Close()
	}()
	db.Update(func(tx *Tx) error {
		b := tx.Bucket([]byte(Wallet_Bucket_Name))
		if nil == b{
			b, err = tx.CreateBucket([]byte(Wallet_Bucket_Name))
			if nil != err{
				return err
			}
		}
		gob.Register(elliptic.P256())
		for addr, wallet := range ws.Wallets{
			buff := new(bytes.Buffer)
			encoder := gob.NewEncoder(buff)
			encoder.Encode(wallet)
			b.Put([]byte(addr), buff.Bytes())
		}
		return nil
	})
}

