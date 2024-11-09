package packet

import (
	"bytes"
	"encoding/binary"
)

func GetPacketBuffer(_type byte, data []byte) (*bytes.Buffer){
    buf := new(bytes.Buffer)
    buf.Write([]byte{_type})
    binary.Write(buf, binary.LittleEndian, int64(len(data)))
    buf.Write(data)
    return buf
}

