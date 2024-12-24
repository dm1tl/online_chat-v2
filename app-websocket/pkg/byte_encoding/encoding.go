package byteencoding

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Int64ToBytes(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Fatalf("Error while encoding: %v", err)
	}
	return buf.Bytes()
}
