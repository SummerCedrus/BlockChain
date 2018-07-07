package block

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct{
	Timestamp int64		//时间戳,自1970-01-01 00:00经过的秒
	PreBlockHash []byte //父hash
	Hash []byte			//自己的hash
	Data []byte			//交易信息
}

func (b *Block) SetHash(){
	tsStr := []byte(strconv.FormatInt(b.Timestamp,10))
	orgStr := bytes.Join([][]byte{tsStr,b.PreBlockHash,b.Data},[]byte{})
	hashStr := sha256.Sum256(orgStr)

	b.Hash = hashStr[:]
}

func NewBlock(preBlockHash []byte, data []byte) *Block{
	newBlock := &Block{
		Timestamp: time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Data:data,
	}
	newBlock.SetHash()
	return newBlock
}

// 创世块
func NewGenesisBlock() *Block{
	nb := NewBlock([]byte{}, []byte("Genesis Block"))
	return nb
}

