package misc

import "strconv"

func Int2Byte(i int64) []byte{
	return []byte(strconv.FormatInt(i,16))
}