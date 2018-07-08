package block

import (
	"time"
	"math/big"
	."misc"
	"bytes"
	"math"
	"crypto/sha256"
	"fmt"
)

type Block struct{
	Timestamp int64		//时间戳,自1970-01-01 00:00经过的秒
	PreBlockHash []byte //父hash
	Hash []byte			//自己的hash
	Data []byte			//交易信息
	Nonce	  int64		//计数器
}

func (b *Block) SetHash(hash []byte){
	b.Hash = hash
}

func (b *Block) SetNonce(nonce int64){
	b.Nonce = nonce
}

type ProofOfWork struct {
	Block *Block
	Target *big.Int
}

func (pow *ProofOfWork) Prepare(nonce int64) []byte{
	data := bytes.Join([][]byte{
		pow.Block.PreBlockHash,
		pow.Block.Data,
		Int2Byte(TargetBit),
		Int2Byte(pow.Block.Timestamp),
		Int2Byte(nonce),
	}, []byte{})

	return data
}

func (pow *ProofOfWork) Work(){
	nonce := int64(0)
	result := new(big.Int)
	for nonce <= math.MaxInt64{
		data := pow.Prepare(nonce)
		hash := sha256.Sum256(data)
		result.SetBytes(hash[:])
		if result.Cmp(pow.Target) == -1 {
			pow.Block.SetHash(hash[:])
			pow.Block.SetNonce(nonce)
			fmt.Printf("%x\n",hash)
			return
		}else{
			nonce++
		}
	}
}

func (pow *ProofOfWork) Validate() bool{
	nonce := pow.Block.Nonce
	result := new(big.Int)
	data := pow.Prepare(nonce)
	hash := sha256.Sum256(data)
	result.SetBytes(hash[:])
	return result.Cmp(pow.Target) == -1
}

func NewBlock(preBlockHash []byte, data []byte) *Block{
	beginSec := time.Now().Unix()
	newBlock := &Block{
		Timestamp: time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Data:data,
	}

	pow := newBlock.NewProofOfWork()

	pow.Work()
	endSec := time.Now().Unix()
	fmt.Printf("take %d second\n",endSec-beginSec)
	return newBlock
}

// 创世块
func NewGenesisBlock() *Block{
	nb := NewBlock([]byte{}, []byte("Genesis Block"))
	return nb
}

func (b *Block)NewProofOfWork() *ProofOfWork{
	target := big.NewInt(1)
	target.Lsh(target, 256-TargetBit)
	pow := &ProofOfWork{
		Block: b,
		Target: target,
	}

	return pow
}



