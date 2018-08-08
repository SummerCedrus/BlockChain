package misc

import (
	"strconv"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

func Int2Byte(i int64) []byte{
	return []byte(strconv.FormatInt(i,16))
}

//序列化struct
func Serialize(v interface{}) ([]byte){
	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(v)
	if nil != err{
		fmt.Printf("Encode failed error[%s]", err.Error())
		return buffer.Bytes()
	}
	return buffer.Bytes()
}
//反序列化struct, e为指针类型
func Deserialize(b []byte, e interface{}) error{
	reader := bytes.NewReader(b)
	dec := gob.NewDecoder(reader)
	err := dec.Decode(e)
	if nil != err{
		fmt.Printf("Encode failed error[%s]", err.Error())
		return err
	}

	return nil
}

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func HashPubKey(pubKey []byte) []byte{
	pubKeyHash := sha256.Sum256(pubKey)
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(pubKeyHash[:])
	result := ripemd160Hasher.Sum(nil)

	return result
}

func CheckSum(payload []byte) []byte{
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:4]
}

// Base58Encode encodes a byte array to Base58
func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}

// Base58Decode decodes Base58-encoded data
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func DecodeAddress(address string) []byte{
	payload := Base58Decode([]byte(address))
	pubKeyHash:= payload[1:len(payload)-4]
	return pubKeyHash
}