package misc

import (
	"strconv"
	"bytes"
	"encoding/gob"
	"fmt"
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
		fmt.Errorf("Encode failed error[%s]", err.Error())
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
		fmt.Errorf("Encode failed error[%s]", err.Error())
		return err
	}

	return nil
}